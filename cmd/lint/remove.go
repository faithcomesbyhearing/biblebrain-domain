package lint

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/faithcomesbyhearing/biblebrain-domain/core/domain/storage"
)

func RemoveFromS3(bucket, filesetId string) {
	// read the specified file
	// read from fs
	inFileName := "toRemove/" + filesetId + "-s3.json"

	in, err := os.Open(inFileName)
	if err != nil {
		fmt.Println("No S3 file for fileset " + filesetId)
		return
	}
	defer in.Close()

	var toDelete []storage.Prefix
	decoder := json.NewDecoder(in)
	if err := decoder.Decode(&toDelete); err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return
	}

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
		outFileName := "removed/" + filesetId + ".json"

		err = os.Rename(inFileName, outFileName)
		if err != nil {
			log.Printf("Unable to move file %v to %v\n", inFileName, outFileName)
		}
	}
}

func RemoveFromDb(filesetId string) {
	// read the specified file
	// read from fs
	inFileName := "toRemove/" + filesetId + "-db.json"

	in, err := os.Open(inFileName)
	if err != nil {
		fmt.Println("No DB file for fileset " + filesetId)
		return
	}
	defer in.Close()

	var toDelete []storage.Prefix
	decoder := json.NewDecoder(in)
	if err := decoder.Decode(&toDelete); err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return
	}

	// delete from bible_files where filename = x
	// establish database connection
	// FIXME: change to a read-only database user
	db, err := sql.Open("mysql", "etl:E87xjHLeXzxtkKr@tcp(127.0.0.1:3306)/dbp_NEWDATA")
	// db, err := sql.Open("mysql", "etl_dev:password@tcp(127.0.0.1:3306)/dbp_TEST")
	if err != nil {
		fmt.Println("sql open failed")
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("sql ping failed")
		panic(err)
	}

	// queries := sqlc.New(db)
	// ctx := context.Background()

	// for _, prefix := range toDelete {
	// 	// execute delete
	// }

	// rows, err := queries.DeleteFilesetFiles(ctx, filesetId)
	// if err != nil {
	// 	panic(err)
	// }
}
