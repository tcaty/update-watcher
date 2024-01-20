package config

type CronJob struct {
	Crontab       string `yaml:"crontab"`
	WithSeconds   bool   `yaml:"withSeconds"`
	ExecImmediate bool   `yaml:"execImmediate"`
}
