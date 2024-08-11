package storage

import (
	"context"
	pb "medical-service/generated/healthAnalytics"
)

type MedicalRecordStorage interface {
	AddMedicalRecord(ctx context.Context, in *pb.AddMedicalRecordRequest) (*pb.MedicalRecord, error)
	GetMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.MedicalRecord, error)
	UpdateMedicalRecord(ctx context.Context, in *pb.UpdateMedicalRecordReq) (*pb.MedicalRecord, error)
	DeleteMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.Message1, error)
	ListMedicalRecords(ctx context.Context, in *pb.MedicalRecordFilter) (*pb.ListMedicalRecord, error)
}

type LifestyleStorageImpl interface {
	AddLifestyleData(ctx context.Context, in *pb.Lifestyle) (*pb.LifestyleResponse, error)
	GetLifestyleData(ctx context.Context, in *pb.LifestyleID) (*pb.LifestyleResponse, error)
	GetAllLifestyleData(ctx context.Context, in *pb.LifestyleFilter) (*pb.AllLifestyles, error)
	UpdateLifestyleData(ctx context.Context, in *pb.UpdateLifestyle) (*pb.LifestyleResponse, error)
	DeleteLifestyleData(ctx context.Context, in *pb.LifestyleID) (*pb.Message, error)
}

type WearableDataStorage interface {
	AddWearableData(ctx context.Context, in *pb.WearableDate) (*pb.WearableDataResponse, error)
	GetWearableData(ctx context.Context, in *pb.WearableDataID) (*pb.WearableDataResponse, error)
	GetAllWearableData(ctx context.Context, in *pb.WearableDataFilter) (*pb.AllWearableData, error)
	UpdateWearableData(ctx context.Context, in *pb.UpdateWearableDate) (*pb.WearableDataResponse, error)
	DeleteWearableData(ctx context.Context, in *pb.WearableDataID) (*pb.Message2, error)
}

type HealthRecommendationsStorage interface {
	GenerateHealthRecommendations(ctx context.Context, in *pb.HealthRecommendationReq) (*pb.HealthRecommendation, error)
	GetHealthRecommendations(ctx context.Context, in *pb.HealthRecommendationID) (*pb.HealthRecommendation, error)
	GetAllHealthRecommendations(ctx context.Context, in *pb.UserID) (*pb.UserHealthRecommendation, error)

	GetRealtimeHealthMonitoring(ctx context.Context, in *pb.Void) (*pb.AllHealthRecommendations, error)
	GetDailyHealthSummary(ctx context.Context, in *pb.Void) (*pb.AllHealthRecommendations, error)
	GetWeeklyHealthSummary(ctx context.Context, in *pb.Void) (*pb.AllHealthRecommendations, error)
}
