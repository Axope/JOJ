package service

import (
	"context"
	"encoding/json"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/model"
	"github.com/Axope/JOJ/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type problemService struct {
}

var ProblemService = new(problemService)

func (p *problemService) GetProblemList(req *request.GetProblemListRequest) ([]model.Problem, error) {
	findOptions := options.Find().SetLimit(req.Length).SetSkip(req.StartIndex - 1)
	findOptions.SetProjection(bson.D{{Key: "testCases", Value: 0}})

	cursor, err := dao.GetProblemColl().Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	} else {
		var results []model.Problem
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
		return results, nil
	}
}

func (p *problemService) GetProblem(req *request.GetProblemRequest) (*model.Problem, error) {
	pid, err := primitive.ObjectIDFromHex(req.PID)
	if err != nil {
		return nil, err
	}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{{Key: "testCases", Value: 0}})

	var result model.Problem
	err = dao.GetProblemColl().FindOne(context.TODO(), bson.D{{Key: "_id", Value: pid}}, findOptions).Decode(&result)
	if err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

func (p *problemService) CreateProblem(req *request.CreateProblemRequest, pid primitive.ObjectID,
	zipPath, testCasesPath string) error {
	// unzip
	if err := utils.Unzip(zipPath, testCasesPath); err != nil {
		return err
	}
	utils.RemoveFile(zipPath)

	var testSamples []model.TestCase
	if err := json.Unmarshal([]byte(req.TestSamplesJson), &testSamples); err != nil {
		return err
	}

	newProblem := model.Problem{
		PID:         pid,
		Title:       req.Title,
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
		Description: req.Description,
		TestSamples: testSamples,

		TestCasesPath: testCasesPath,
	}
	log.Logger.Debug("new problem", log.Any("newProblem", newProblem))
	// mongo
	insertResult, err := dao.GetProblemColl().InsertOne(context.TODO(), newProblem)
	if err != nil {
		return err
	}

	if insertResult.InsertedID.(primitive.ObjectID) != pid {
		panic("unknown error: pid not equal")
	}

	return nil
}
