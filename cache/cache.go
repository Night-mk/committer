package cache

import (
	"fmt"
	"sync"
)

type msg struct {
	Key   string
	Value []byte
}

type Cache struct {
	store map[uint64]msg
	mu    sync.RWMutex
}

func New() *Cache {
	hashtable := make(map[uint64]msg)
	return &Cache{store: hashtable}
}

func (c *Cache) Set(index uint64, key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[index] = msg{key, value}
}

func (c *Cache) Get(index uint64) (string, []byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	message, ok := c.store[index]
	return message.Key, message.Value, ok
}

func (c *Cache) Delete(index uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, index)
}

/**
	dynamic sharding修改
	针对分片的缓存CacheShard
 */
type msgS struct {
	Source uint32
	Target uint32
	Keylist []byte
	Ch chan int
	Phase int // phase=1表示在propose阶段，phase=2表示在commit阶段
}

type CacheShard struct {
	store map[uint64]*msgS
	mu    sync.RWMutex
}

func NewCache() *CacheShard {
	hashtable := make(map[uint64]*msgS)
	return &CacheShard{store: hashtable}
}

// index标识事务的唯一编号
func (c *CacheShard) Set(index uint64, source uint32, target uint32, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[index] = &msgS{source, target, value, make(chan int), 0}
}

func (c *CacheShard) Get(index uint64) (uint32, uint32, []byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	message, ok := c.store[index]
	return message.Source, message.Target, message.Keylist, ok
}

func (c *CacheShard) Delete(index uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, index)
}

// 增加带有channel的设置
func (c *CacheShard) SetWithC(index uint64, source uint32, target uint32, value []byte, ch chan int, phase int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[index] = &msgS{source, target, value, ch, phase}
}

// 处理交易
func (c *CacheShard) SetPhase(index uint64, i int){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[index].Phase = i
}


// 获取全部没有返回过的交易
func (c *CacheShard) GetAllUnHandled() []*msgS{
	// 读取时加锁
	c.mu.RLock()
	defer c.mu.RUnlock()
	var unHandledMsgs []*msgS
	//fmt.Println("map总长度：",len(c.store))
	if len(c.store) != 0{
		for key, value := range c.store{
			// 判断channel值是否是1，是的话就将消息加入返回队列
			if value.Phase == 1{
				fmt.Println("获取未处理的交易编号： ", key)

				unHandledMsgs = append(unHandledMsgs, value)
				c.SetPhase(key, 2)
			}
		}
	}

	return unHandledMsgs
}

// 处理当前全部交易
func (c *CacheShard) ProposeTxHandle(){
	// 读取时加锁
	c.mu.RLock()
	defer c.mu.RUnlock()
	//fmt.Println("Handle loop map总长度：",len(c.store))
	if len(c.store) > 0{
		//count := 0
		for key, value := range c.store{
			// 判断phase值是否是1，是的话就处理消息
			//if value.Phase == 1{
			//	fmt.Println("正在处理的交易编号： ", key)
			//	value.Ch <- 2
			//	c.SetPhase(key, 2) // 读的时候不应该写相同的map
			//}
			// 并行方式
			go func(key uint64, value *msgS) { // 试试并行
				// 判断phase值是否是1，是的话就处理消息
				if value.Phase == 1{
					c.SetPhase(key, 2) // 读的时候不应该写相同的map
					//fmt.Println("正在处理的交易编号： ", key)
					value.Ch <- 2
				}
			}(key, value)
		}
	}

}