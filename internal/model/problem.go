package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TestCase struct {
	Input       string `bson:"input"`
	Output      string `bson:"output"`
	Note        string `bson:"note,omitempty"`
	Explanation string `bson:"explanation,omitempty"`
}

type Problem struct {
	PID         primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	TimeLimit   int                `bson:"timeLimit"`
	MemoryLimit int                `bson:"memoryLimit"`
	Description string             `bson:"description"`
	TestSamples []TestCase         `bson:"testSamples"`

	// options
	DataRange string   `bson:"dataRange,omitempty"`
	Point     int      `bson:"point,omitempty"`
	Tags      []string `bson:"tags,omitempty"`
	Tutorial  string   `bson:"tutorial,omitempty"`

	// hide
	TestCases []TestCase `bson:"testCases" json:"-"`
}
