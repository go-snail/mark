package mark

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
		if err := Run(&ESConfig{Url:"http://192.168.1.51:9200",Index:"test"});err != nil {
			fmt.Println(err)
		}
		mes := map[string]interface{}{
		"tid":       100001,
		"did":       "123abc",
		"user":      "lvxx@radacat",
		"message":   "测试",
		"pid":       123,
		"timestamp": time.Now(),
		}
		//Mark(mes)
	   for i := 0;i<10000;i++{
	   		Mark(mes)
	   }
		time.Sleep(10 *  time.Second)
}

func Benchmark_mark(b *testing.B) {
	b.StopTimer()
	if err := Run(&ESConfig{Url:"http://192.168.1.6:9200",Index:"test"}); err != nil {
		fmt.Println(err)
		return
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mes := map[string]interface{}{
			"tid":       100001,
			"did":       "123abc",
			"user":      "lvxx@radacat",
			"message":   "测试",
			"pid":       123,
			"timestamp": time.Now(),
		}
		Mark(mes)
	}
}
