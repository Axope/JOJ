package configs

type MongoConfig struct {
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Database       string `yaml:"database"`
	ProblemColl    string `yaml:"problemColl"`
	SubmissionColl string `yaml:"submissionColl"`
}
