package dao

import (
	"context"
	"fmt"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	client      *mongo.Client
	problemColl *mongo.Collection
	submissionColl  *mongo.Collection
)

func InitMongo() error {
	cfg := configs.GetMongoConfig()
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port
	database := cfg.Database
	problemCollName := cfg.ProblemColl
	submissionCollName := cfg.SubmissionColl
	log.Logger.Sugar().Debugf("database %v, problemCollName %v, submissionCollName %v",
		database, problemCollName, submissionCollName)

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/", username, password, host, port)
	log.Logger.Debug(uri)

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return fmt.Errorf("mongo ping error: %v", err)
	}
	log.Logger.Debug("", zap.Any("result", result))

	db := client.Database(database)
	problemColl = db.Collection(problemCollName)
	submissionColl = db.Collection(submissionCollName)

	// {
	// 	// debug
	// 	cursor, err := problemColl.Find(context.Background(), bson.D{{}})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	var results []model.Problem
	// 	if err = cursor.All(context.TODO(), &results); err != nil {
	// 		panic(err)
	// 	}
	// 	log.Logger.Debug("", zap.Any("slice", results))
	// }

	return nil
}

func GetClient() *mongo.Client {
	return client
}
func GetProblemColl() *mongo.Collection {
	return problemColl
}
func GetSubmissionColl() *mongo.Collection {
	return submissionColl
}
