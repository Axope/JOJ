package service

import (
	"context"
	"encoding/json"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/middleware/rabbitmq"
	"github.com/Axope/JOJ/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type submitService struct {
}

var SubmitService = new(submitService)

func (s *submitService) Submit(req *request.SubmitRequest) error {
	pid, err := primitive.ObjectIDFromHex(req.PID)
	if err != nil {
		return err
	}

	findOptions := options.FindOne()
	findOptions.SetProjection(bson.D{
		{Key: "timeLimit", Value: 1},
		{Key: "memoryLimit", Value: 1},
		{Key: "testSamples", Value: 1},
		{Key: "testCases", Value: 1},
	})
	var result model.Problem
	err = dao.GetProblemColl().FindOne(context.TODO(), bson.D{{Key: "_id", Value: pid}}, findOptions).Decode(&result)
	if err != nil {
		return err
	}

	newSubmission := model.Submission{
		SID:        primitive.NewObjectID(),
		UID:        req.UID,
		PID:        result.PID,
		SubmitTime: req.SubmitTime,
		Lang:       req.Lang,
		Status:     model.PENDING,
		SubmitCode: req.SubmitCode,
	}
	insertResult, err := dao.GetSubmissionColl().InsertOne(context.TODO(), newSubmission)
	if err != nil {
		return err
	}

	judgeReq := request.JudgeRequest{
		SID:         insertResult.InsertedID.(primitive.ObjectID).Hex(),
		TimeLimit:   result.TimeLimit,
		MemoryLimit: result.MemoryLimit,
		TestSamples: result.TestCases,
		TestCases:   result.TestCases,
		Lang:        req.Lang,
		SubmitCode:  req.SubmitCode,
	}
	msg, err := json.Marshal(judgeReq)
	if err != nil {
		return err
	}
	if err := rabbitmq.SendMsgByJson(msg); err != nil {
		return err
	}
	log.Logger.Info("send judge request success", log.Any("judgeReq", judgeReq))
	return nil
}
