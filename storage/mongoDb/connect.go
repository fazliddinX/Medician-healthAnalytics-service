package mongoDb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"medical-service/pkg/config"
)

func ConnectMongoDB(cfg config.Config) (*mongo.Database, error) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(cfg.MONGO_URL))

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database("UserProgress"), nil
}

//func TestConnectMongoDB() (*mongo.Database, error) {
//	client, err := mongo.Connect(context.Background(),
//		options.Client().ApplyURI("mongodb://localhost:27017"))
//
//	if err != nil {
//		return nil, err
//	}
//
//	err = client.Ping(context.TODO(), nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return client.Database("UserProgress"), nil
//}
