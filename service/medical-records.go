package service

import (
	"context"
	"log/slog"
	pb "medical-service/generated/healthAnalytics"
	"medical-service/storage"
)

type MedicalRecords struct {
	pb.UnimplementedMedicalRecordsServiceServer
	storage.MedicalRecordStorage
	log *slog.Logger
}

func NewMedicalRecords(storage storage.MedicalRecordStorage, log *slog.Logger) *MedicalRecords {
	return &MedicalRecords{log: log, MedicalRecordStorage: storage}
}

func (m *MedicalRecords) AddMedicalRecord(ctx context.Context, in *pb.AddMedicalRecordRequest) (*pb.MedicalRecord, error) {
	res, err := m.MedicalRecordStorage.AddMedicalRecord(ctx, in)
	if err != nil {
		m.log.Error("Failed to add medical record", "error", err)
		return nil, err
	}

	return res, nil
}

func (m *MedicalRecords) GetMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.MedicalRecord, error) {
	res, err := m.MedicalRecordStorage.GetMedicalRecord(ctx, in)
	if err != nil {
		m.log.Error("Failed to get medical record", "error", err)
		return nil, err
	}

	return res, nil
}

func (m *MedicalRecords) UpdateMedicalRecord(ctx context.Context, in *pb.UpdateMedicalRecordReq) (*pb.MedicalRecord, error) {
	res, err := m.MedicalRecordStorage.UpdateMedicalRecord(ctx, in)
	if err != nil {
		m.log.Error("Failed to update medical record", "error", err)
		return nil, err
	}

	return res, nil
}

func (m *MedicalRecords) DeleteMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.Message1, error) {
	res, err := m.MedicalRecordStorage.DeleteMedicalRecord(ctx, in)
	if err != nil {
		m.log.Error("Failed to delete medical record", "error", err)
		return nil, err
	}

	return res, nil
}

func (m *MedicalRecords) ListMedicalRecords(ctx context.Context, in *pb.MedicalRecordFilter) (*pb.ListMedicalRecord, error) {
	res, err := m.MedicalRecordStorage.ListMedicalRecords(ctx, in)
	if err != nil {
		m.log.Error("Failed to list medical records", "error", err)
		return nil, err
	}

	return res, nil
}
