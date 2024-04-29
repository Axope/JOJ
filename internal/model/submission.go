package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LangSet int32

const (
	CPP    LangSet = iota //"Cpp"
	JAVA                  //"Java"
	PYTHON                //"Python"
	GO                    //"Go"
)

type StatusSet string

const (
	UNSUBMIT  StatusSet = "UnSubmit"
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
	UKE StatusSet = "Unknown Error"
)

type Submission struct {
	SID           primitive.ObjectID `bson:"_id,omitempty" json:"sid"`
	UID           uint               `bson:"uid" json:"uid"`
	PID           primitive.ObjectID `bson:"pid" json:"pid"`
	SubmitTime    time.Time          `bson:"submitTime" json:"submitTime"`
	Lang          LangSet            `bson:"lang" json:"lang"`
	Status        StatusSet          `bson:"status" json:"status"`
	RunningTime   int                `bson:"runningTime" json:"runningTime"`
	RunningMemory int                `bson:"runningMemory" json:"runningMemory"`
	SubmitCode    string             `bson:"submitCode" json:"submitCode"`

	// options
	Point int `bson:"point,omitempty" json:"point,omitempty"`
}
