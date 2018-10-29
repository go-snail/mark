package mark

import (
	"github.com/pkg/errors"
	"strconv"
	"sync"
	"time"
)

const (
	eventBufferSize = 1000 * 1024
)

type Feilds map[string]interface{}

type event struct {
	Value     map[string]interface{} `json:"value"`
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
				break
			} else {
				r.writer.write(evt)
			}
		default :
			break
		}
	}
}

func Mark(feild Feilds) {
	evt := reg.evtBuf.Get().(*event)
	evt.Value = feild
	evt.Value["timestamp"] = time.Now()
	evt.Value["mid"] = strconv.Itoa(int(GetAssetID()))
	select {
	case reg.eventBus <- evt:
	default:
		errors.New("metrics eventBus is full.")
	}
}
