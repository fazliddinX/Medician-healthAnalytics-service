package service

import (
	"context"
	"log/slog"
	pb "medical-service/generated/healthAnalytics"
	"medical-service/storage"
)

type Lifestyles struct {
	pb.UnimplementedLifestyleServiceServer
	storage.LifestyleStorage
	log *slog.Logger
}

func NewLifestyle(store storage.LifestyleStorage, log *slog.Logger) *Lifestyles {
	return &Lifestyles{log: log, LifestyleStorage: store}
}

func (l *Lifestyles) AddLifestyleData(ctx context.Context, in *pb.Lifestyle) (*pb.LifestyleResponse, error) {
	res, err := l.LifestyleStorage.AddLifestyleData(ctx, in)
	if err != nil {
		l.log.Error("Error adding lifestyle data", "error", err)
		return nil, err
	}

	return res, nil
}

func (l *Lifestyles) GetLifestyleData(ctx context.Context, in *pb.LifestyleID) (*pb.LifestyleResponse, error) {
	res, err := l.LifestyleStorage.GetLifestyleData(ctx, in)
	if err != nil {
		l.log.Error("Error getting lifestyle data", "error", err)
		return nil, err
	}

	return res, nil
}

func (l *Lifestyles) GetAllLifestyleData(ctx context.Context, in *pb.LifestyleFilter) (*pb.AllLifestyles, error) {
	res, err := l.LifestyleStorage.GetAllLifestyleData(ctx, in)
	if err != nil {
		l.log.Error("Error getting lifestyle data", "error", err)
		return nil, err
	}

	return res, nil
}

func (l *Lifestyles) UpdateLifestyleData(ctx context.Context, in *pb.UpdateLifestyle) (*pb.LifestyleResponse, error) {
	res, err := l.LifestyleStorage.UpdateLifestyleData(ctx, in)
	if err != nil {
		l.log.Error("Error updating lifestyle data", "error", err)
		return nil, err
	}

	return res, nil
}

func (l *Lifestyles) DeleteLifestyleData(ctx context.Context, in *pb.LifestyleID) (*pb.Message, error) {
	res, err := l.LifestyleStorage.DeleteLifestyleData(ctx, in)
	if err != nil {
		l.log.Error("Error deleting lifestyle data", "error", err)
		return nil, err
	}

	return res, nil
}
