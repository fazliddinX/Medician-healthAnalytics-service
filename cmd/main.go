package main

import (
	"log"
	"medical-service/cmd/goroutines"
	"medical-service/pkg/config"
	"medical-service/pkg/logs"
	"medical-service/queue/rabbitMQ"
	"medical-service/storage/mongoDb"
)

func main() {
	cfg := config.Load()
	logger := logs.InitLogger()

	db, err := mongoDb.ConnectMongoDB(cfg)
	if err != nil {
		logger.Error("Error in MongoDB", "error", err)
		log.Fatal(err)
	}

	rabbit, err := rabbitMQ.ConnectToRabbit(cfg)
	if err != nil {
		logger.Error("Error in RabbitMQ", "error", err)
		log.Fatal(err)
	}

	goroutine := goroutines.NewGoroutines(db, rabbit, logger, cfg)
	wait := make(chan bool)

	go goroutine.RunGRPC()
	go goroutine.CreateMedicalRecords()
	go goroutine.DeleteMedicalRecords()
	go goroutine.CreateLifestyle()
	go goroutine.DeleteLifestyle()
	go goroutine.CreateWearableData()
	go goroutine.DeleteWearableData()

	<-wait
}
