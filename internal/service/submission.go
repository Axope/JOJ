package service

import (
	"context"

	"github.com/Axope/JOJ/common/log"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/internal/dao"
	"github.com/Axope/JOJ/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type submissionService struct {
}

var SubmissionService = new(submissionService)

func (s *submissionService) GetSubmissionList(req *request.GetSubmissionListRequest) ([]model.Submission, error) {
	pid, err := primitive.ObjectIDFromHex(req.PID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "uid", Value: req.UID}, {Key: "pid", Value: pid}}
	log.LoggerSugar.Debugf("uid = %v, pid = %v", req.UID, pid)
	cursor, err := dao.GetSubmissionColl().Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	} else {
		var results []model.Submission
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
		return results, nil
	}
}

// func (s *submissionService) GetSubmissionListByUID(uid uint) ([]model.Submission, error) {

// }

// func (s *submissionService) GetSubmissionListByPID(pid primitive.ObjectID) ([]model.Submission, error) {

// }
