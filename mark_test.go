package mark

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
		if err := Run(Config{Url:"",Scheme:"",Index:""});err != nil {
			fmt.Println(err)
		}
		mes := make(map[string]interface{})
		mes["123"] = "qwe"
		mes["456"] = "asdf"
		Mark(mes)

}

func Benchmark_mark(b *testing.B) {
	b.StopTimer()
	if err := Run(Config{Url:"http://192.168.1.6:9200",Scheme:"http"}); err != nil {
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
