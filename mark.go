package mark

import (
	"time"
	"sync"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
	"github.com/prometheus/common/log"
	"motan-go/core"
)

const (
	eventBufferSize = 1024 * 100
)



type event struct {
	Tid       int32     `json:"tid"`
	Pid       int32     `json:"pid"`
	Did       string    `json:"did"`
	Level     int32     `json:"level"` //日志类型
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
	Test      string    `json:"test"`
	Message   string    `json:"message"`
}

type reporter struct {
	stopping int32
	eventBus chan *event
	interval time.Duration
	writer  Writer
	evtBuf *sync.Pool
}

var (
	sinkDuration = time.Second * 5
	reg          = &reporter{
		stopping: 0,
		eventBus: make(chan *event, eventBufferSize),
		evtBuf:   &sync.Pool{New: func() interface{} { return new(event) }},
	}
)

func Run(cmdRoot string) error {
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		peerpath := filepath.Join(p, "src/"+cmdRoot)
		viper.AddConfigPath(peerpath)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return err
	}

	t := viper.GetString("mark.type")
	url := viper.GetString("mark.url")
	scheme := viper.GetString("mark.scheme")
	writer := getWriter(t, url, scheme)
	reg.writer = writer
	go reg.eventLoop()
	return nil
}

func getWriter(t, url, scheme string ) Writer {
	switch t {
	case "elasticSearch":
		return NewESClient(url, scheme)
	case "mongo":
	default:
		return nil

	}
	return nil
}

func (r *reporter)eventLoop() {
	for {
		select {
		case evt := <-r.eventBus:
			if err := r.writer.write(evt); err != nil {
				log.Errorf("mark writer error : %v", err)
			}

		}
	}
}


func Mark(map[string]interface{}) {
	//TODO 构造消息
	evt := reg.evtBuf.Get().(*event)


	//TODO 将消息写入到evtBuf中
	select {
	case reg.eventBus <- evt:
	default:
		log.Errorf("metrics eventBus is full.")
	}
}


