/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 8/27/21$ 7:33 PM$
 **/
package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/cache"
	"time"
)

type Process struct {
	NodeCache *cache.CacheShard
	ProcessEndCh chan int
}

func (p *Process)NewProcessor(c *cache.CacheShard){
	p.NodeCache = c
	p.ProcessEndCh = make(chan int)
}

//func (p *Process) SetCh(i int){
//	p.ProcessStartCh <- i
//}
//
func (p *Process) CloseCh(){
	close(p.ProcessEndCh)
}

func (p *Process) loop(){
	// 处理cache中的消息
	for{
		// 每2秒获取一次所有未处理的交易
		report := time.NewTicker(1 * time.Second)
		defer report.Stop()
		select {
		case  <-report.C:
			//fmt.Println("正在处理交易")
			//log.Info("正在处理交易")
			p.NodeCache.ProposeTxHandle()

			//msgs := p.NodeCache.GetAllUnHandled()
			//if len(msgs)>0{
			//	for _, msg := range msgs{
			//		msg.Ch <- 2
			//	}
			//}
		case <-p.ProcessEndCh:
			log.Info("收到结束信号")
			//fmt.Println("收到结束信号")
			return
		}
	}
}