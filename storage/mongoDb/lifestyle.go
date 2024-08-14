package mongoDb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	pb "medical-service/generated/healthAnalytics"
	"medical-service/storage"
	"time"
)

type Lifestyle struct {
	collection *mongo.Collection
}

func NewLifestyle(collection *mongo.Collection) storage.LifestyleStorage {
	return &Lifestyle{collection}
}

func (l *Lifestyle) AddLifestyleData(ctx context.Context, in *pb.Lifestyle) (*pb.LifestyleResponse, error) {
	lifestyle := &pb.LifestyleResponse{
		Id:           primitive.NewObjectID().Hex(),
		UserId:       in.UserId,
		DataType:     in.DataType,
		DataValue:    in.DataValue,
		RecordedData: in.RecordedData,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err := l.collection.InsertOne(ctx, lifestyle)
	if err != nil {
		return nil, err
	}

	return lifestyle, nil
}

func (l *Lifestyle) GetLifestyleData(ctx context.Context, in *pb.LifestyleID) (*pb.LifestyleResponse, error) {
	var lifestyle pb.LifestyleResponse
	err := l.collection.FindOne(ctx, bson.M{"id": in.Id}).Decode(&lifestyle)
	if err != nil {
		return nil, err
	}

	return &lifestyle, nil
}

func (l *Lifestyle) GetAllLifestyleData(ctx context.Context, in *pb.LifestyleFilter) (*pb.AllLifestyles, error) {
	if in.Limit == 0 {
		in.Limit = 10
	}

	opts := options.Find().SetLimit(in.Limit).SetSkip(in.Offset)

	cursor, err := l.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var lifestyles []*pb.LifestyleResponse
	for cursor.Next(ctx) {
		var lifestyle pb.LifestyleResponse
		if err := cursor.Decode(&lifestyle); err != nil {
			return nil, err
		}
		lifestyles = append(lifestyles, &lifestyle)
	}

	return &pb.AllLifestyles{Lifestyles: lifestyles}, nil
}

func (l *Lifestyle) UpdateLifestyleData(ctx context.Context, in *pb.UpdateLifestyle) (*pb.LifestyleResponse, error) {
	update := bson.M{
		"$set": bson.M{
			"dataType":  in.DataType,
			"dataValue": in.DataValue,
			"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	var updatedLifestyle pb.LifestyleResponse
	err := l.collection.FindOneAndUpdate(ctx, bson.M{"id": in.Id}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedLifestyle)
	if err != nil {
		return nil, err
	}

	return &updatedLifestyle, nil
}

func (l *Lifestyle) DeleteLifestyleData(ctx context.Context, in *pb.LifestyleID) (*pb.Message, error) {
	_, err := l.collection.DeleteOne(ctx, bson.M{"id": in.Id})
	if err != nil {
		return nil, err
	}

	return &pb.Message{Message: "Lifestyle data deleted successfully"}, nil
}
