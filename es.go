package mark

import (
	"context"
	log "github.com/cihub/seelog"
	"github.com/olivere/elastic"
	"time"
)

const (
	Type  = "log"
)

type ESClient struct {
	URL    string //es地址
	Scheme string //监视时使用的协议，默认是http
	Index  string //索引
	Es     *elastic.Client
}

func NewESClient(index, url, scheme string) *ESClient {
	log.Infof("index %s url %s scheme %s", index, url, scheme)
	ctx := context.Background()
	var err error
	if scheme == "" {
		scheme = elastic.DefaultScheme
	}
	if url == "" {
		url = elastic.DefaultURL
	}
	if index == "" {
		index = Index
	}
	esclient, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetScheme(scheme),
		elastic.SetHealthcheck(true), //true时, 设置健康检查
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetHealthcheckTimeout(1*time.Second),
		elastic.SetHealthcheckTimeoutStartup(2*time.Second),
		elastic.SetSniff(false), //true的时候,设置监测interval:SetSnifferInterval,SetSnifferTimeoutStartup,SetSnifferTimeout
		elastic.SetSendGetBodyAs("GET"),
	)
	if err != nil {
		log.Errorf("创建es客户端失败:", err)
		return nil
	}
	info, code, err := esclient.Ping(url).Do(ctx)
	if err != nil {
		log.Error("ping es服务失败:", err)
		return nil
	}
	log.Infof("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	return &ESClient{url, scheme, index, esclient}
}

func (es *ESClient) write(event *event) error {
	add, err := es.Es.Index().Index(Index).Type(Type).BodyJson(event).Do(context.Background())
	if err != nil {
		log.Error("写入es的日志失败:", err)
		return err
	}
	log.Infof("index (%s),type (%s) ,id (%s),result (%s)", add.Index, add.Type, add.Id, add.Result)
	return nil
}
