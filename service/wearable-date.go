package service

import (
	"context"
	"log/slog"
	pb "medical-service/generated/healthAnalytics"
	"medical-service/storage"
)

type WearableDate struct {
	pb.UnimplementedWearableDataServer
	storage.WearableDataStorage
	log *slog.Logger
}

func NewWearableDate(store storage.WearableDataStorage, log *slog.Logger) *WearableDate {
	return &WearableDate{log: log, WearableDataStorage: store}
}

func (w *WearableDate) AddWearableData(ctx context.Context, in *pb.WearableDate) (*pb.WearableDataResponse, error) {
	res, err := w.WearableDataStorage.AddWearableData(ctx, in)
	if err != nil {
		w.log.Error("Failed to add wearable data", "error", err)
		return nil, err
	}

	return res, nil
}

func (w *WearableDate) GetWearableData(ctx context.Context, in *pb.WearableDataID) (*pb.WearableDataResponse, error) {
	res, err := w.WearableDataStorage.GetWearableData(ctx, in)
	if err != nil {
		w.log.Error("Failed to get data", "error", err)
		return nil, err
	}

	return res, nil
}

func (w *WearableDate) GetAllWearableData(ctx context.Context, in *pb.WearableDataFilter) (*pb.AllWearableData, error) {
	res, err := w.WearableDataStorage.GetAllWearableData(ctx, in)
	if err != nil {
		w.log.Error("Failed to get data", "error", err)
		return nil, err
	}

	return res, nil
}

func (w *WearableDate) UpdateWearableData(ctx context.Context, in *pb.UpdateWearableDate) (*pb.WearableDataResponse, error) {
	res, err := w.WearableDataStorage.UpdateWearableData(ctx, in)
	if err != nil {
		w.log.Error("Failed to update data", "error", err)
		return nil, err
	}

	return res, nil
}

func (w *WearableDate) DeleteWearableData(ctx context.Context, in *pb.WearableDataID) (*pb.Message2, error) {
	res, err := w.WearableDataStorage.DeleteWearableData(ctx, in)
	if err != nil {
		w.log.Error("Failed to delete data", "error", err)
		return nil, err
	}

	return res, nil
}
