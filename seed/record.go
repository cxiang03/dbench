package dbench

import (
	"strconv"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type RawRecord struct {
	UUID      string
	Price     string
	TimeStamp string
	PostCode  string
	PType     string
	IsNew     string
	Duration  string
	Addr1     string
	Addr2     string
	Street    string
	Locality  string
	Town      string
	District  string
	County    string
	D         string
	E         string
}

type Record struct {
	bun.BaseModel `bun:"table:prices"`
	ID            int64  `bun:",pk,autoincrement"`
	UUID          string `json:"uuid"`
	Price         int64  `json:"price"`
	TimeStamp     int64  `json:"time_stamp"`
	PostCode      string `json:"postcode"`
	PType         string `json:"p_type"`
	IsNew         bool   `json:"is_new"`
	Duration      string `json:"duration"`
	Addr1         string `json:"addr1"`
	Addr2         string `json:"addr2"`
	Street        string `json:"street"`
	Locality      string `json:"locality"`
	Town          string `json:"town"`
	District      string `json:"district"`
	County        string `json:"county"`
}

func ParseRecord(line []string) *Record {
	raw := RawRecord{
		UUID:      line[0],
		Price:     line[1],
		TimeStamp: line[2],
		PostCode:  line[3],
		PType:     line[4],
		IsNew:     line[5],
		Duration:  line[6],
		Addr1:     line[7],
		Addr2:     line[8],
		Street:    line[9],
		Locality:  line[10],
		Town:      line[11],
		District:  line[12],
		County:    line[13],
		// D:         line[14],
		// E:         line[15],
	}
	uuid := strings.TrimPrefix(raw.UUID, "{")
	uuid = strings.TrimSuffix(uuid, "}")
	price, _ := strconv.Atoi(raw.Price)
	timestamp, _ := time.Parse("2006-01-02 15:04", raw.TimeStamp)
	return &Record{
		UUID:      uuid,
		Price:     int64(price),
		TimeStamp: timestamp.Unix(),
		PostCode:  truncate(raw.PostCode, 10),
		PType:     truncate(raw.PType, 10),
		IsNew:     raw.IsNew == "Y",
		Duration:  truncate(raw.Duration, 10),
		Addr1:     truncate(raw.Addr1, 100),
		Addr2:     truncate(raw.Addr2, 100),
		Street:    truncate(raw.Street, 100),
		Locality:  truncate(raw.Locality, 100),
		Town:      truncate(raw.Town, 100),
		District:  truncate(raw.District, 100),
		County:    truncate(raw.County, 100),
	}
}

func truncate(s string, length int) string {
	if len(s) < length {
		return s
	}
	return s[:length]
}
