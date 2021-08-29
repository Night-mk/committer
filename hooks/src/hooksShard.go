/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 7/21/21$ 7:32 AM$
 **/
package src

import (
	log "github.com/sirupsen/logrus"
	pb "github.com/vadiminshakov/committer/protocol"
)

func ProposeShard(req *pb.ProposeRequest) bool {
	log.Infof("propose hook on height %d is OK", req.Index)
	return true
}

func CommitShard(req *pb.CommitRequest) bool {
	log.Infof("commit hook on height %d is OK", req.Index)
	return true
}