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
	var filesetId, bucket string
	flag.StringVar(&filesetId, "filesetId", "", "filesetId to be pruned")
	flag.StringVar(&bucket, "bucket", "", "bucket name (eg dbp-prod)")
	flag.Parse()
	if len(filesetId) == 0 {
		fmt.Println("provide filesetid as command line arg")
		return
	}
	if len(bucket) == 0 {
		bucket = "dbp-staging"
	}

	// read the specified file
	// read from fs
	inFileName := "delete/" + filesetId + ".json"

	in, err := os.Open(inFileName)
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

	// fmt.Println("Slice read from JSON file:", toDelete)

	profile := "dbs"

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
		outFileName := "processed/" + filesetId + ".json"

		err = os.Rename(inFileName, outFileName)
		if err != nil {
			log.Printf("Unable to move file %v to %v\n", inFileName, outFileName)
		}
	}
}
