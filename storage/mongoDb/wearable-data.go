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

type WearableData struct {
	collection *mongo.Collection
}

func NewWearableData(collection *mongo.Collection) storage.WearableDataStorage {
	return &WearableData{collection}
}

func (w *WearableData) AddWearableData(ctx context.Context, in *pb.WearableDate) (*pb.WearableDataResponse, error) {
	wearable := &pb.WearableDataResponse{
		Id:           primitive.NewObjectID().Hex(),
		UserId:       in.UserId,
		DeviceType:   in.DeviceType,
		DataType:     in.DataType,
		DataValue:    in.DataValue,
		RecordedTime: in.RecordedTime,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err := w.collection.InsertOne(ctx, wearable)
	if err != nil {
		return nil, err
	}

	return wearable, nil
}

func (w *WearableData) GetWearableData(ctx context.Context, in *pb.WearableDataID) (*pb.WearableDataResponse, error) {
	var wearable pb.WearableDataResponse
	err := w.collection.FindOne(ctx, bson.M{"id": in.Id}).Decode(&wearable)
	if err != nil {
		return nil, err
	}

	return &wearable, nil
}

func (w *WearableData) GetAllWearableData(ctx context.Context, in *pb.WearableDataFilter) (*pb.AllWearableData, error) {
	if in.Limit == 0 {
		in.Limit = 10
	}
	opts := options.Find().SetLimit(in.Limit).SetSkip(in.Offset)

	filter := bson.M{}
	if in.DeviceType != "" {
		filter["deviceType"] = in.DeviceType
	}

	cursor, err := w.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var wearableData []*pb.WearableDataResponse
	for cursor.Next(ctx) {
		var wearable pb.WearableDataResponse
		if err := cursor.Decode(&wearable); err != nil {
			return nil, err
		}
		wearableData = append(wearableData, &wearable)
	}

	return &pb.AllWearableData{WearableData: wearableData}, nil
}

func (w *WearableData) UpdateWearableData(ctx context.Context, in *pb.UpdateWearableDate) (*pb.WearableDataResponse, error) {
	update := bson.M{
		"$set": bson.M{
			"deviceType": in.DeviceType,
			"dataType":   in.DataType,
			"dataValue":  in.DataValue,
			"updatedAt":  time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	var updatedWearable pb.WearableDataResponse
	err := w.collection.FindOneAndUpdate(ctx, bson.M{"id": in.Id}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedWearable)
	if err != nil {
		return nil, err
	}

	return &updatedWearable, nil
}

func (w *WearableData) DeleteWearableData(ctx context.Context, in *pb.WearableDataID) (*pb.Message2, error) {
	_, err := w.collection.DeleteOne(ctx, bson.M{"id": in.Id})
	if err != nil {
		return nil, err
	}

	return &pb.Message2{Message: "Wearable data deleted successfully"}, nil
}
