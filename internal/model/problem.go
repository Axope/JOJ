package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TestCase struct {
	Input       string `bson:"input" json:"input" form:"input"`
	Output      string `bson:"output" json:"output" form:"output"`
	Note        string `bson:"note,omitempty" json:"note,omitempty" form:"note,omitempty"`
	Explanation string `bson:"explanation,omitempty" json:"explanation,omitempty" form:"explanation,omitempty"`
}

type Problem struct {
	PID          primitive.ObjectID `bson:"_id,omitempty" json:"pid"`
	Title        string             `bson:"title" json:"title"`
	TimeLimit    int64              `bson:"timeLimit" json:"timeLimit"`
	MemoryLimit  int64              `bson:"memoryLimit" json:"memoryLimit"`
	Description  string             `bson:"description" json:"description"`
	TestSamples  []TestCase         `bson:"testSamples" json:"testSamples"`
	InputFormat  string             `bson:"inputFormat" json:"inputFormat"`
	OutputFormat string             `bson:"outputFormat" json:"outputFormat"`

	// options
	DataRange string   `bson:"dataRange,omitempty" json:"dataRange,omitempty"`
	Point     int      `bson:"point,omitempty" json:"point,omitempty"`
	Tags      []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Tutorial  string   `bson:"tutorial,omitempty" json:"tutorial,omitempty"`

	// hide
	// TestCases []TestCase `bson:"testCases" json:"-"`
	TestCasesPath string `bson:"testCasesPath" json:"-"`
}
