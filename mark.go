package mark

import (
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	eventBufferSize = 1024 * 100
)

type Feilds map[string]interface{}

type event struct {
	Value  map[string]interface{} `value`
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
		peerpath := filepath.Join(p, "src/"+ cmdRoot)
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
		log.Info("not created mongo client")
		return nil
	default:
		return nil
	}
	return nil
}

func (r *reporter) eventLoop() {
	for {
		select {
		case evt,ok := <-r.eventBus:
			if !ok {
				log.Error("read eventBus chan failed")
			}else{
				if err := r.writer.write(evt); err != nil {
					log.Errorf("mark writer error : %v", err)
				}
			}
		}
	}
}


func Mark(feild Feilds) {
	//TODO 构造消息
	evt := reg.evtBuf.Get().(*event)
    evt.Value = feild
	//TODO 将消息写入到evtBuf中
	select {
	case reg.eventBus <- evt:
	default:
		log.Errorf("metrics eventBus is full.")
	}
}


