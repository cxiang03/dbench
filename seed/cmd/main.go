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
	csv, dsn, esa := parseArgs()
	_ = esa

	db := dbench.NewDB(dsn)
	ctx := context.Background()
	now := ""
	err := db.NewSelect().NewRaw("SELECT NOW()").Scan(ctx, &now)
	log.Println("select now:", now)
	log.Println("error:", err)
	if err != nil {
		panic("database is not ready, please check your DSN")
	}

	// es := dbench.NewEsClient(esa)
	// err = es.Ping()
	// log.Println("error:", err)
	// if err != nil {
	// 	panic("elasticsearch is not ready, please check your elastic search address")
	// }

	redis := dbench.NewRedis()
	err = redis.Ping(ctx).Err()
	log.Println("error:", err)
	if err != nil {
		panic("redis is not ready, please check your redis address")
	}

	meili := dbench.NewMeili()
	v, err := meili.Version()
	log.Println("version", v, "error:", err)
	if err != nil {
		panic("meili is not ready, please check your meili address")
	}
	if err = meili.UpdateFilterableAttributes(ctx); err != nil {
		panic("failed to update filterable attributes")
	}

	linesCh := readByBatch(csv, 1000)
	for lines := range linesCh {
		records := make([]*dbench.Record, 0, len(lines))
		for _, line := range lines {
			records = append(records, dbench.ParseRecord(line))
		}
		// if err := db.Insert(ctx, records); err != nil {
		// 	log.Println("error:", err)
		// 	panic("failed to insert records")
		// }
		// if err := es.Insert(ctx, records); err != nil {
		// 	log.Println("error:", err)
		// 	panic("failed to insert records")
		// }
		if err := redis.Insert(ctx, records); err != nil {
			log.Println("error:", err)
			panic("failed to insert records")
		}
		if err := meili.Insert(ctx, records); err != nil {
			log.Println("error:", err)
			panic("failed to insert records")
		}
		log.Println("one batch finished, len=", len(records))
	}
}

func parseArgs() (string, string, string) {
	filePath := flag.String("f", "./pp-monthly-update-new-version.csv", "CSV data to load")
	databaseDSN := flag.String("d", "root:password@tcp(127.0.0.1:3306)/dbench", "Database DSN")
	elasticsearchDSN := flag.String("e", "http://103.3.60.74:9200", "Elasticsearch DSN")
	flag.Parse()
	return *filePath, *databaseDSN, *elasticsearchDSN
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
