package mark

import (
	log "github.com/cihub/seelog"
	"sync"
	"time"
)

const (
	eventBufferSize = 1024 * 1000
)

type Feilds map[string]interface{}

type event struct {
	Value     map[string]interface{} `value`
	TimeStamp time.Time              `timestamp`
}

type reporter struct {
	stopping int32
	eventBus chan *event
	interval time.Duration
	writer   Writer
	evtBuf   *sync.Pool
}

var (
	sinkDuration = time.Second * 5
	reg          = &reporter{
		stopping: 0,
		eventBus: make(chan *event, eventBufferSize),
		evtBuf:   &sync.Pool{New: func() interface{} { return new(event) }},
	}
)

type Config struct {
	Url    string
	Scheme string
	Index  string
}

const (
	Url    = "http://127.0.0.1:9200"
	Scheme = "http"
	Index  = "iot"
)

func Run(i ConfigI) error {
	reg.writer = i.Init()
	go reg.eventLoop()
	return nil
}

func (r *reporter) eventLoop() {
	for {
		select {
		case evt, ok := <-r.eventBus:
			if !ok {
				log.Error("read eventBus chan failed")
				break
			} else {
				r.writer.write(evt)
			}
		}
	}
}

func Mark(feild Feilds) {
	evt := reg.evtBuf.Get().(*event)
	evt.Value = feild
	evt.TimeStamp = time.Now()
	select {
	case reg.eventBus <- evt:
	default:
		log.Errorf("metrics eventBus is full.")
	}
}
