package configs

type LogConfig struct {
	Level      string
	Path       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}
