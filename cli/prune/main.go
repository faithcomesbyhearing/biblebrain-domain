package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/faithcomesbyhearing/biblebrain-domain/core/domain/storage"
)

func main() {
	var filesetId string
	flag.StringVar(&filesetId, "filesetId", "", "filesetId to be pruned")
	flag.Parse()
	if len(filesetId) == 0 {
		fmt.Println("provide filesetid as command line arg")
		return
	}

	// read the specified file
	// read from fs
	in, err := os.Open("in/NYJBIBO1DA.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer in.Close()

	var toDelete []storage.Prefix
	decoder := json.NewDecoder(in)
	if err := decoder.Decode(&toDelete); err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return
	}

	fmt.Println("Slice read from JSON file:", toDelete)

	// prepare to navigate through S3
	// FIXME: change to a read-only S3 profile
	bucket := "dbp-staging"
	profile := "dbsxx"

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)

	svc := s3.NewFromConfig(cfg)

	var keys []types.ObjectIdentifier

	for _, prefix := range toDelete {
		keys = append(keys, types.ObjectIdentifier{Key: aws.String(prefix.String())})
	}

	output, err := svc.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{Objects: keys},
	})

	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", bucket, err)
	} else {
		log.Printf("Deleted %v objects.\n", len(output.Deleted))
	}
}
