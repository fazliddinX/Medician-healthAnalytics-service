package service

import (
	"context"
	"log/slog"
	pb "medical-service/generated/healthAnalytics"
	"medical-service/storage"
)

type Health struct {
	pb.UnimplementedHealthRecommendationsServiceServer
	storage.HealthRecommendationsStorage
	log *slog.Logger
}

func NewHealth(storage1 storage.HealthRecommendationsStorage, log *slog.Logger) *Health {
	return &Health{HealthRecommendationsStorage: storage1, log: log}
}

// ----------------- Health Recommendations --------------------
func (h *Health) GenerateHealthRecommendations(ctx context.Context, in *pb.HealthRecommendationReq) (*pb.HealthRecommendation, error) {
	res, err := h.HealthRecommendationsStorage.GenerateHealthRecommendations(ctx, in)
	if err != nil {
		h.log.Error("Failed to generate health recommendations", "error", err)
		return nil, err
	}

	return res, nil
}

func (h *Health) GetHealthRecommendations(ctx context.Context, in *pb.HealthRecommendationID) (*pb.HealthRecommendation, error) {
	res, err := h.HealthRecommendationsStorage.GetHealthRecommendations(ctx, in)
	if err != nil {
		h.log.Error("Failed to get health recommendations", "error", err)
		return nil, err
	}

	return res, nil
}

func (h *Health) GetAllHealthRecommendations(ctx context.Context, in *pb.UserID) (*pb.UserHealthRecommendation, error) {
	res, err := h.HealthRecommendationsStorage.GetAllHealthRecommendations(ctx, in)
	if err != nil {
		h.log.Error("Failed to get all health recommendations", "error", err)
		return nil, err
	}

	return res, nil
}

// ----------------------------------- Health Monitoring ---------------------------------------------------------

func (h *Health) GetRealtimeHealthMonitoring(ctx context.Context, in *pb.Void) (*pb.MonitoringRealTime, error) {
	res, err := h.HealthRecommendationsStorage.GetRealtimeHealthMonitoring(ctx, in)
	if err != nil {
		h.log.Error("Failed to get realtime health monitoring", "error", err)
		return nil, err
	}

	return res, nil
}

func (h *Health) GetDailyHealthSummary(ctx context.Context, in *pb.Void) (*pb.Monitoring, error) {
	res, err := h.HealthRecommendationsStorage.GetDailyHealthSummary(ctx, in)
	if err != nil {
		h.log.Error("Failed to get daily health summary", "error", err)
		return nil, err
	}

	return res, nil
}

func (h *Health) GetWeeklyHealthSummary(ctx context.Context, in *pb.Void) (*pb.Monitoring, error) {
	res, err := h.HealthRecommendationsStorage.GetWeeklyHealthSummary(ctx, in)
	if err != nil {
		h.log.Error("Failed to get weekly health summary", "error", err)
		return nil, err
	}

	return res, nil
}
