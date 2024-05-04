package model

import (
	"time"

	pb "github.com/Axope/JOJ/protocol"
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
	UNACCEPT  StatusSet = "UnAccept"
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
	CID           primitive.ObjectID `bson:"cid,omitempty" json:"cid,omitempty"`
	SubmitTime    time.Time          `bson:"submitTime" json:"submitTime"`
	Lang          LangSet            `bson:"lang" json:"lang"`
	Status        StatusSet          `bson:"status" json:"status"`
	RunningTime   int                `bson:"runningTime" json:"runningTime"`
	RunningMemory int                `bson:"runningMemory" json:"runningMemory"`
	SubmitCode    string             `bson:"submitCode" json:"submitCode"`

	// options
	Point int `bson:"point,omitempty" json:"point,omitempty"`
}

func GetStatusSet(status pb.StatusSet) StatusSet {
	switch status {
	case pb.StatusSet_CE:
		return CE
	case pb.StatusSet_AC:
		return AC
	case pb.StatusSet_WA:
		return WA
	case pb.StatusSet_TLE:
		return TLE
	case pb.StatusSet_MLE:
		return MLE
	case pb.StatusSet_RE:
		return RE
	case pb.StatusSet_OLE:
		return OLE
	case pb.StatusSet_UKE:
		return UKE
	}
	return UKE
}
