/**
 * @Author Night-mk
 * @Description //TODO
 * @Date 8/29/21$ 12:06 AM$
 **/
package config

import (
	"fmt"
	"testing"
)

func test(){
	config := Get()
	for _, item := range config.Followers{
		fmt.Println(item)
	}
	//config := GetConfig()
	//for key, item := range config.Followers{
	//	for _, item1 := range item{
	//		fmt.Println(key)
	//		fmt.Println(item1)
	//	}
	//}
}

func TestGetConfig(t *testing.T) {
	test()
}