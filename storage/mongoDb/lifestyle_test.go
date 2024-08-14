package mongoDb

import (
	"context"
	"testing"
	"time"

	pb "medical-service/generated/healthAnalytics"
)

func TestAddLifestyleData(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-lifestyle")
	st := NewLifestyle(collection)

	req := &pb.Lifestyle{
		UserId:       "user-123",
		DataType:     "steps",
		DataValue:    "10000",
		RecordedData: time.Now().Format("2006-01-02 15:04:05"),
	}

	res, err := st.AddLifestyleData(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	if res.Id == "" {
		t.Fatal("id should not be empty")
	}

	_, err = st.DeleteLifestyleData(context.Background(), &pb.LifestyleID{Id: res.Id})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLifestyleData(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-lifestyle")
	st := NewLifestyle(collection)

	req := &pb.Lifestyle{
		UserId:       "user-123",
		DataType:     "steps",
		DataValue:    "10000",
		RecordedData: time.Now().Format("2006-01-02 15:04:05"),
	}

	addedData, err := st.AddLifestyleData(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := st.GetLifestyleData(context.Background(), &pb.LifestyleID{Id: addedData.Id})
	if err != nil {
		t.Fatal(err)
	}

	if res.Id != addedData.Id {
		t.Fatalf("expected id %s, got %s", addedData.Id, res.Id)
	}

	_, err = st.DeleteLifestyleData(context.Background(), &pb.LifestyleID{Id: res.Id})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllLifestyleData(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-lifestyle")
	st := NewLifestyle(collection)

	req1 := &pb.Lifestyle{
		UserId:       "user-123",
		DataType:     "steps",
		DataValue:    "10000",
		RecordedData: time.Now().Format("2006-01-02 15:04:05"),
	}

	req2 := &pb.Lifestyle{
		UserId:       "user-124",
		DataType:     "calories",
		DataValue:    "2500",
		RecordedData: time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err = st.AddLifestyleData(context.Background(), req1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.AddLifestyleData(context.Background(), req2)
	if err != nil {
		t.Fatal(err)
	}

	filter := &pb.LifestyleFilter{Limit: 2, Offset: 0}
	res, err := st.GetAllLifestyleData(context.Background(), filter)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Lifestyles) != 2 {
		t.Fatalf("expected 2 records, got %d", len(res.Lifestyles))
	}

	_, err = st.DeleteLifestyleData(context.Background(), &pb.LifestyleID{Id: res.Lifestyles[0].Id})
	if err != nil {
		t.Fatal(err)
	}

	_, err = st.DeleteLifestyleData(context.Background(), &pb.LifestyleID{Id: res.Lifestyles[1].Id})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateLifestyleData(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-lifestyle")
	st := NewLifestyle(collection)

	req := &pb.Lifestyle{
		UserId:       "user-123",
		DataType:     "steps",
		DataValue:    "10000",
		RecordedData: time.Now().Format("2006-01-02 15:04:05"),
	}

	addedData, err := st.AddLifestyleData(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	updateReq := &pb.UpdateLifestyle{
		Id:        addedData.Id,
		DataType:  "updated-steps",
		DataValue: "15000",
	}

	updatedRes, err := st.UpdateLifestyleData(context.Background(), updateReq)
	if err != nil {
		t.Fatal(err)
	}

	if updatedRes.DataType != updateReq.DataType || updatedRes.DataValue != updateReq.DataValue {
		t.Fatalf("update failed: expected %v, got %v", updateReq, updatedRes)
	}

	_, err = st.DeleteLifestyleData(context.Background(), &pb.LifestyleID{Id: updatedRes.Id})
	if err != nil {
		t.Fatal(err)
	}
}
