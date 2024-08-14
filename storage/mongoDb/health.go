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

type Health struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewHealth(collection *mongo.Collection, db *mongo.Database) storage.HealthRecommendationsStorage {
	return &Health{collection: collection, db: db}
}

func (h *Health) GenerateHealthRecommendations(ctx context.Context, in *pb.HealthRecommendationReq) (*pb.HealthRecommendation, error) {
	recommendation := &pb.HealthRecommendation{
		Id:                 primitive.NewObjectID().Hex(),
		UserId:             in.UserId,
		RecommendationType: in.RecommendationType,
		Description:        in.Description,
		Priority:           in.Priority,
		CreatedAt:          time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:          time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err := h.collection.InsertOne(ctx, recommendation)
	if err != nil {
		return nil, err
	}

	return recommendation, nil
}

func (h *Health) GetHealthRecommendations(ctx context.Context, in *pb.HealthRecommendationID) (*pb.HealthRecommendation, error) {
	var recommendation pb.HealthRecommendation

	err := h.collection.FindOne(ctx, bson.M{"id": in.Id}).Decode(&recommendation)
	if err != nil {
		return nil, err
	}

	return &recommendation, nil
}

func (h *Health) GetAllHealthRecommendations(ctx context.Context, in *pb.UserID) (*pb.UserHealthRecommendation, error) {
	cursor, err := h.collection.Find(ctx, bson.M{"userid": in.Id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recommendations []*pb.HealthRecommendationReq
	for cursor.Next(ctx) {
		var recommendation pb.HealthRecommendationReq
		if err := cursor.Decode(&recommendation); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, &recommendation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.UserHealthRecommendation{HealthRecommends: recommendations}, nil
}

// ------------------------------------ Monitoring ----------------------------------------------------------

func (h *Health) GetRealtimeHealthMonitoring(ctx context.Context, in *pb.Void) (*pb.MonitoringRealTime, error) {

	lifeCol := h.db.Collection("health-lifestyle")
	medRecordCol := h.db.Collection("health-records")
	wearableCol := h.db.Collection("health-wearable")

	filter := bson.M{}
	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var lifestyle pb.Lifestyles
	err := lifeCol.FindOne(ctx, filter, opts).Decode(&lifestyle)
	if err != nil {
		return nil, err
	}

	var medRecords pb.MedicalRecords
	err = medRecordCol.FindOne(ctx, filter, opts).Decode(&medRecords)
	if err != nil {
		return nil, err
	}

	var wearableDate pb.WearableDates
	err = wearableCol.FindOne(ctx, filter, opts).Decode(&wearableDate)
	if err != nil {
		return nil, err
	}

	res := &pb.MonitoringRealTime{
		Message:      "Last created Dates",
		MedRecords:   &medRecords,
		Lifestyle:    &lifestyle,
		WearableData: &wearableDate,
	}

	return res, nil
}

func (h *Health) GetDailyHealthSummary(ctx context.Context, in *pb.Void) (*pb.Monitoring, error) {
	lifeCol := h.db.Collection("health-lifestyle")
	medRecordCol := h.db.Collection("health-records")
	wearableCol := h.db.Collection("health-wearable")

	filter := bson.M{}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	lifestyleCursor, err := lifeCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer lifestyleCursor.Close(ctx)

	var lifestyles []*pb.Lifestyles
	if err = lifestyleCursor.All(ctx, &lifestyles); err != nil {
		return nil, err
	}

	medRecordCursor, err := medRecordCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer medRecordCursor.Close(ctx)

	var medicalRecords []*pb.MedicalRecords
	if err = medRecordCursor.All(ctx, &medicalRecords); err != nil {
		return nil, err
	}

	wearableCursor, err := wearableCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer wearableCursor.Close(ctx)

	var wearableData []*pb.WearableDates
	if err = wearableCursor.All(ctx, &wearableData); err != nil {
		return nil, err
	}

	res := &pb.Monitoring{
		Message:      "All created Dates",
		MedRecords:   medicalRecords,
		Lifestyle:    lifestyles,
		WearableData: wearableData,
	}

	return res, nil
}

func (h *Health) GetWeeklyHealthSummary(ctx context.Context, in *pb.Void) (*pb.Monitoring, error) {
	lifeCol := h.db.Collection("health-lifestyle")
	medRecordCol := h.db.Collection("health-records")
	wearableCol := h.db.Collection("health-wearable")

	filter := bson.M{}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	lifestyleCursor, err := lifeCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer lifestyleCursor.Close(ctx)

	var lifestyles []*pb.Lifestyles
	if err = lifestyleCursor.All(ctx, &lifestyles); err != nil {
		return nil, err
	}

	medRecordCursor, err := medRecordCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer medRecordCursor.Close(ctx)

	var medicalRecords []*pb.MedicalRecords
	if err = medRecordCursor.All(ctx, &medicalRecords); err != nil {
		return nil, err
	}

	wearableCursor, err := wearableCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer wearableCursor.Close(ctx)

	var wearableData []*pb.WearableDates
	if err = wearableCursor.All(ctx, &wearableData); err != nil {
		return nil, err
	}

	res := &pb.Monitoring{
		Message:      "All created Dates",
		MedRecords:   medicalRecords,
		Lifestyle:    lifestyles,
		WearableData: wearableData,
	}

	return res, nil
}
