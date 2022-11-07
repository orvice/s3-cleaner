package config

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketConfigs   []BucketConfig
}

type BucketConfig struct {
	Name   string
	Prefix []string
}
