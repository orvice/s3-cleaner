package biz

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/orvice/s3-cleaner/internal/config"
	log "github.com/sirupsen/logrus"
)

var (
	cli *minio.Client
)

func Clean(cfg config.Config) {
	endpoint := cfg.Endpoint
	accessKeyID := cfg.AccessKeyID
	secretAccessKey := cfg.SecretAccessKey
	useSSL := cfg.UseSSL

	log.Info("endpoint", endpoint)
	var err error
	// Initialize minio client object.
	cli, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	for _, bucket := range cfg.BucketConfigs {
		cleanBucket(ctx, bucket.Name, bucket.Prefix)
	}
}

func cleanBucket(ctx context.Context, bucket string, prefix []string) {
	if len(prefix) == 0 {
		cleanBucketOjbect(ctx, bucket, "")
		return
	}

	for _, p := range prefix {
		cleanBucketOjbect(ctx, bucket, p)
	}
}

func cleanBucketOjbect(ctx context.Context, bucket, prefix string) {
	log.Infof("list objects %s %s", bucket, prefix)
	resp := cli.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range resp {
		log.Info(obj.Key)

		if obj.LastModified.Before(time.Now().Add(-time.Hour * 24 * 90)) {
			log.Infof("remove %s last modify %s", obj.Key, obj.LastModified.String())
			removeOjbect(ctx, bucket, obj.Key)
		}
	}

}

func removeOjbect(ctx context.Context, bucket, key string) {
	err := cli.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		log.Errorf("remove ojb %s error %s", key, err.Error())
	}
}
