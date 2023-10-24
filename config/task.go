package config

type Task struct {
	Concurrency int `mapstructure:"concurrency" json:"concurrency" yaml:"concurrency"`
}
