/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 7/21/21$ 7:14 AM$
 **/
package peer

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/pkg/errors"
	pb "github.com/vadiminshakov/committer/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)

type CommitClientShard struct {
	Connection pb.CommitClient
	Tracer     *zipkin.Tracer
}

// New creates instance of peer client.
// 'addr' is a coordinator network address (host + port).
func NewClient(addr string, tracer *zipkin.Tracer) (*CommitClientShard, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	// 增加最小超时时间MinConnectTimeout
	connParams := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay: 100 * time.Millisecond,
			MaxDelay:  5 * time.Second,
		},
		MinConnectTimeout: 200 * time.Millisecond,
	}
	if tracer != nil {
		conn, err = grpc.Dial(addr, grpc.WithConnectParams(connParams), grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)))
	} else {
		conn, err = grpc.Dial(addr, grpc.WithConnectParams(connParams), grpc.WithInsecure())
	}
	if err != nil {
		time.Sleep(10*time.Millisecond)
		return nil, errors.Wrap(err, "failed to connect")
	}
	return &CommitClientShard{Connection: pb.NewCommitClient(conn), Tracer: tracer}, nil
}

func (client *CommitClientShard) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Propose")
		defer span.Finish()
	}
	//fmt.Println("[Client] execute client propose: [index=",req.Index,"]")
	return client.Connection.Propose(ctx, req)
}

func (client *CommitClientShard) Precommit(ctx context.Context, req *pb.PrecommitRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Precommit")
		defer span.Finish()
	}
	//fmt.Println("execute client precommit")
	return client.Connection.Precommit(ctx, req)
}

func (client *CommitClientShard) Commit(ctx context.Context, req *pb.CommitRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Commit")
		defer span.Finish()
	}
	//fmt.Println("[Client] execute client commit: [index=",req.Index,"]")
	return client.Connection.Commit(ctx, req)
}

// NodeInfo gets info about current node height
func (client *CommitClientShard) NodeInfo(ctx context.Context) (*pb.Info, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "NodeInfo")
		defer span.Finish()
	}
	return client.Connection.NodeInfo(ctx, &empty.Empty{})
}

// Get queries value of specific key
func (client *CommitClientShard) Get(ctx context.Context, key string) (*pb.Value, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Get")
		defer span.Finish()
	}
	return client.Connection.Get(ctx, &pb.Msg{Key: key})
}

/**
	dynamic sharding修改
**/
// StateTransfer将一组state从 source followers转到 target followers
// 2PC StateTransfer client端协议
func (client *CommitClientShard) StateTransfer(ctx context.Context, source uint32, target uint32, keylist []byte, index uint64) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "StateTransfer")
		defer span.Finish()
	}
	//fmt.Println("execute client stateTransfer")
	return client.Connection.StateTransfer(ctx, &pb.TransferRequest{Source: source, Target: target, Keylist: keylist, Index: index})
}
