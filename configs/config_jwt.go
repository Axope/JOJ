package configs

type JWTConfig struct {
	SigningKey string `yaml:"signingKey"`
	Expire     string `yaml:"expire"`
}
