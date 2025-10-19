package internal

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient() *minio.Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		log.Fatalln("MINIO_ENDPOINT is not set")
	}
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	if accessKeyID == "" {
		log.Fatalln("ACCESS_KEY_ID is not set")
	}
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		log.Fatalln("SECRET_ACCESS_KEY is not set")
	}
	useSSL := os.Getenv("USE_SSL") != "false"
	ignoreCerts := os.Getenv("SKIP_TLS_VALIDATION") == "true"

	fmt.Println("Creating MinIO client with endpoint:", endpoint)
	fmt.Println("Using SSL:", useSSL)
	fmt.Println("Ignoring certificate errors:", ignoreCerts)

	// Initialize minio client object.
	tlsConfig := &tls.Config{}
	if ignoreCerts {
		tlsConfig.InsecureSkipVerify = true
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
			DisableCompression:    true,
		},
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}
