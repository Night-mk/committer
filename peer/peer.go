package peer

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/pkg/errors"
	pb "github.com/vadiminshakov/committer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)

type CommitClient struct {
	Connection pb.CommitClient
	Tracer     *zipkin.Tracer
}

// New creates instance of peer client.
// 'addr' is a coordinator network address (host + port).
func New(addr string, tracer *zipkin.Tracer) (*CommitClient, error) {
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
	return &CommitClient{Connection: pb.NewCommitClient(conn), Tracer: tracer}, nil
}

func (client *CommitClient) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Propose")
		defer span.Finish()
	}
	fmt.Println("execute client propose")
	return client.Connection.Propose(ctx, req)
}

func (client *CommitClient) Precommit(ctx context.Context, req *pb.PrecommitRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Precommit")
		defer span.Finish()
	}
	fmt.Println("execute client precommit")
	return client.Connection.Precommit(ctx, req)
}

func (client *CommitClient) Commit(ctx context.Context, req *pb.CommitRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Commit")
		defer span.Finish()
	}
	fmt.Println("execute client commit: [index=",req.Index,"]")
	return client.Connection.Commit(ctx, req)
}

// Put sends key/value pair to peer (it should be a coordinator).
// The coordinator reaches consensus and all peers commit the value.
func (client *CommitClient) Put(ctx context.Context, key string, value []byte) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Put")
		defer span.Finish()
	}
	fmt.Println("execute client put")
	return client.Connection.Put(ctx, &pb.Entry{Key: key, Value: value})
}

// NodeInfo gets info about current node height
func (client *CommitClient) NodeInfo(ctx context.Context) (*pb.Info, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "NodeInfo")
		defer span.Finish()
	}
	return client.Connection.NodeInfo(ctx, &empty.Empty{})
}

// Get queries value of specific key
func (client *CommitClient) Get(ctx context.Context, key string) (*pb.Value, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Get")
		defer span.Finish()
	}
	return client.Connection.Get(ctx, &pb.Msg{Key: key})
}
