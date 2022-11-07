package biz

import (
	"context"
	"sync"
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

	wg := new(sync.WaitGroup)
	wg.Add(len(cfg.BucketConfigs))
	for _, bucket := range cfg.BucketConfigs {
		go func(cfg config.BucketConfig) {
			defer wg.Done()
			cleanBucket(ctx, cfg.Name, cfg.Prefix)
		}(bucket)
	}
	wg.Wait()
}

func cleanBucket(ctx context.Context, bucket string, prefix []string) {
	if len(prefix) == 0 {
		cleanBucketOjbect(ctx, bucket, "")
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(prefix))
	for _, p := range prefix {
		go func(prefix string) {
			defer wg.Done()
			cleanBucketOjbect(ctx, bucket, prefix)
		}(p)
	}
	wg.Wait()
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
