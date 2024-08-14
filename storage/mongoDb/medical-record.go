package mongoDb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	pb "medical-service/generated/healthAnalytics"
	"medical-service/storage"
	"time"
)

type MedicalRecordRepo struct {
	col *mongo.Collection
}

func NewMedicalRecordRepo(col *mongo.Collection) storage.MedicalRecordStorage {
	return &MedicalRecordRepo{col: col}
}

func (m *MedicalRecordRepo) AddMedicalRecord(ctx context.Context, in *pb.AddMedicalRecordRequest) (*pb.MedicalRecord, error) {
	res := &pb.MedicalRecord{
		Id:          in.Id,
		UserId:      in.UserId,
		RecordType:  in.RecordType,
		RecordDate:  in.RecordDate,
		Description: in.Description,
		DoctorId:    in.DoctorId,
		Attachments: in.Attachments,
		CreatedAt:   time.Now().Format("2006/01/02"),
		UpdatedAt:   time.Now().Format("2006/01/02"),
	}

	_, err := m.col.InsertOne(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *MedicalRecordRepo) GetMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.MedicalRecord, error) {
	var wearable pb.MedicalRecord
	err := m.col.FindOne(ctx, bson.M{"id": in.Id}).Decode(&wearable)
	if err != nil {
		return nil, err
	}

	return &wearable, nil
}

func (m *MedicalRecordRepo) UpdateMedicalRecord(ctx context.Context, in *pb.UpdateMedicalRecordReq) (*pb.MedicalRecord, error) {
	update := bson.M{
		"$set": bson.M{
			"recordType": in.RecordType,
			"recordDate": in.RecordDate,
			"updatedAt":  time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	res := &pb.MedicalRecord{}
	err := m.col.FindOneAndUpdate(ctx, bson.M{"id": in.Id}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("medical record with ID %s not found", in.Id)
		}
		return nil, err
	}

	return res, nil
}

func (m *MedicalRecordRepo) DeleteMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.Message1, error) {
	_, err := m.col.DeleteOne(ctx, bson.M{"id": in.Id})
	if err != nil {
		return nil, err
	}

	return &pb.Message1{Message: "medical record deleted"}, nil
}

func (m *MedicalRecordRepo) ListMedicalRecords(ctx context.Context, in *pb.MedicalRecordFilter) (*pb.ListMedicalRecord, error) {
	filter := bson.M{}
	if in.Description != "" {
		filter["description"] = in.Description
	}
	if in.DoctorId != "" {
		filter["doctorId"] = in.DoctorId
	}

	if in.Limit == 0 {
		in.Limit = 10
	}

	opts := options.Find().SetLimit(in.Limit).SetSkip(in.Offset)

	cursor, err := m.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []*pb.MedicalRecord
	for cursor.Next(ctx) {
		var record pb.MedicalRecord
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return &pb.ListMedicalRecord{MedicalRecords: records}, nil
}
