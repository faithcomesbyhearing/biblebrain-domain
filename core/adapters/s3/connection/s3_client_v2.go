package connection

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func GetS3Client() *s3.Client {
	// create custom S3 resolver pointing to s3-local if ENVVAR IS_OFFLINE is set (to any value)
	_, isOffline := os.LookupEnv("IS_OFFLINE")
	s3LocalResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if isOffline && service == s3.ServiceID && region == "us-west-2" {
			log.Println("OFFLINE, so using S3 local endpoint")
			endpointURL := os.Getenv("AWS_SERVERLESS_S3_LOCAL_HOST")
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               endpointURL,
				HostnameImmutable: true,
				SigningRegion:     "us-west-2",
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	s3LocalCredentialsProvider := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		if isOffline {
			log.Println("OFFLINE, so using S3 local credentials")
			return aws.Credentials{
				AccessKeyID:     "S3RVER",
				SecretAccessKey: "S3RVER",
			}, nil
		}
		log.Println("not OFFLINE, so returning error indicating credentials not available from this provider")
		return aws.Credentials{}, &aws.EndpointNotFoundError{}

	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithDefaultRegion("us-west-2"),
		config.WithEndpointResolverWithOptions(s3LocalResolver),
		// config.WithCredentialsProvider(s3LocalCredentialsProvider),
	)
	cfg.Credentials = newOfflineAwareCredentialsProviderChain(
		s3LocalCredentialsProvider,
		cfg.Credentials,
	)
	if err != nil {
		log.Panic(err)
	}
	// Load the Shared AWS Configuration (~/.aws/config)

	// Create an Amazon S3 service client
	return s3.NewFromConfig(cfg)
}

func newOfflineAwareCredentialsProviderChain(providers ...aws.CredentialsProvider) aws.CredentialsProvider {
	return aws.NewCredentialsCache(
		aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			var errs []error

			for _, p := range providers {
				creds, err := p.Retrieve(ctx)
				if err == nil {
					log.Printf("resolving credentials providers... returning creds access key: %s, ", creds.AccessKeyID)
					return creds, nil
				}

				errs = append(errs, err)
			}

			return aws.Credentials{}, fmt.Errorf("no valid providers in chain: %s", errs)
		}),
	)
}

func ListObjectsByBucketV2(bucket string) (*s3.ListObjectsV2Output, error) {
	s3Client := GetS3Client()
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	result, err := s3Client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		log.Panic(err)
	}

	return result, nil
}

func ListObjectsByBucketAndPrefix(s3Client *s3.Client, bucket, prefix string) (*s3.ListObjectsV2Output, error) {
	log.Printf("ListObjectsByBucketAndPrefix. bucket: %s, prefix: %s", bucket, prefix)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	result, err := s3Client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		log.Panic(err)
	}

	return result, nil
}

func DownloadFileV2(objectKey string, bucket string, s3Client *s3.Client) ([]byte, error) {

	buffer := manager.NewWriteAtBuffer([]byte{})

	downloader := manager.NewDownloader(s3Client)

	numBytes, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}

	if numBytes < 1 {
		return nil, errors.New("zero bytes written to memory")
	}

	return buffer.Bytes(), nil
}

func UploadFile(file *bytes.Reader, key string, bucket string, s3Client *s3.Client) *manager.UploadOutput {
	uploader := manager.NewUploader(s3Client)

	log.Printf("uploadFile. bucket: %s, key: %s, filesize: %d", bucket, key, file.Len())
	output, uerr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLBucketOwnerFullControl,
		Body:   file,
	})

	if uerr != nil {
		log.Panic(uerr)
	}

	return output
}
