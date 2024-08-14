package goroutines

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	pb "medical-service/generated/healthAnalytics"
	"medical-service/pkg/config"
	"medical-service/queue"
	"medical-service/service"
	"medical-service/storage/mongoDb"
	"net"
)

type Goroutines interface {
	RunGRPC()
	// MedicalRecords
	CreateMedicalRecords()
	DeleteMedicalRecords()
	// Lifestyle
	CreateLifestyle()
	DeleteLifestyle()
	// WearableData
	CreateWearableData()
	DeleteWearableData()
}

func NewGoroutines(db *mongo.Database, rabbit *amqp.Connection, log *slog.Logger, cfg config.Config) Goroutines {
	return &goroutines{db: db, rabbit: rabbit, cfg: cfg, logger: log}
}

type goroutines struct {
	logger *slog.Logger
	cfg    config.Config
	db     *mongo.Database
	rabbit *amqp.Connection
}

func (g *goroutines) RunGRPC() {
	listner, err := net.Listen("tcp", g.cfg.GRPC_SERVER_PORT)
	if err != nil {
		g.logger.Error("Error in net/Listen", "error", err)
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	storageMedicalRecords := mongoDb.NewMedicalRecordRepo(g.db.Collection("health-records"))
	storageLifestyle := mongoDb.NewLifestyle(g.db.Collection("health-lifestyle"))
	storageHealthRecommend := mongoDb.NewHealth(g.db.Collection("health-recommend"), g.db)
	storageWearableData := mongoDb.NewWearableData(g.db.Collection("health-wearable"))

	serviceMedicalRecords := service.NewMedicalRecords(storageMedicalRecords, g.logger)
	serviceLifestyle := service.NewLifestyle(storageLifestyle, g.logger)
	serviceHealthRecommend := service.NewHealth(storageHealthRecommend, g.logger)
	serviceWearableData := service.NewWearableDate(storageWearableData, g.logger)

	pb.RegisterMedicalRecordsServiceServer(grpcServer, serviceMedicalRecords)
	pb.RegisterLifestyleServiceServer(grpcServer, serviceLifestyle)
	pb.RegisterHealthRecommendationsServiceServer(grpcServer, serviceHealthRecommend)
	pb.RegisterWearableDataServer(grpcServer, serviceWearableData)

	log.Println("service is running on port :  " + g.cfg.GRPC_SERVER_PORT)
	log.Fatal(grpcServer.Serve(listner))
}

//------------------------- Medical Records ---------------------------------------------------------

func (g *goroutines) CreateMedicalRecords() {

	storage := mongoDb.NewMedicalRecordRepo(g.db.Collection("health-records"))
	medicalRecords := service.NewMedicalRecords(storage, g.logger)

	medicalQ := queue.NewMedicalRecordConsume(g.logger, g.rabbit)
	ch, err := medicalQ.CreateConsume()
	if err != nil {
		g.logger.Error("Error in CreateConsume", "error", err)
		log.Fatal(err)
	}

	log.Println("rabbitMQ reading from queue: ", "CreateMedicalRecords")

	for msg := range ch {
		log.Println("11111")
		var create pb.AddMedicalRecordRequest

		err = json.Unmarshal([]byte(msg.Body), &create)
		if err != nil {
			g.logger.Error("Error in Unmarshal", "error", err)
			log.Fatal(err)
		}

		_, err = medicalRecords.AddMedicalRecord(context.Background(), &create)
		if err != nil {
			g.logger.Error("Error in AddMedicalRecord", "error", err)
			log.Fatal(err)
		}
	}
}

func (g *goroutines) DeleteMedicalRecords() {

	storage := mongoDb.NewMedicalRecordRepo(g.db.Collection("health-records"))
	medicalRecords := service.NewMedicalRecords(storage, g.logger)

	medicalQ := queue.NewMedicalRecordConsume(g.logger, g.rabbit)
	ch, err := medicalQ.DeleteConsume()
	if err != nil {
		g.logger.Error("Error in CreateConsume", "error", err)
		log.Fatal(err)
	}

	log.Println("rabbitMQ reading from queue: ", "DeleteMedicalRecords")

	for msg := range ch {

		var id pb.MedicalRecordID

		err = json.Unmarshal([]byte(msg.Body), &id)
		if err != nil {
			g.logger.Error("Error in Unmarshal", "error", err)
			log.Fatal(err)
		}
		log.Println(id)
		_, err = medicalRecords.DeleteMedicalRecord(context.Background(), &id)
		if err != nil {
			g.logger.Error("Error in AddMedicalRecord", "error", err)
			log.Fatal(err)
		}
	}
}

// ----------------------------------- Lifestyle ----------------------------------------------------

func (g *goroutines) CreateLifestyle() {

	storage := mongoDb.NewLifestyle(g.db.Collection("health-lifestyle"))
	medicalRecords := service.NewLifestyle(storage, g.logger)

	medicalQ := queue.NewLifestyleConsumer(g.logger, g.rabbit)
	ch, err := medicalQ.CreateConsume()
	if err != nil {
		g.logger.Error("Error in CreateConsume", "error", err)
		log.Fatal(err)
	}

	log.Println("rabbitMQ reading from queue: ", "CreateLifestyle")

	for msg := range ch {

		var create pb.Lifestyle
		err = json.Unmarshal([]byte(msg.Body), &create)
		if err != nil {
			g.logger.Error("Error in Unmarshal", "error", err)
			log.Fatal(err)
		}

		_, err = medicalRecords.AddLifestyleData(context.Background(), &create)
		if err != nil {
			g.logger.Error("Error in AddMedicalRecord", "error", err)
			log.Fatal(err)
		}
	}
}

func (g *goroutines) DeleteLifestyle() {

	storage := mongoDb.NewLifestyle(g.db.Collection("health-lifestyle"))
	medicalRecords := service.NewLifestyle(storage, g.logger)

	medicalQ := queue.NewLifestyleConsumer(g.logger, g.rabbit)
	ch, err := medicalQ.DeleteConsume()
	if err != nil {
		g.logger.Error("Error in CreateConsume", "error", err)
		log.Fatal(err)
	}

	log.Println("rabbitMQ reading from queue: ", "DeleteLifestyle")

	for msg := range ch {
		log.Println("4444444")
		var id pb.LifestyleID

		err = json.Unmarshal([]byte(msg.Body), &id)
		if err != nil {
			g.logger.Error("Error in Unmarshal", "error", err)
			log.Fatal(err)
		}
		log.Println(id)
		_, err = medicalRecords.DeleteLifestyleData(context.Background(), &id)
		if err != nil {
			g.logger.Error("Error in AddMedicalRecord", "error", err)
			log.Fatal(err)
		}
	}
}

// ---------------------------- Wearable Data --------------------------------------------------------

func (g *goroutines) CreateWearableData() {

	storage := mongoDb.NewWearableData(g.db.Collection("health-wearable"))
	medicalRecords := service.NewWearableDate(storage, g.logger)

	medicalQ := queue.NewWearableData(g.logger, g.rabbit)
	ch, err := medicalQ.CreateConsume()
	if err != nil {
		g.logger.Error("Error in CreateConsume", "error", err)
		log.Fatal(err)
	}

	log.Println("rabbitMQ reading from queue: ", "CreateWearableData")

	for msg := range ch {
		log.Println("555555")
		var create pb.WearableDate
		err = json.Unmarshal([]byte(msg.Body), &create)
		if err != nil {
			g.logger.Error("Error in Unmarshal", "error", err)
			log.Fatal(err)
		}

		_, err = medicalRecords.AddWearableData(context.Background(), &create)
		if err != nil {
			g.logger.Error("Error in AddMedicalRecord", "error", err)
			log.Fatal(err)
		}
	}
}

func (g *goroutines) DeleteWearableData() {

	storage := mongoDb.NewWearableData(g.db.Collection("health-wearable"))
	medicalRecords := service.NewWearableDate(storage, g.logger)

	medicalQ := queue.NewWearableData(g.logger, g.rabbit)
	ch, err := medicalQ.DeleteConsume()
	if err != nil {
		g.logger.Error("Error in CreateConsume", "error", err)
		log.Fatal(err)
	}

	log.Println("rabbitMQ reading from queue: ", "DeleteWearableData")

	for msg := range ch {

		var id pb.WearableDataID

		err = json.Unmarshal([]byte(msg.Body), &id)
		if err != nil {
			g.logger.Error("Error in Unmarshal", "error", err)
			log.Fatal(err)
		}
		log.Println(id)
		_, err = medicalRecords.DeleteWearableData(context.Background(), &id)
		if err != nil {
			g.logger.Error("Error in AddMedicalRecord", "error", err)
			log.Fatal(err)
		}
	}
}
