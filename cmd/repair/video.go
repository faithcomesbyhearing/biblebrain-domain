package repair

import (
	"context"
	"database/sql"
	"fmt"

	sqlc "github.com/faithcomesbyhearing/biblebrain-domain/adapters/mysql/generated"

	_ "github.com/go-sql-driver/mysql"
)

func StitchHls() {
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

	rows, err := queries.HlsStitchingDrillDown(ctx, "INDALA")
	if err != nil {
		panic(err)
	}

	for row := range rows {
		fmt.Println(row)
	}
}
