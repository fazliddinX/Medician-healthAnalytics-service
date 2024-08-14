package mongoDb

import (
	"context"
	"log"
	"testing"
	"time"

	pb "medical-service/generated/healthAnalytics"
)

func TestGetWearableData(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-wearable")
	st := NewWearableData(collection)

	// Вставляем тестовые данные
	req := &pb.WearableDate{
		UserId:       "Hello world",
		DeviceType:   "bittanarse",
		DataType:     "aaaaaa",
		DataValue:    "bbbbbbbbbb",
		RecordedTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	addedData, err := st.AddWearableData(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	res, err := st.GetWearableData(context.Background(), &pb.WearableDataID{Id: addedData.Id})
	if err != nil {
		t.Fatal(err)
	}

	if res.Id != addedData.Id {
		log.Fatalf("expected id %s, got %s", addedData.Id, res.Id)
	}

	_, err = st.DeleteWearableData(context.Background(), &pb.WearableDataID{Id: addedData.Id})
}

func TestUpdateWearableData(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-wearable")
	st := NewWearableData(collection)

	req := &pb.WearableDate{
		UserId:       "user-123",
		DeviceType:   "bittanarse",
		DataType:     "krc",
		DataValue:    "Usha narsa",
		RecordedTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	addedData, err := st.AddWearableData(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	updateReq := &pb.UpdateWearableDate{
		Id:         addedData.Id,
		DeviceType: "updated-device",
		DataType:   "updated-type",
		DataValue:  "updated-value",
	}

	updatedRes, err := st.UpdateWearableData(context.Background(), updateReq)
	if err != nil {
		t.Fatal(err)
	}

	if updatedRes.DeviceType != updateReq.DeviceType || updatedRes.DataType != updateReq.DataType || updatedRes.DataValue != updateReq.DataValue {
		log.Fatalf("update failed: expected %v, got %v", updateReq, updatedRes)
	}

	_, err = st.DeleteWearableData(context.Background(), &pb.WearableDataID{Id: updatedRes.Id})
	if err != nil {
		t.Fatal(err)
	}
}
