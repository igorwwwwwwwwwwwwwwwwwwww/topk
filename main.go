package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	_ "github.com/duckdb/duckdb-go/v2"
)

var k = flag.Int("k", 10, "limit to this many top values")
var other = flag.Bool("other", false, "include sum count of remaining values")

func main() {
	flag.Parse()

	ctx := context.Background()

	filename := "/dev/stdin"
	if flag.NArg() > 0 {
		filename = flag.Arg(0)
	}

	db, err := sql.Open("duckdb", "")
	if err != nil {
		slog.Error("failed to open db", "err", err)
		os.Exit(1)
	}

	query := `
		WITH all_data AS (
			SELECT #1 as item
			FROM read_csv(?1, header=FALSE, delim='\n')
		),
		counts AS (
			SELECT item, COUNT(*) as cnt
			FROM all_data
			GROUP BY item
		),
		topk AS (
			SELECT item, cnt, ROW_NUMBER() OVER (ORDER BY cnt DESC) as rn
			FROM counts
		),
		results AS (
			SELECT item, cnt FROM topk WHERE rn <= ?2
			UNION ALL
			SELECT 'OTHER', COALESCE(SUM(cnt), 0) FROM topk WHERE rn > ?2 AND ?3
		)
		SELECT
			item,
			cnt,
			(SELECT MAX(LENGTH(item)) FROM results) as max_len,
			CAST(cnt * 50.0 / MAX(cnt) OVER () AS INTEGER) as bar_width
		FROM results
		ORDER BY CASE WHEN item = 'OTHER' THEN 1 ELSE 0 END, cnt DESC`

	rows, err := db.QueryContext(ctx, query, filename, *k, *other)
	if err != nil {
		log.Fatalf("could not query db: %s", err.Error())
	}
	defer rows.Close()

	var (
		item                    string
		count, maxLen, barWidth int
	)

	for rows.Next() {
		if err := rows.Scan(&item, &count, &maxLen, &barWidth); err != nil {
			log.Fatalf("could not get row: %s", err.Error())
		}

		if item == "OTHER" && count == 0 {
			continue
		}

		bar := strings.Repeat("∎", barWidth)
		fmt.Printf("%-*s  %6d  %s\n", maxLen, item, count, bar)
	}
}
