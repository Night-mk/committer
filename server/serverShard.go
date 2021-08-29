/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 7/21/21$ 7:24 AM$
 **/
package server

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/cache"
	"github.com/vadiminshakov/committer/config"
	"github.com/vadiminshakov/committer/db"
	"github.com/vadiminshakov/committer/peer"
	pb "github.com/vadiminshakov/committer/protocol"
	"github.com/vadiminshakov/committer/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type OptionShard func(server *ServerShard) error

const COORDINATOR_SHARDID = uint32(10000)

type StateData struct {
	Address string
	State string
}

// commitServer端，实现所有服务接口
// Server holds server instance, node config and connections to followers (if it's a coordinator node)
type ServerShard struct {
	pb.UnimplementedCommitServer
	Addr                 string
	ShardID              uint32 // 当前server的shardID, 暂定coordinator的shardID是10000
	FollowerType         pb.FollowerType
	Followers            []*peer.CommitClientShard
	ShardFollowers       map[uint32][]*peer.CommitClientShard // 带shard的全体follower集合
	Config               *config.ConfigShard
	GRPCServer           *grpc.Server
	DB                   db.Database
	DBPath               string
	ProposeShardHook          func(req *pb.ProposeRequest) bool
	CommitShardHook           func(req *pb.CommitRequest) bool
	NodeCache            *cache.CacheShard
	Height               uint64
	cancelCommitOnHeight map[uint64]bool
	mu                   sync.RWMutex
	Tracer               *zipkin.Tracer

	Processor            Process
}

// 2PC 第一阶段
func (s *ServerShard) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.Response, error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "ProposeHandle")
		defer span.Finish()
	}
	//fmt.Println("[server] ",s.Addr, " execute server propose")
	// 设置在某个Height是否需要cancel
	s.SetProgressForCommitPhase(req.Index, false)

	return s.ProposeHandler(ctx, req, s.ProposeShardHook)
}

//func (s *ServerShard) Precommit(ctx context.Context, req *pb.PrecommitRequest) (*pb.Response, error) {
//	var span zipkin.Span
//	if s.Tracer != nil {
//		span, _ = s.Tracer.StartSpanFromContext(ctx, "PrecommitHandle")
//		defer span.Finish()
//	}
//	fmt.Println("execute server precommit")
//	if s.Config.CommitType == THREE_PHASE {
//		ctx, _ = context.WithTimeout(context.Background(), time.Duration(s.Config.Timeout)*time.Millisecond)
//		go func(ctx context.Context) {
//		ForLoop:
//			for {
//				select {
//				case <-ctx.Done():
//					md := metadata.Pairs("mode", "autocommit")
//					ctx := metadata.NewOutgoingContext(context.Background(), md)
//					if !s.HasProgressForCommitPhase(req.Index) {
//						s.Commit(ctx, &pb.CommitRequest{Index: s.Height})
//						log.Warn("commit without coordinator after timeout")
//					}
//					break ForLoop
//				}
//			}
//		}(ctx)
//	}
//	return s.PrecommitHandler(ctx, req)
//}

// 2PC 第二阶段
func (s *ServerShard) Commit(ctx context.Context, req *pb.CommitRequest) (resp *pb.Response, err error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "CommitHandle")
		defer span.Finish()
	}
	//fmt.Println("[Server] execute server commit: [index=",req.Index,"]")
	if s.Config.CommitType == THREE_PHASE {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "no metadata")
		}

		meta := md["mode"]

		if len(meta) == 0 {
			s.SetProgressForCommitPhase(s.Height, true) // set flag for cancelling 'commit without coordinator' action, coz coordinator responded actually
			resp, err = s.CommitHandler(ctx, req, s.CommitShardHook, s.DB)
			if err != nil {
				return nil, err
			}
			if resp.Type == pb.AckType_ACK {
				atomic.AddUint64(&s.Height, 1)
			}
			return
		}

		if req.IsRollback {
			s.rollback()
		}
	} else {
		// 两阶段commit
		//fmt.Println("[Server] Commit (two phase)")
		resp, err = s.CommitHandler(ctx, req, s.CommitShardHook, s.DB)

		if err != nil {
			return
		}
		if resp.Type == pb.AckType_ACK {
			atomic.AddUint64(&s.Height, 1)
		}
		return
	}
	return
}

func (s *ServerShard) Get(ctx context.Context, req *pb.Msg) (*pb.Value, error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "GetHandle")
		defer span.Finish()
	}
	// 从database通过key获取数据
	value, err := s.DB.Get(req.Key)
	if err != nil {
		return nil, err
	}
	return &pb.Value{Value: value}, nil
}

/**
	dynamic sharding修改
**/
// 2PC StateTransfer server端协议，coordinator调用
// Single 调用
func (s *ServerShard) StateTransferS(ctx context.Context, req *pb.TransferRequest) (*pb.Response, error){
	var (
		response *pb.Response
		err      error
		span     zipkin.Span
		//resp_num int
	)
	// trace
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "StateTransferHandle")
		defer span.Finish()
	}
	// 选择使用的协议，目前暂定只用2PC(propose, commit)
	var ctype pb.CommitType
	if s.Config.CommitType == THREE_PHASE {
		ctype = pb.CommitType_THREE_PHASE_COMMIT
	} else {
		ctype = pb.CommitType_TWO_PHASE_COMMIT
	}

	/**=======2PC协议流程=======*/
	// propose阶段

	// 打印address list
	// 反序列化req.Keylist为string slice, log输出, 调用ethereum的rlp编码进行解码
	var deState StateData
	err = rlp.DecodeBytes(req.Keylist, &deState)
	if err != nil{
		log.Info("decode error: ", err)
	}
	//fmt.Printf("[StateTransfer] Propose phase: index=[%d] source=[%d], target=[%d], state=[%+v] \n", req.Index, req.Source, req.Target, deState.Address)

	// 根据source和target选择followers组
	currentFollowers := s.ObtainFollowers(req.Source, req.Target)
	//fmt.Println("selected followers: ", s.ShardFollowers[req.Source])
	//resp_num = 0
	// 对每个follower循环处理（可以并发？）
	for _, follower := range currentFollowers {
		if s.Tracer != nil {
			span, ctx = s.Tracer.StartSpanFromContext(ctx, "Propose")
		}
		var (
			response *pb.Response
			err      error
		)
		for response == nil || response != nil && response.Type == pb.AckType_NACK {
			// 每个follower（source, target）用grpc调用coordinator给他们发送消息
			response, err = follower.Propose(ctx, &pb.ProposeRequest{Source: req.Source, Target: req.Target, Keylist: req.Keylist, CommitType: ctype, Index: req.Index})

			// 停止trace
			if s.Tracer != nil && span != nil {
				span.Finish()
			}

			// 根据coordinator的Height更新当前follower的Height参数
			if response != nil && response.Index > s.Height {
				log.Warnf("сoordinator has stale height [%d], update to [%d] and try to send again", s.Height, response.Index)
				s.Height = response.Index
			}

			// 统计response数量
			//if response != nil && response.Type==pb.AckType_ACK{
			//	resp_num += 1
			//}
		}

		if err != nil {
			log.Errorf(err.Error())
			fmt.Println("propose phase error")
			return &pb.Response{Type: pb.AckType_NACK}, nil
		}
		// coordinator server设置cache，用于记录index=Height的事务
		// index设置成交易的绝对编号
		s.NodeCache.Set(req.Index, req.Source, req.Target, req.Keylist)
		if response.Type != pb.AckType_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}

	// the coordinator got all the answers, so it's time to persist msg and send commit command to followers
	_, _, statelist, ok := s.NodeCache.Get(req.Index)
	//source, target, keylist, ok := s.NodeCache.Get(s.Height)
	if !ok {
		fmt.Println("[StateTransfer] can not find msg in cache for index=[", req.Index,"]")
		return nil, status.Error(codes.Internal, "can't to find msg in the coordinator's cache")
	}
	// coordinator server将事务写入
	if err = s.DB.Put(deState.Address, statelist); err != nil {
		return &pb.Response{Type: pb.AckType_NACK}, status.Error(codes.Internal, "failed to save msg on coordinator")
	}

	// commit阶段
	//fmt.Printf("[StateTransfer] Commit phase: index=[%d] source=[%d], target=[%d], state=[%s] \n", req.Index, req.Source, req.Target, deState)
	for _, follower := range currentFollowers{
		if s.Tracer != nil {
			span, ctx = s.Tracer.StartSpanFromContext(ctx, "Commit")
		}
		// coordinator获取到commit的response
		response, err = follower.Commit(ctx, &pb.CommitRequest{Index: req.Index})

		if s.Tracer != nil && span != nil {
			span.Finish()
		}

		if err != nil {
			log.Errorf(err.Error())
			return &pb.Response{Type: pb.AckType_NACK}, nil
		}
		if response.Type != pb.AckType_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}

	fmt.Printf("committed state for tx_index=[%d] \n", s.Height)
	// increase height for next round
	atomic.AddUint64(&s.Height, 1)

	return &pb.Response{Type: pb.AckType_ACK}, nil
}

/**
	stateTransfer 并行化方案
 */
func (s *ServerShard) StateTransfer(ctx context.Context, req *pb.TransferRequest) (*pb.Response, error){
	var (
		response *pb.Response
		err      error
		span     zipkin.Span
	)
	// trace
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "StateTransferHandle")
		defer span.Finish()
	}
	// 选择使用的协议，目前暂定只用2PC(propose, commit)
	var ctype pb.CommitType
	if s.Config.CommitType == THREE_PHASE {
		ctype = pb.CommitType_THREE_PHASE_COMMIT
	} else {
		ctype = pb.CommitType_TWO_PHASE_COMMIT
	}

	/**=======2PC协议流程=======*/
	// propose阶段
	// 打印address list
	// 反序列化req.Keylist为string slice, log输出, 调用ethereum的rlp编码进行解码
	var deState StateData
	err = rlp.DecodeBytes(req.Keylist, &deState)
	if err != nil{
		log.Info("decode error: ", err)
	}
	//fmt.Printf("[StateTransfer] Propose phase: index=[%d] source=[%d], target=[%d], state=[%+v] \n", req.Index, req.Source, req.Target, deState.Address)

	// 根据source和target选择followers组
	currentFollowers := s.ObtainFollowers(req.Source, req.Target)

	// 对每个follower并发处理
	var wg sync.WaitGroup
	proposeResChan := make(chan *pb.Response, 1000)
	proposeErrChan := make(chan error, 1000)
	//proposeDone := make(chan bool)
	for _, follower := range currentFollowers {
		wg.Add(1)
		go func(follower *peer.CommitClientShard) {
			if s.Tracer != nil {
				span, ctx = s.Tracer.StartSpanFromContext(ctx, "Propose")
			}
			var (
				response *pb.Response
				err      error
			)
			for response == nil || response != nil && response.Type == pb.AckType_NACK {
				// 每个follower（source, target）用grpc调用coordinator给他们发送消息
				response, err = follower.Propose(ctx, &pb.ProposeRequest{Source: req.Source, Target: req.Target, Keylist: req.Keylist, CommitType: ctype, Index: req.Index})
				// 停止trace
				if s.Tracer != nil && span != nil {
					span.Finish()
				}
				// 根据coordinator的Height更新当前follower的Height参数
				if response != nil && response.Index > s.Height {
					log.Warnf("сoordinator has stale height [%d], update to [%d] and try to send again", s.Height, response.Index)
					s.Height = response.Index
				}
			}

			proposeResChan <- response
			proposeErrChan <- err

			if err != nil {
				//log.Errorf(err.Error())
				//fmt.Println("propose phase error")
				return
				//return &pb.Response{Type: pb.AckType_NACK}, nil
			}
			// coordinator server设置cache，用于记录index=Height的事务
			// index设置成交易的绝对编号
			s.NodeCache.Set(req.Index, req.Source, req.Target, req.Keylist)

			if response.Type != pb.AckType_ACK {
				fmt.Println("follower not acknowledged msg")
				return
				//return nil, status.Error(codes.Internal, "follower not acknowledged msg")
			}
			wg.Done()
		}(follower)
	}

	// 等待所有follower反馈完成，再进入commit过程
	wg.Wait()
	close(proposeErrChan)
	close(proposeResChan)

	for err := range proposeErrChan{
		//fmt.Println("=============execute proposeErrChan=====================")
		if err != nil {
			log.Errorf(err.Error())
			return &pb.Response{Type: pb.AckType_NACK}, nil
		}
	}
	for res := range proposeResChan{
		//fmt.Println("=============execute proposeResChan=====================")
		if res.Type != pb.AckType_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}

	//fmt.Printf(s.Config.Nodeaddr,"PROPOSE END index=[", req.Index,"]")

	// the coordinator got all the answers, so it's time to persist msg and send commit command to followers
	_, _, statelist, ok := s.NodeCache.Get(req.Index)
	//source, target, keylist, ok := s.NodeCache.Get(s.Height)
	if !ok {
		fmt.Println("[StateTransfer] can not find msg in cache for index=[", req.Index,"]")
		return nil, status.Error(codes.Internal, "can't to find msg in the coordinator's cache")
	}
	// coordinator server将事务写入
	if err = s.DB.Put(deState.Address, statelist); err != nil {
		return &pb.Response{Type: pb.AckType_NACK}, status.Error(codes.Internal, "failed to save msg on coordinator")
	}

	// commit阶段
	//fmt.Printf("[StateTransfer] Commit phase: index=[%d] source=[%d], target=[%d], state=[%s] \n", req.Index, req.Source, req.Target, deState)
	var wg1 sync.WaitGroup
	commitResChan := make(chan *pb.Response, 1000)
	commitErrChan := make(chan error, 1000)
	for _, follower := range currentFollowers{
		wg1.Add(1)
		go func(follower *peer.CommitClientShard) {
			if s.Tracer != nil {
				span, ctx = s.Tracer.StartSpanFromContext(ctx, "Commit")
			}
			// coordinator获取到commit的response
			response, err = follower.Commit(ctx, &pb.CommitRequest{Index: req.Index})
			if s.Tracer != nil && span != nil {
				span.Finish()
			}
			commitResChan <- response
			commitErrChan <- err
			if err != nil {
				log.Errorf(err.Error())
				return
				//return &pb.Response{Type: pb.AckType_NACK}, nil
			}
			if response.Type != pb.AckType_ACK {
				return
				//return nil, status.Error(codes.Internal, "follower not acknowledged msg")
			}
			wg1.Done()
		}(follower)
	}

	wg1.Wait()
	close(commitResChan)
	close(commitErrChan)

	// 对所有follower节点的返回值进行处理
	for res := range commitResChan{
		if res.Type != pb.AckType_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}
	for err := range commitErrChan{
		if err != nil {
			log.Errorf(err.Error())
			fmt.Println("propose phase error")
			return &pb.Response{Type: pb.AckType_NACK}, nil
		}
	}

	log.Infof("committed state for tx_index=[%d] \n", req.Index)
	fmt.Printf("committed state for tx_index=[%d] \n", req.Index)
	// increase height for next round
	atomic.AddUint64(&s.Height, 1)

	// 返回调用端最终结果
	return &pb.Response{Type: pb.AckType_ACK}, nil
}


// 选择source和target shard的节点作为follower, server.Config保存全部的followers的信息
/**
	原本的config
	map[string][]*config.Config{
		COORDINATOR_TYPE: {
			{Nodeaddr: "localhost:3000", Role: "coordinator",
				Followers: []string{"localhost:3001", "localhost:3002", "localhost:3003", "localhost:3004", "localhost:3005"},
				Whitelist: whitelist, CommitType: "two-phase", Timeout: 1000, WithTrace: false},
			{Nodeaddr: "localhost:5000", Role: "coordinator",
				Followers: []string{"localhost:3001", "localhost:3002", "localhost:3003", "localhost:3004", "localhost:3005"},
				Whitelist: whitelist, CommitType: "three-phase", Timeout: 1000, WithTrace: false},
		}}
	有shard的config
	map[string][]*config.Config{
		COORDINATOR_TYPE: {
			{Nodeaddr: "localhost:3000", Role: "coordinator",
				Followers: map[uint32][]string{
					0: []string{"localhost:3001", "localhost:3002"},
					1: []string{"localhost:3003", "localhost:3004"},
					2: []string{"localhost:3005", "localhost:3006"},
				},
				Whitelist: whitelist, CommitType: "two-phase", Timeout: 1000, WithTrace: false},
		}}
**/
func (s *ServerShard)ObtainFollowers(source uint32, target uint32) (follower []*peer.CommitClientShard){
	for _, node := range s.ShardFollowers[source]{
		follower = append(follower, node)
	}
	for _, node := range s.ShardFollowers[target]{
		follower = append(follower, node)
	}
	return follower
}

// 根据config，创建server端instance
func NewCommitServerShard(conf *config.ConfigShard, opts ...OptionShard) (*ServerShard, error){
	var wg sync.WaitGroup
	wg.Add(1)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true, // Seems like automatic color detection doesn't work on windows terminals
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	// 初始化serverShard
	//server := &ServerShard{Addr: conf.Nodeaddr, DBPath: conf.DBPath}
	server := &ServerShard{DBPath: conf.DBPath}
	// 对coordinator和follower分开初始化
	if conf.Role == "coordinator"{
		server.Addr = conf.Nodeaddr
		server.ShardID = COORDINATOR_SHARDID
	}else{
		// follower分解nodeaddr
		addrWithShard := strings.SplitN(conf.Nodeaddr, "+", 2)
		server.Addr = addrWithShard[1]
		intNum, _ := strconv.Atoi(addrWithShard[0])
		server.ShardID = uint32(intNum)
	}

	var err error
	// 查看option是否有问题
	for _, option := range opts {
		err = option(server)
		if err != nil {
			return nil, err
		}
	}
	// 管理trace
	if conf.WithTrace {
		nodeaddr := conf.Nodeaddr
		if conf.Role != "coordinator"{
			nodeaddr = strings.SplitN(conf.Nodeaddr, "+", 2)[1]
		}
		// get Zipkin tracer
		server.Tracer, err = trace.Tracer(fmt.Sprintf("%s:%s", conf.Role, conf.Nodeaddr), nodeaddr)
		if err != nil {
			return nil, err
		}
	}

	// 初始化配置里的follower client（coordinator中的follower）
	server.ShardFollowers = make(map[uint32][]*peer.CommitClientShard)
	for key, nodes := range conf.Followers {
		for _, node := range nodes{
			//fmt.Println("config init node: ", node)
			cli, err := peer.NewClient(node, server.Tracer)
			if err != nil {
				return nil, err
			}
			server.Followers = append(server.Followers, cli)
			server.ShardFollowers[key] = append(server.ShardFollowers[key], cli)
		}
	}

	// 在server.config中标记coordinator
	server.Config = conf
	if conf.Role == "coordinator" {
		server.Config.Coordinator = server.Addr
	}

	// 初始化DB, cache, rollback参数
	server.DB, err = db.New(conf.DBPath)
	if err != nil {
		return nil, err
	}
	server.NodeCache = cache.NewCache()
	server.cancelCommitOnHeight = map[uint64]bool{}

	if server.Config.CommitType == TWO_PHASE {
		log.Info("two-phase-commit mode enabled")
	} else {
		log.Info("three-phase-commit mode enabled")
	}

	err = checkServerField(server)

	server.Processor.NewProcessor(server.NodeCache)
	wg.Done()
	wg.Wait()

	// 执行消息本地处理操作(只有follower执行)
	if conf.Role == "follower"{
		go server.Processor.loop()
	}

	return server, err
}

// 查看是否有DB
func checkServerField(server *ServerShard) error {
	if server.DB == nil {
		return errors.New("database is not selected")
	}
	return nil
}

// Run starts non-blocking GRPC server
func (s *ServerShard) Run(opts ...grpc.UnaryServerInterceptor) {
	var err error

	if s.Config.WithTrace {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...), grpc.StatsHandler(zipkingrpc.NewServerHandler(s.Tracer)))
	} else {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))
	}
	pb.RegisterCommitServer(s.GRPCServer, s)

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listening on tcp://%s", s.Addr)

	go s.GRPCServer.Serve(l)
}

// Stop stops server
func (s *ServerShard) Stop() {
	if s.Config.Role == "follower"{
		s.Processor.ProcessEndCh <- 1
	}
	// stop的时候关闭processor的通道
	s.Processor.CloseCh()
	log.Info("stopping server")
	s.GRPCServer.GracefulStop()
	if err := s.DB.Close(); err != nil {
		log.Infof("failed to close db, err: %s\n", err)
	}
	log.Info("server stopped")
}

func (s *ServerShard) rollback() {
	s.NodeCache.Delete(s.Height)
}

func (s *ServerShard) SetProgressForCommitPhase(height uint64, docancel bool) {
	s.mu.Lock()
	s.cancelCommitOnHeight[height] = docancel
	s.mu.Unlock()
}

func (s *ServerShard) HasProgressForCommitPhase(height uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cancelCommitOnHeight[height]
}