package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlc "github.com/faithcomesbyhearing/biblebrain-domain/mysql/generated"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetFilesForBible(t *testing.T) {
	// establish database connection
	db, err := sql.Open("mysql", "etl_dev:password@tcp(127.0.0.1:3306)/dbp_TEST")
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

	rows, err := queries.GetFilesForBible(ctx, "NYJBIB")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		fmt.Printf("Type: %s, BibleId: %s, Filesetid: %s, Filename: %s\n", row.Type, row.BibleID, row.ID, row.FileName)
	}
}
