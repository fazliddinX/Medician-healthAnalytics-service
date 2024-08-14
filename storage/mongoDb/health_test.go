package mongoDb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	pb "medical-service/generated/healthAnalytics"
	"testing"
)

func TestGenerateHealthRecommendations(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-recommend")
	st := NewHealth(collection)

	req := &pb.HealthRecommendationReq{
		UserId:             "user-123",
		RecommendationType: "diet",
		Description:        "Eat more vegetables",
		Priority:           5,
	}

	res, err := st.GenerateHealthRecommendations(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	if res.Id == "" {
		t.Fatal("id should not be empty")
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"id": res.Id})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetHealthRecommendations(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-recommend")
	st := NewHealth(collection)

	req := &pb.HealthRecommendationReq{
		UserId:             "user-123",
		RecommendationType: "exercise",
		Description:        "Daily running for 30 minutes",
		Priority:           8,
	}

	addedData, err := st.GenerateHealthRecommendations(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := st.GetHealthRecommendations(context.Background(), &pb.HealthRecommendationID{Id: addedData.Id})
	if err != nil {
		t.Fatal(err)
	}

	if res.Id != addedData.Id {
		t.Fatalf("expected id %s, got %s", addedData.Id, res.Id)
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"id": res.Id})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllHealthRecommendations(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-recommend")
	st := NewHealth(collection)

	req1 := &pb.HealthRecommendationReq{
		UserId:             "user-123",
		RecommendationType: "diet",
		Description:        "Eat more vegetables",
		Priority:           5,
	}

	req2 := &pb.HealthRecommendationReq{
		UserId:             "user-124",
		RecommendationType: "exercise",
		Description:        "Daily running for 30 minutes",
		Priority:           8,
	}

	_, err = st.GenerateHealthRecommendations(context.Background(), req1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.GenerateHealthRecommendations(context.Background(), req2)
	if err != nil {
		t.Fatal(err)
	}

	filter := &pb.UserID{Id: "user-123"}
	res, err := st.GetAllHealthRecommendations(context.Background(), filter)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.HealthRecommends) == 0 {
		t.Fatalf("expected at least 1 record, got %d", len(res.HealthRecommends))
	}

	for _, rec := range res.HealthRecommends {
		_, err = collection.DeleteOne(context.Background(), bson.M{"id": rec.UserId})
		if err != nil {
			t.Fatal(err)
		}
	}
}
