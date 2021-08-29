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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/db"
	pb "github.com/vadiminshakov/committer/protocol"
)

// follower server处理propose逻辑
func (s *ServerShard) ProposeHandler(ctx context.Context, req *pb.ProposeRequest, hook func(req *pb.ProposeRequest) bool) (*pb.Response, error) {
	var (
		response *pb.Response
	)
	// follower处理propose时初始化type, 是source还是target
	s.FollowerType = pb.FollowerType_SOURCE
	if req.Target == s.ShardID{s.FollowerType = pb.FollowerType_TARGET}

	// hook函数检查请求是否正确（目前都直接返回true）
	if hook(req) {
		/**
			重要逻辑： 同步执行分片之间的数据传输协议，只有：数据传输完毕+验证 后才能return
		 */
		//fmt.Printf("[ProposeHandler] propose received: SourceShard=[%d], TargetShard=[%d] \n", req.Source, req.Target)


		// 执行数据传输dataTransmission (思考怎么异步执行，等传输结束的时候再返回)
		// follower server节点设置cache事务，记录事务的index，便于commit阶段使用
		//log.Infof("received: %s=%s\n", req.Key, string(req.Value))

		// 创建一个新channel放入cache
		proposeCh := make(chan int)

		s.NodeCache.SetWithC(req.Index, req.Source, req.Target, req.Keylist, proposeCh, 1)

		// follower server节点设置cache状态数据[state data + proof] （还没写）
		response = &pb.Response{Type: pb.AckType_ACK, Index: req.Index, Ftype: s.FollowerType}

		/**
			source & target 不同的处理
		 */
		// source follower锁定state
		if s.FollowerType==pb.FollowerType_SOURCE{
			log.Infof("Source freeze the state of tx: %d", req.Index)
			// 等待区块链进程通知propose处理完成
			var freezeState bool
			freezeState = true
			Loop:
				for {
					select {
					case sig := <- proposeCh: // 接收到事件处理完成的消息，就顺序执行后续操作
						if sig == 2{ // 写入2表示成功
							//fmt.Println(s.Config.Nodeaddr, " Get proposeCh channel: index=[", req.Index,"]" )
							close(proposeCh)
							break Loop
						}
					//case <-time.After(time.Duration(10)*time.Second): // 超时，事件失败
					//	freezeState = false
					//	close(proposeCh)
					//	break Loop
					}
				}

			// 等待执行成功的事件通知
			if freezeState{ // 锁定成功
				log.Info("[ProposeHandler]freeze success")
				//fmt.Println("[ProposeHandler] freeze success")
			}else{ // 锁定失败
				log.Info("[ProposeHandler] freeze failure")
				//fmt.Println("[ProposeHandler] freeze failure")
				response = &pb.Response{Type: pb.AckType_NACK, Index: req.Index, Ftype: s.FollowerType}
			}

		}else{ // target follower 验证验证state
			log.Infof("Target verify the state of tx: %d", req.Index)
			var verifyState bool
			verifyState = true
			if verifyState{ // 验证成功
				log.Info("[ProposeHandler] verify success")
				//fmt.Println("[ProposeHandler] verify success")
			}else{ // 验证失败
				log.Info("[ProposeHandler] verify failure")
				//fmt.Println("[ProposeHandler] verify failure")
				response = &pb.Response{Type: pb.AckType_NACK, Index: req.Index, Ftype: s.FollowerType}
			}
		}

	} else {
		// 不满足hook请求时，返回NACK
		response = &pb.Response{Type: pb.AckType_NACK, Index: req.Index, Ftype: s.FollowerType}
	}
	// 如果follower server的height比请求的还高，该请求是无效请求
	// 不关注index和height的处理顺序问题
	//if s.Height > req.Index {
	//	response = &pb.Response{Type: pb.AckType_NACK, Index: s.Height, Ftype: s.FollowerType}
	//}

	//fmt.Println("PROPOSE END", s.Config.Nodeaddr, "index=[", req.Index,"]")
	return response, nil
}

//func (s *ServerShard) PrecommitHandler(ctx context.Context, req *pb.PrecommitRequest) (*pb.Response, error) {
//	return &pb.Response{Type: pb.AckType_ACK}, nil
//}

// follower server 处理commit逻辑
func (s *ServerShard) CommitHandler(ctx context.Context, req *pb.CommitRequest, hook func(req *pb.CommitRequest) bool, db db.Database) (*pb.Response, error) {
	var (
		response *pb.Response
	)

	if hook(req) {
		//fmt.Printf("[CommitHandler] %s Committing tx Index: %d\n", s.Addr, req.Index)

		// 查看cache中是否有request相关的数据
		_, _, stateList, ok := s.NodeCache.Get(req.Index)
		//source, target, keylist, ok := s.NodeCache.Get(req.Index)
		if !ok {
			s.NodeCache.Delete(req.Index)
			fmt.Println("[CommitHandler] TX index: ",req.Index," Error: no value in node cache")
			return &pb.Response{
				Type:  pb.AckType_NACK}, errors.New(fmt.Sprintf("no value in node cache on the index %d", req.Index))
		}
		/**
		source & target 不同的处理
		*/
		// source follower将事务写入2PC db
		var deState StateData
		if err := rlp.DecodeBytes(stateList, &deState); err != nil{
			log.Info("decode error: ", err)
		}
		if s.FollowerType==pb.FollowerType_SOURCE{
			log.Info("Source write tx db success")
			//fmt.Println("Source write tx db success")
			if err := db.Put(deState.Address, stateList); err != nil {
				return nil, err
			}
		}else{ // target follower 将事务写入2PC db
			log.Info("Target write tx db success")
			//fmt.Println("Target write tx db success")
			// target将状态数据插入区块链
			if err := db.Put(deState.Address, stateList); err != nil {
				return nil, err
			}
		}

		// commit之后，将交易从cache删除
		s.NodeCache.Delete(req.Index)
		response = &pb.Response{Type: pb.AckType_ACK, Index: s.Height, Ftype: s.FollowerType}
	} else {
		s.NodeCache.Delete(req.Index)
		response = &pb.Response{Type: pb.AckType_NACK, Index: s.Height, Ftype: s.FollowerType}
	}
	// 结束processor执行
	//s.Processor.ProcessEndCh <- 1

	return response, nil
}