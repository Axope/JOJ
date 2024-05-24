package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/configs"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/manager"
	"github.com/Axope/JOJ/internal/middleware/rabbitmq"
	"github.com/Axope/JOJ/internal/model"
	pb "github.com/Axope/JOJ/protocol"
	"github.com/Axope/JOJ/utils"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
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
		Lang:       model.LangSet(req.Lang),
		Status:     model.PENDING,
		SubmitCode: req.SubmitCode,
	}
	if req.CID != "" {
		cid, err := primitive.ObjectIDFromHex(req.CID)
		if err != nil {
			return err
		}
		newSubmission.CID = cid
	}
	insertResult, err := dao.GetSubmissionColl().InsertOne(context.TODO(), newSubmission)
	if err != nil {
		return err
	}

	cfg := configs.GetDatasConfig()
	// judgeReq := request.JudgeRequest{
	// 	SID:         insertResult.InsertedID.(primitive.ObjectID).Hex(),
	// 	PID:         result.PID.Hex(),
	// 	TimeLimit:   result.TimeLimit,
	// 	MemoryLimit: result.MemoryLimit,
	// 	TestCases:   utils.GetStringSlice(cfg.DirPath + pid.Hex() + "/" + cfg.TestCasesListFile),
	// 	Lang:        req.Lang,
	// 	SubmitCode:  req.SubmitCode,
	// }
	// msg, err := json.Marshal(judgeReq)
	// if err != nil {
	// 	return err
	// }
	// if err := rabbitmq.SendMsgByJson(msg); err != nil {
	// 	return err
	// }
	judgeReq := &pb.Judge{
		Sid:             insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Pid:             result.PID.Hex(),
		Uid:             strconv.FormatUint(uint64(req.UID), 10),
		Cid:             req.CID,
		TimeLimit:       result.TimeLimit,
		MemoryLimit:     result.MemoryLimit,
		TestCases:       utils.GetStringSlice(cfg.DirPath + pid.Hex() + "/" + cfg.TestCasesListFile),
		Lang:            pb.LangSet(req.Lang),
		SubmitCode:      req.SubmitCode,
		SubmitTimestamp: req.SubmitTime.Unix(),
	}
	msg, err := proto.Marshal(judgeReq)
	if err != nil {
		return err
	}
	if err := rabbitmq.SendMsgByProtobuf(msg); err != nil {
		return err
	}
	log.Logger.Info("send judge request success", log.Any("judgeReq", judgeReq))
	return nil
}

func (s *submitService) HandleSubmitResult(msgs <-chan amqp091.Delivery) {
	for msg := range msgs {
		judgeResult := &pb.JudgeResult{}
		// TODO: Unmarshal error, redeliver the message
		if err := proto.Unmarshal(msg.Body, judgeResult); err != nil {
			log.Logger.Error("json.Unmarshal", log.Any("err", err))
			break
		}
		log.Logger.Debug("judge result", log.Any("judgeResult", judgeResult))

		// process
		sid, err := primitive.ObjectIDFromHex(judgeResult.Sid)
		if err != nil {
			panic(err)
		}
		filter := bson.D{{Key: "_id", Value: sid}}
		status := model.GetStatusSet(judgeResult.Status)
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "status", Value: status},
			{Key: "executeTime", Value: judgeResult.ExecuteTime},
			{Key: "executeMemory", Value: judgeResult.ExecuteMemory},
		}}}
		log.LoggerSugar.Debugf("mongo filter(%v) update(%v)", filter, update)
		_, err = dao.GetSubmissionColl().UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Logger.Error("mongo update error", log.Any("judgeResult", judgeResult))
			panic(err)
		}
		if judgeResult.Cid != "" {
			contest, err := manager.ContestManager.GetContest(judgeResult.Cid)
			if err != nil {
				log.Logger.Error("contest miss", log.Any("cid", judgeResult.Cid))
				panic(err)
			}
			if model.GetStatusSet(judgeResult.Status) == model.AC {
				uid, err := strconv.ParseUint(judgeResult.Uid, 10, 64)
				if err != nil {
					log.Logger.Error("judgeResult uid parse error", log.Any("uid", judgeResult.Uid))
					panic(err)
				}
				pid, err := primitive.ObjectIDFromHex(judgeResult.Pid)
				if err != nil {
					log.Logger.Error("judgeResult pid parse error", log.Any("pid", judgeResult.Pid))
					panic(err)
				}
				submitTime := time.Unix(judgeResult.SubmitTimestamp, 0)

				if err := contest.Accept(uint(uid), pid, submitTime); err != nil {
					log.Logger.Error("contest accept error")
					panic(err)
				}
			} else if model.GetStatusSet(judgeResult.Status) == model.CE {
				log.Logger.Debug("CE")
			} else if model.GetStatusSet(judgeResult.Status) == model.UKE {
				log.Logger.Debug("UKE")
			} else {
				uid, err := strconv.ParseUint(judgeResult.Uid, 10, 64)
				if err != nil {
					log.Logger.Error("judgeResult uid parse error", log.Any("uid", judgeResult.Uid))
					panic(err)
				}
				pid, err := primitive.ObjectIDFromHex(judgeResult.Pid)
				if err != nil {
					log.Logger.Error("judgeResult pid parse error", log.Any("pid", judgeResult.Pid))
					panic(err)
				}

				if err := contest.Fail(uint(uid), pid); err != nil {
					log.Logger.Error("contest fail error")
					panic(err)
				}
			}
		}
		if err = msg.Ack(false); err != nil {
			panic(err)
		}
	}
}
