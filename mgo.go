package mark

//暂未实现

type MongoClient struct{
	MgClient  string
}

func NewMongoClient() *MongoClient {
	return &MongoClient{MgClient:""}
}

func (MC *MongoClient)  write(event *event)  error{
	return nil
}