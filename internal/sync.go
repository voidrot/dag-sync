package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

func SyncDags() {
	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		log.Fatalln("BUCKET_NAME is not set")
	}
	targetDir := os.Getenv("TARGET_DIR")
	if targetDir == "" {
		log.Fatalln("TARGET_DIR is not set")
	}
	bucketPrefix := os.Getenv("BUCKET_PREFIX")

	client := NewClient()

	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    bucketPrefix,
	}

	log.Printf("Attempting to sync DAG's from %s/%s\n", bucket, bucketPrefix)

	for object := range client.ListObjects(context.Background(), bucket, opts) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		reader, err := client.GetObject(context.Background(), bucket, object.Key, minio.GetObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		defer reader.Close()

		localFile, err := os.Create(fmt.Sprintf("%s/%s", targetDir, object.Key))
		if err != nil {
			log.Fatalln(err)
		}
		defer localFile.Close()

		stat, err := reader.Stat()
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
			log.Fatalln(err)
		}
		log.Printf("Downloaded %s to %s/%s\n", object.Key, targetDir, object.Key)
	}

}
