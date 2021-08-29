/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 8/28/21$ 10:43 PM$
 **/
package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/openzipkin/zipkin-go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vadiminshakov/committer/config"
	"github.com/vadiminshakov/committer/hooks"
	"github.com/vadiminshakov/committer/peer"
	pb "github.com/vadiminshakov/committer/protocol"
	"github.com/vadiminshakov/committer/server"
	"github.com/vadiminshakov/committer/trace"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	COORDINATOR_TYPE = "coordinator"
	FOLLOWER_TYPE    = "follower"
	BADGER_DIR       = "/tmp/badger"
)

const (
	NOT_BLOCKING = iota
	BLOCK_ON_PRECOMMIT_FOLLOWERS
	BLOCK_ON_PRECOMMIT_COORDINATOR
)
const (
	TXNUM = 10000
)

var (
	whitelist = []string{"127.0.0.1"}
	nodes     = map[string][]*config.ConfigShard{
		COORDINATOR_TYPE: {
			{Nodeaddr: "localhost:3000", Role: "coordinator",
				Followers: map[uint32][]string{
					uint32(0): {"localhost:3001", "localhost:3002"},
					uint32(1): {"localhost:3003", "localhost:3004"},
					uint32(2): {"localhost:3005", "localhost:3006"},
				},
				Whitelist: whitelist, CommitType: "two-phase", Timeout: 1000, WithTrace: false},
		},
		FOLLOWER_TYPE: {
			&config.ConfigShard{Nodeaddr: "0+localhost:3001", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
			&config.ConfigShard{Nodeaddr: "0+localhost:3002", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
			&config.ConfigShard{Nodeaddr: "1+localhost:3003", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
			&config.ConfigShard{Nodeaddr: "1+localhost:3004", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
			&config.ConfigShard{Nodeaddr: "2+localhost:3005", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
			&config.ConfigShard{Nodeaddr: "2+localhost:3006", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
		},
	}
	testTable = obtainTestData(shardList, stateList)
)

// 测试数据准备
type StateTransferMsg struct{
	source uint32
	target uint32
	keylist []byte
}
var shardList = []uint32{
	uint32(0), uint32(1), uint32(2),
}
// 需要rlp.encode和decode的结构体，所有元素首字母都要大写！！
type StateDataTest struct {
	Address string
	State string
}
// 这里直接使用harmony的地址, 传输的数据是包含state data的map数据
var stateList = []StateDataTest{
	{"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed", "balance=50"},
	{"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359", "balance=29.59"},
	{"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB", "balance=100"},
	{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", "balance=3000"},
}
//var testTable = obtainTestData(shardList, stateList)

// 设置测试数据
func obtainTestData(shard []uint32, stateList []StateDataTest) []StateTransferMsg{
	var enStateList = make([][]byte, len(stateList))
	fmt.Println("state ori: ", stateList)
	// 获取address的rlp编码
	for i, state := range stateList{
		enState, _ := rlp.EncodeToBytes(state)
		//fmt.Printf("%v → %X\n", state, enState)
		var deState StateDataTest
		err := rlp.DecodeBytes(enState, &deState)
		fmt.Printf("State:err=%+v,val=%+v\n", err, deState)
		enStateList[i] = enState
	}
	var testTable = []StateTransferMsg{
		{shard[0], shard[1], enStateList[0]},
		{shard[0], shard[1], enStateList[1]},
		{shard[0], shard[1], enStateList[2]},
		{shard[1], shard[2], enStateList[3]},
	}

	// 构建更多交易
	TxNum := TXNUM
	for i:=0; i<TxNum-4; i++{
		testTable = append(testTable, testTable[0])
	}

	return testTable
}

func TestDatacreate(t *testing.T){
	testdata := obtainTestData(shardList, stateList)
	t.Log(testdata)
}

// 测试2PC for shard
func TestTwoPhase(t *testing.T){
	log.SetLevel(log.FatalLevel)

	done := make(chan struct{})
	go startnodes(NOT_BLOCKING, done)
	time.Sleep(5 * time.Second) // wait for coordinators and followers to start and establish connections
	fmt.Println("============ STARTED NODES ============")
	//var height uint64 = 0

	// 测试2PC
	for _, coordConfig := range nodes[COORDINATOR_TYPE] {
		// commit type
		if coordConfig.CommitType == "two-phase" {
			log.Println("***\nTEST IN TWO-PHASE MODE\n***")
		} else {
			log.Println("***\nTEST IN THREE-PHASE MODE\n***")
		}
		// 设置trace
		var (
			tracer *zipkin.Tracer
			err    error
		)
		if coordConfig.WithTrace { // 对coordinator地址做trace
			tracer, err = trace.Tracer("client", coordConfig.Nodeaddr)
			if err != nil {
				t.Errorf("no tracer, err: %v", err)
			}
		}
		// 创建coordinator client
		c, err := peer.NewClient(coordConfig.Nodeaddr, tracer)
		if err != nil {
			t.Error(err)
		}

		// 控制coordinator并发的协程数量
		coordPCh := make(chan struct{}, 1500)

		startTime := time.Now() // 获取当前时间
		// 调用stateTransfer请求
		// 使用goroutine并行调用
		var wg sync.WaitGroup
		for index, state := range testTable{
			wg.Add(1)
			coordPCh <- struct{}{}
			go func(s StateTransferMsg, i int) {
				defer wg.Done()
				resp, err := c.StateTransfer(context.Background(), s.source, s.target, s.keylist, uint64(i))

				if err != nil {
					t.Error(err)
				}
				//errChannel <- err
				//respChannel <- resp
				if resp.Type != pb.AckType_ACK {
					fmt.Println("err response type: ", resp.Type)
					t.Error("msg is not acknowledged")
				}
				<- coordPCh
			}(state, index+1)
		}

		wg.Wait()
		fmt.Println("=========== STATETRANSFER END ===========")
		//endTime := time.Now()
		elapsed := time.Since(startTime)
		//fmt.Println("2PC time consumed: ", elapsed)

		// 关闭接收channel
		//close(respChannel)
		//close(errChannel)
		//wgResp.Wait()

		// 验证调用的正确性，连接follower看db中是否有相关数据
		for _, node := range nodes[FOLLOWER_TYPE]{
			if coordConfig.WithTrace {
				tracer, err = trace.Tracer(fmt.Sprintf("%s:%s", coordConfig.Role, coordConfig.Nodeaddr), coordConfig.Nodeaddr)
				if err != nil {
					t.Errorf("no tracer, err: %v", err)
				}
			}
			// 初始化follower client的时候需要对Nodeaddr拆分，因为添加了shard标识
			addrWithShard := strings.SplitN(node.Nodeaddr, "+", 2)
			shardNum, _ := strconv.Atoi(addrWithShard[0])
			cli, err := peer.NewClient(addrWithShard[1], tracer)
			assert.NoError(t, err, "err not nil")

			for _, state := range testTable{
				var deState StateDataTest
				if err := rlp.DecodeBytes(state.keylist, &deState); err!=nil{
					t.Errorf("decode failure, err: %v", err)
				}
				// 判断shard是否参与数据的2PC
				if state.source != uint32(shardNum) && state.target != uint32(shardNum){ break }

				// check values added by nodes 查看state是否都被
				var resp *pb.Value
				resp = new(pb.Value)
				resp, err = cli.Get(context.Background(), deState.Address)
				if assert.Equal(t, resp.Value, state.keylist) && assert.NoError(t, err, "err not nil") {
					t.Log(node.Nodeaddr," values success added by nodes")
				}else{
					t.Log(node.Nodeaddr," values failed added by nodes")
				}
			}
		}

		fmt.Println("===== TX NUM:",len(testTable),"=====")
		// 计算纳秒
		duration := elapsed.Seconds()
		tps := TXNUM/duration
		fmt.Println("2PC time consumed: ", duration)
		fmt.Println("2PC TPS: ", tps)
	}
	done <- struct{}{}
	time.Sleep(2 * time.Second)
}


// 启动节点
func startnodes(block int, done chan struct{}){
	// 创建节点的DB存储目录
	COORDINATOR_BADGER := fmt.Sprintf("%s%s%d", BADGER_DIR, "coordinator", time.Now().UnixNano())
	FOLLOWER_BADGER := fmt.Sprintf("%s%s%d", BADGER_DIR, "follower", time.Now().UnixNano())

	os.Mkdir(COORDINATOR_BADGER, os.FileMode(0777)) // 赋予权限777, 要输入0777才能生效
	os.Mkdir(FOLLOWER_BADGER, os.FileMode(0777))

	var blocking grpc.UnaryServerInterceptor
	switch block {
	case BLOCK_ON_PRECOMMIT_FOLLOWERS:
		blocking = server.PrecommitBlockALL
	case BLOCK_ON_PRECOMMIT_COORDINATOR:
		blocking = server.PrecommitBlockCoordinator
	}

	// 启动所有follower server
	for i, node := range nodes[FOLLOWER_TYPE] {
		// create db dir
		os.Mkdir(node.DBPath, os.FileMode(0777))
		node.DBPath = fmt.Sprintf("%s%s%s", FOLLOWER_BADGER, strconv.Itoa(i), "~")
		// start follower
		hooks, err := hooks.GetHookF()
		if err != nil {
			panic(err)
		}
		// 创建follower时，只需要follower自己的配置信息
		followerServer, err := server.NewCommitServerShard(node, hooks...)
		if err != nil {
			panic(err)
		}

		if block == 0 { // 没有阻塞的情况
			go followerServer.Run(server.WhiteListCheckerShard)
		} else {
			go followerServer.Run(server.WhiteListCheckerShard, blocking)
		}
		defer followerServer.Stop() //结束函数之前调用
	}
	time.Sleep(3 * time.Second)

	// 启动coordinator server
	for i, coordConfig := range nodes[COORDINATOR_TYPE] {
		// create db dir
		os.Mkdir(coordConfig.DBPath, os.FileMode(0777))
		coordConfig.DBPath = fmt.Sprintf("%s%s%s", COORDINATOR_BADGER, strconv.Itoa(i), "~")
		// start coordinator
		hooks, err := hooks.GetHookF()
		if err != nil {
			panic(err)
		}
		// 创建coordinator时，需要传入自己配置+相关follower的信息作为client端
		coordServer, err := server.NewCommitServerShard(coordConfig, hooks...)
		if err != nil {
			panic(err)
		}

		if block == 0 {
			go coordServer.Run(server.WhiteListCheckerShard)
		} else {
			go coordServer.Run(server.WhiteListCheckerShard, blocking)
		}

		defer coordServer.Stop()
	}

	<- done
	// prune
	os.RemoveAll(BADGER_DIR)
}

// 测试getConfig的分词
func TestGetConfig(t *testing.T){
	str := "0+localhost:3001,0+localhost:3002,1+localhost:3003,1+localhost:3004,2+localhost:3005,2+localhost:3006"

	shardArray := strings.Split(str, ",")
	followersArray := make(map[uint32][]string)
	for _, item := range shardArray{
		//fmt.Println(item)
		follower := strings.SplitN(item, ">", 2)
		intNum, _ := strconv.Atoi(follower[0])
		followersArray[uint32(intNum)]=append(followersArray[uint32(intNum)], follower[1])
	}
	for key, item := range followersArray{
		for _, item1 := range item{
			fmt.Println(key)
			fmt.Println(item1)
		}
	}
}