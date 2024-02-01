package config

import "github.com/spf13/viper"

// -- CronJob configuraion

type CronJob struct {
	Crontab       string `yaml:"crontab"`
	WithSeconds   bool   `yaml:"withSeconds"`
	ExecImmediate bool   `yaml:"execImmediate"`
}

func setDefaultCronJobValues() {
	viper.SetDefault("crontab.crontab", "0 */12 * * *")
	viper.SetDefault("crontab.withSeconds", false)
	viper.SetDefault("crontab.execImmediate", false)
}
