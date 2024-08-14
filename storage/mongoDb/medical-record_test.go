package mongoDb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	pb "medical-service/generated/healthAnalytics"
	"testing"
	"time"
)

func TestAddMedicalRecord(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-records")
	st := NewMedicalRecordRepo(collection)

	req := &pb.AddMedicalRecordRequest{
		UserId:      "user-123",
		RecordType:  "diagnosis",
		RecordDate:  "2024-08-10",
		Description: "High blood pressure",
		DoctorId:    "doctor-456",
		Attachments: []string{"attachment1.jpg"},
	}

	res, err := st.AddMedicalRecord(context.Background(), req)
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

func TestGetMedicalRecord(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-records")
	st := NewMedicalRecordRepo(collection)

	objectID := primitive.NewObjectID()

	req := &pb.AddMedicalRecordRequest{
		UserId:      "user-123",
		RecordType:  "diagnosis",
		RecordDate:  "2024-08-10",
		Description: "High blood pressure",
		DoctorId:    "doctor-456",
		Attachments: []string{"attachment1.jpg"},
	}

	addedData := &pb.MedicalRecord{
		Id:          objectID.Hex(),
		UserId:      req.UserId,
		RecordType:  req.RecordType,
		RecordDate:  req.RecordDate,
		Description: req.Description,
		DoctorId:    req.DoctorId,
		Attachments: req.Attachments,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err = collection.InsertOne(context.Background(), addedData)
	if err != nil {
		t.Fatal(err)
	}
	a := objectID.Hex()
	log.Println(a)
	res, err := st.GetMedicalRecord(context.Background(), &pb.MedicalRecordID{Id: a})
	if err != nil {
		t.Fatal(err)
	}
	if res.UserId == "" {
		t.Fatal("user id should not be empty")
	}

	if res.Id != addedData.Id {
		log.Fatalf("expected id %s, got %s", addedData.Id, res.Id)
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"id": objectID})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMedicalRecord(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-records")
	st := NewMedicalRecordRepo(collection)

	req := &pb.AddMedicalRecordRequest{
		UserId:      "user-123",
		RecordType:  "diagnosis",
		RecordDate:  "2024-08-10",
		Description: "High blood pressure",
		DoctorId:    "doctor-456",
		Attachments: []string{"attachment1.jpg"},
	}

	addedData, err := st.AddMedicalRecord(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	updateReq := &pb.UpdateMedicalRecordReq{
		Id:         addedData.Id,
		RecordType: "treatment",
		RecordDate: "2024-08-11",
	}

	updatedRes, err := st.UpdateMedicalRecord(context.Background(), updateReq)
	if err != nil {
		t.Fatal(err)
	}

	if updatedRes.RecordType != updateReq.RecordType || updatedRes.RecordDate != updateReq.RecordDate {
		log.Fatalf("update failed: expected %v, got %v", updateReq, updatedRes)
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"id": updatedRes.Id})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListMedicalRecords(t *testing.T) {
	db, err := TestConnectMongoDB()
	if err != nil {
		t.Fatal(err)
	}

	collection := db.Collection("health-records")
	st := NewMedicalRecordRepo(collection)

	req1 := &pb.AddMedicalRecordRequest{
		UserId:      "user-123",
		RecordType:  "diagnosis",
		RecordDate:  "2024-08-10",
		Description: "High blood pressure",
		DoctorId:    "doctor-456",
		Attachments: []string{"attachment1.jpg"},
	}

	req2 := &pb.AddMedicalRecordRequest{
		UserId:      "user-124",
		RecordType:  "treatment",
		RecordDate:  "2024-08-11",
		Description: "Diabetes management",
		DoctorId:    "doctor-789",
		Attachments: []string{"attachment2.jpg"},
	}

	_, err = st.AddMedicalRecord(context.Background(), req1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.AddMedicalRecord(context.Background(), req2)
	if err != nil {
		t.Fatal(err)
	}

	filter := &pb.MedicalRecordFilter{
		Limit: 10,
	}

	res, err := st.ListMedicalRecords(context.Background(), filter)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.MedicalRecords) < 2 {
		t.Fatalf("expected at health-records, got %d", len(res.MedicalRecords))
	}

	_, err = collection.DeleteMany(context.Background(), bson.M{"userId": bson.M{"$in": []string{"user-123", "user-124"}}})
	if err != nil {
		t.Fatal(err)
	}
}
