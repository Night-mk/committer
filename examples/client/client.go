package main

import (
	"context"
	"fmt"
	"github.com/openzipkin/zipkin-go"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/peer"
	pb "github.com/vadiminshakov/committer/proto"
	"github.com/vadiminshakov/committer/trace"
	"sync"
	"time"
)

const addr = "localhost:3000"

// 压力测试{
func pressureTest(addr string, tracer *zipkin.Tracer) {
	// create a set of clients
	const clientsNum = 100
	var clients [clientsNum]*peer.CommitClient
	var err [clientsNum]error
	// 创建多个client instance
	for i:=0; i<clientsNum; i++{
		clients[i], err[i] = peer.New(addr, tracer)
		if err[i] != nil {
			panic(err[i])
		}
	}
	time.Sleep(3* time.Second)
}

// 多个clients并发提交10组数据，测试总时间消耗
func multiPutConcurrent(){
	msgNum := 10
	respDone := make(chan *pb.Response, msgNum)
	errDone := make(chan error, msgNum)
	//doneSig := make(chan int)
	var testtable = map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
		"key4": []byte("value4"),
		"key5": []byte("value5"),
		"key6": []byte("value6"),
		"key7": []byte("value7"),
		"key8": []byte("value8"),
		"key9": []byte("value9"),
		"key10": []byte("value10"),
	}

	const clientsNum = 1000
	var clients [clientsNum]*peer.CommitClient
	var errs [clientsNum]error
	tracer, err := trace.Tracer("client", addr)
	if err != nil {
		panic(err)
	}
	// 创建多个client instance
	for i:=0; i<clientsNum; i++{
		clients[i], errs[i] = peer.New(addr, tracer)
		if errs[i] != nil {
			panic(errs[i])
		}
	}
	time.Sleep(3* time.Second) // 给足够的时间创建client
	log.Info("create clients end")
	var go_sync sync.WaitGroup
	// 启动多个client并发发送10条消息
	for _, cli := range clients{
		go_sync.Add(1)
		go func() {
			// 添加panic错误处理
			defer func() {
				if e := recover(); e != nil {
					fmt.Println("panic err: ", e)
				}
			}()
			defer go_sync.Done()
			for key, value := range testtable{
				resp, err := cli.Put(context.Background(), key, value)
				respDone <- resp
				errDone <- err
				fmt.Println("execute: ",key)
			}
		}()
	}
	for i:=0; i<len(clients)*len(testtable); i++{
		go_sync.Add(1)
		go func() {
			// 添加panic错误处理
			defer func() {
				if e := recover(); e != nil {
					fmt.Println("panic err: ", e)
				}
			}()
			defer go_sync.Done()
			select {
			case e := <-errDone:
				if e != nil {
					panic(err)
				}
			case r := <-respDone:
				fmt.Println("response: ", r.Type)
				if r.Type != pb.Type_ACK {
					panic("msg is not acknowledged")
				}
			}
		}()
	}
	go_sync.Wait()
}

// 超过1000个客户端时zipkin连接存在问题
// 测试元素
func multiPut(){
	var testtable = map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
		"key4": []byte("value4"),
		"key5": []byte("value5"),
		"key6": []byte("value6"),
		"key7": []byte("value7"),
		"key8": []byte("value8"),
		"key9": []byte("value9"),
		"key10": []byte("value10"),
	}
	const clientsNum = 1
	var clients [clientsNum]*peer.CommitClient
	var errs [clientsNum]error
	tracer, err := trace.Tracer("client", addr)
	if err != nil {
		panic(err)
	}
	// 创建多个client instance
	for i:=0; i<clientsNum; i++{
		clients[i], errs[i] = peer.New(addr, tracer)
		//clients[i], errs[i] = peer.New(addr, nil)
		if errs[i] != nil {
			panic(errs[i])
		}
	}
	// 循环多个client发送多条消息
	for _, cli := range clients{
		for key, value := range testtable{
			resp, err := cli.Put(context.Background(), key, value)
			if err != nil {
				panic(err)
			}
			if resp.Type != pb.Type_ACK {
				panic("msg is not acknowledged")
			}
			fmt.Println("execute: ",key)
		}
	}
}

func main()  {
	// 测试多个client并发put多条数据
	//multiPutConcurrent()
	// 测试多个client循环put多条数据
	multiPut()
}

//func main() {
//	tracer, err := trace.Tracer("client", addr)
//	if err != nil {
//		panic(err)
//	}
//	cli, err := peer.New(addr, tracer)
//	if err != nil {
//		panic(err)
//	}
//	resp, err := cli.Put(context.Background(), "key3", []byte("1111"))
//	if err != nil {
//		panic(err)
//	}
//	if resp.Type != pb.Type_ACK {
//		panic("msg is not acknowledged")
//	}
//
//	// read committed keys
//	//key, err := cli.Get(context.Background(), "key3")
//	//if err != nil {
//	//	panic(err)
//	//}
//	//fmt.Println(string(key.Value))
//}
