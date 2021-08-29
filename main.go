/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 8/28/21$ 10:41 PM$
 **/
package main

import (
	"github.com/vadiminshakov/committer/config"
	"github.com/vadiminshakov/committer/hooks"
	"github.com/vadiminshakov/committer/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	conf := config.GetConfig()
	hooks, err := hooks.GetHookF()
	if err != nil {
		panic(err)
	}
	s, err := server.NewCommitServerShard(conf, hooks...)
	if err != nil {
		panic(err)
	}

	// Run函数启动一个non-blocking GRPC server
	s.Run(server.WhiteListCheckerShard)
	<-ch // 如果一直没有系统信号，就一直等待？
	s.Stop()
}