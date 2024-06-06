package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sqlc "github.com/faithcomesbyhearing/biblebrain-domain/adapters/mysql/generated"
	"github.com/faithcomesbyhearing/biblebrain-domain/core/domain/storage"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var filesetId string
	flag.StringVar(&filesetId, "filesetId", "", "filesetId to be pruned")
	flag.Parse()
	if len(filesetId) == 0 {
		fmt.Println("provide filesetid as command line arg")
		return
	}

	extensions := []string{"mp3", "webm"}

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

	queries := sqlc.New(db)
	ctx := context.Background()

	rows, err := queries.GetFilesForFileset(ctx, filesetId)
	if err != nil {
		panic(err)
	}

	prefixRoot, err := queries.GetStoragePrefixForFileset(ctx, filesetId)
	if err != nil {
		panic(err)
	}
	metadata := []storage.Prefix{}
	for _, row := range rows {
		// fmt.Printf("Type: %s, BibleId: %s, Filesetid: %s, Filename: %s\n", row.Type, row.BibleID, row.ID, row.FileName)
		metadata = append(metadata, storage.NewPrefix("audio", row.BibleID, row.ID, row.FileName))
	}

	slices.SortFunc(metadata, func(a, b storage.Prefix) int {
		return storage.Compare(a, b)
	})

	// fmt.Printf("metadata.. %v, len=%d cap=%d\n", metadata, len(metadata), cap(metadata))
	fmt.Printf("metadata.. len=%d\n", len(metadata))

	// prepare to navigate through S3
	// FIXME: change to a read-only S3 profile
	bucket := "dbp-prod"
	profile := "dbs"

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)

	svc := s3.NewFromConfig(cfg)

	var continuationToken *string

	media := []storage.Prefix{}
	for {
		// List objects in the specified S3 bucket
		resp, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(prefixRoot),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			fmt.Println("Unable to list items in bucket:", err)
			os.Exit(1)
		}

		// A01___01
		for _, item := range resp.Contents {
			prefix := storage.Parse(aws.ToString(item.Key))

			// TODO: fancy this up as a filter
			ext := extractExtension(prefix.Filename)
			if len(ext) > 0 && slices.Contains(extensions, ext) {
				media = append(media, storage.Parse(aws.ToString(item.Key)))
			}

		}

		// check for more results to fetch
		if resp.NextContinuationToken == nil {
			break
		}
		continuationToken = resp.NextContinuationToken
	}

	// sort
	slices.SortFunc(media, func(a, b storage.Prefix) int {
		return storage.Compare(a, b)
	})

	// fmt.Printf("media...  %v, len=%d cap=%d\n", media, len(media), cap(media))
	fmt.Printf("media.. len=%d\n", len(media))

	// now find media that is not represented in metadata
	orphans := []storage.Prefix{}
	for _, m := range media {
		found := slices.ContainsFunc(metadata, func(p storage.Prefix) bool {
			return storage.Compare(p, m) == 0
		})
		if !found {
			orphans = append(orphans, m)
		}
	}

	fmt.Printf("orphans.. len=%d\n", len(orphans))

	// write to fs
	file, err := os.Create(filesetId + ".json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(orphans); err != nil {
		fmt.Println("Error encoding data to JSON:", err)
	}
}

func extractExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}
