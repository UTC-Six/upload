package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Minio struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		BucketName      string
		MaxUploadSize   int64
	}
}
