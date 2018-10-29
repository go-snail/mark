package mark

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"time"
)

const (
	Type = "log"
)

//elastic search
type ESClient struct {
	Index string
	Es    *elastic.Client
}

func NewESClients(es *ESConfig) *ESClient {
	if es.Open {
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
			return nil
		}
		_, _, err = esclient.Ping(es.Url).Do(ctx)
		if err != nil {
			return nil
		}
		return &ESClient{es.Index, esclient}
	}
	return &ESClient{}
}

func (es *ESClient) write(event *event) error {
	if es.Es == nil {
		return errors.New("es client not connection")
	}
	_, err := es.Es.Index().Index(es.Index).Type(Type).BodyJson(event).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
