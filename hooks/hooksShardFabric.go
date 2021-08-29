/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 7/21/21$ 8:26 AM$
 **/
package hooks

import (
	"github.com/vadiminshakov/committer/hooks/src"
	pb "github.com/vadiminshakov/committer/protocol"
	serverShard "github.com/vadiminshakov/committer/server"
)

type ProposeShardHook func(req *pb.ProposeRequest) bool
type CommitShardHook func(req *pb.CommitRequest) bool

func GetHookF() ([]serverShard.OptionShard, error) {
	proposeShardHook := func(f ProposeShardHook) func(*serverShard.ServerShard) error {
		return func(server *serverShard.ServerShard) error {
			server.ProposeShardHook = f
			return nil
		}
	}

	commitShardHook := func(f CommitShardHook) func(*serverShard.ServerShard) error {
		return func(server *serverShard.ServerShard) error {
			server.CommitShardHook = f
			return nil
		}
	}

	return []serverShard.OptionShard{proposeShardHook(src.ProposeShard), commitShardHook(src.CommitShard)}, nil
}