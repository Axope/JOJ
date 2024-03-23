package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LangSet string

const (
	CPP    LangSet = "Cpp"
	JAVA   LangSet = "Java"
	PYTHON LangSet = "Python"
	GO     LangSet = "Go"
)

type StatusSet string

const (
	PENDING   StatusSet = "Pending"
	COMPILING StatusSet = "Compiling"
	JUDGING   StatusSet = "Judging"

	CE  StatusSet = "Compile Error"
	AC  StatusSet = "Accept"
	WA  StatusSet = "Wrong Answer"
	TLE StatusSet = "Time Limit Exceeded"
	MLE StatusSet = "Memory Limit Exceeded"
	RE  StatusSet = "Runtime Error"
	OLE StatusSet = "Output Limit Exceeded"
	UKE StatusSet = "Unknow Error"
)

type Submission struct {
	SID           primitive.ObjectID `bson:"_id,omitempty"`
	UID           uint               `bson:"uid"`
	PID           primitive.ObjectID `bson:"pid"`
	SubmitTime    time.Time          `bson:"submitTime"`
	Lang          LangSet            `bson:"lang"`
	Status        StatusSet          `bson:"status"`
	RunningTime   int                `bson:"runningTime"`
	RunningMemory int                `bson:"runningMemory"`
	SubmitCode    string             `bson:"submitCode"`

	// options
	Point int `bson:"point,omitempty"`
}
