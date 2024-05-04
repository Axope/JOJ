package contest

type RankListData struct {
	Uid string               `bson:"uid" json:"uid"`
	Pss []ProblemSolveStatus `bson:"pss" json:"pss"`
}
