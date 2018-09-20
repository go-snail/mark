package mark

type Writer interface {
	write(*event) error
}

type ConfigI interface {
	Init() Writer
}

type MonGoConfig struct {
	Url   string

}

func (MC *MonGoConfig) Init() Writer {
	return  NewMongoClient()
}

type ESConfig struct {
	Url      string
	Index    string
	Scheme   string
	Username string
	Password string
}

func (EC *ESConfig) Init() Writer {
	if EC.Index == ""{
		EC.Index = "iot"
	}
	if EC.Scheme == ""{
		EC.Scheme = "http"
	}
	if EC.Url == ""{
		EC.Url = "http://127.0.0.1:9200"
	}
	return NewESClients(EC)
}

//type ESSource struct {
//	 ESConfig
//}
//func (ES *ESSource) GetType() int {
//	return 1
//}
//
//func (ES *ESSource) GetConfig() ESConfig {
//    return ESConfig{}
//}
//
//type MongoSource struct {
//	MonGoConfig
//}
//
//func (MS *MongoSource) GetType() int {
//	return 2
//}
//
//func (MS *MongoSource) GetConfig() ConfigI {
//	 return nil
//}
//
//type Source interface {
//	GetConfig() ConfigI
//}
