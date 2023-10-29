package main

import (
	"context"
	"encoding/csv"
	"flag"
	"log"
	"os"

	"github.com/cxiang03/dbench"
)

func main() {
	csv, dsn := parseArgs()

	db := dbench.NewDB(dsn)
	ctx := context.Background()
	now := ""
	err := db.NewSelect().NewRaw("SELECT NOW()").Scan(ctx, &now)
	log.Println("select now:", now)
	log.Println("error:", err)
	if err != nil {
		panic("database is not ready, please check your DSN")
	}

	linesCh := readByBatch(csv, 1000)
	for lines := range linesCh {
		records := make([]*dbench.Record, 0, len(lines))
		for _, line := range lines {
			records = append(records, dbench.ParseRecord(line))
		}
		if err := db.Insert(ctx, records); err != nil {
			log.Println("error:", err)
			panic("failed to insert records")
		}
		log.Println("one batch finished, len=", len(records))
	}
}

func parseArgs() (string, string) {
	filePath := flag.String("f", "./pp-monthly-update-new-version.csv", "CSV data to load")
	databaseDSN := flag.String("d", "root:password@tcp(127.0.0.1:3306)/dbench", "Database DSN")
	flag.Parse()
	return *filePath, *databaseDSN
}

func readByBatch(filename string, batchSize int) <-chan [][]string {
	rst := make(chan [][]string)
	file, err := os.Open(filename)
	if err != nil {
		close(rst)
		return nil
	}

	go func() {
		defer file.Close()
		defer close(rst)

		reader := csv.NewReader(file)
		for {
			batch := make([][]string, 0, batchSize)

			for i := 0; i < batchSize; i++ {
				line, err := reader.Read()
				if err != nil {
					break
				}
				batch = append(batch, line)
			}

			rst <- batch
			if len(batch) < batchSize {
				break
			}
		}
	}()

	return rst
}
