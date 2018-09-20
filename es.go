package mark

import (
	"context"
	log "github.com/cihub/seelog"
	"github.com/olivere/elastic"
	"time"
)

const (
	Type = "log"
)

type ESClients struct {
	Index string
	Es    *elastic.Client
}

func NewESClients(es *ESConfig) *ESClients {
	log.Infof("es config %s:", es)
	ctx := context.Background()
	var err error
	esclient, err := elastic.NewClient(
		elastic.SetURL(es.Url),
		elastic.SetScheme(es.Scheme),
		elastic.SetHealthcheck(true), //true时, 设置健康检查
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetHealthcheckTimeout(1*time.Second),
		elastic.SetHealthcheckTimeoutStartup(2*time.Second),
		elastic.SetSniff(false), //true的时候,设置监测interval:SetSnifferInterval,SetSnifferTimeoutStartup,SetSnifferTimeout
		elastic.SetSendGetBodyAs("GET"),
		elastic.SetBasicAuth(es.Username, es.Password),
	)
	if err != nil {
		log.Errorf("创建es客户端失败:", err)
		return nil
	}
	info, code, err := esclient.Ping(es.Url).Do(ctx)
	if err != nil {
		log.Error("ping es服务失败:", err)
		return nil
	}
	log.Infof("es returned with code %d and version %s", code, info.Version.Number)
	return &ESClients{es.Index, esclient}
}

func (es *ESClients) write(event *event) error {
	add, err := es.Es.Index().Index(es.Index).Type(Type).BodyJson(event).Do(context.Background())
	if err != nil {
		log.Error("写入es的日志失败:", err)
		return err
	}
	log.Infof("index (%s),type (%s) ,id (%s),result (%s)", add.Index, add.Type, add.Id, add.Result)
	return nil
}
