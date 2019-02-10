package convert

import (
	"database/sql"
	"fmt"
	"github.com/jszwec/csvutil"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
)

type Record struct {
	Id          string
	Title       string
	Body        string
	Tags_string string
}

func (*Record) CreateRecord(id, title, body, tags_string string) Record {
	return Record{Id: id, Title: title, Body: body, Tags_string: tags_string}
}
func (r Record) ToStringArray() []string {
	return []string{r.Id, r.Title, r.Body, r.Tags_string}
}
func (r Record) ToInsertString() string {
	sql := "insert into my_blog(id,title,body,tags_string)"
	sql += fmt.Printf(" values(%s,'%s','%s','%s')",
		r.Id,
		r.Title,
		r.Body,
		r.Tags_string)
	return sql
}

func ConvertCSV(dbname, csvPath string) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(
		`SELECT id, title, body, tags_string FROM my_blog`,
	)
	if err != nil {
		panic(err)
	}
	var recs []*Record
	defer rows.Close()
	for rows.Next() {
		rec := new(Record)
		if err := rows.Scan(&rec.Id, &rec.Title, &rec.Body, &rec.Tags_string); err != nil {
			log.Fatal("rows.Scan()", err)
			return
		}
		recs = append(recs, rec)
	}
	b, err := csvutil.Marshal(recs)
	if err != nil {
		fmt.Println("error:", err)
	}
	file, err := os.Create(csvPath)
	if err != nil {
		// Openエラー処理
	}
	defer file.Close()
	file.Write(b)
}
func ConvertDB(dbname, csvPath string) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		panic(err)
	}
	// CSVを構造体配列に変換する。
	data, err := ioutil.ReadFile(csvPath)
	if err != nil {
		panic(err)
	}
	var recs []*Record
	if err := csvutil.Unmarshal(data, &recs); err != nil {
		fmt.Println("error:", err)
	}
	// 構造体配列をSQLに変換し、対象DBへ挿入。
	for _, rec := range recs {
		db.Exec(rec.ToInsertString())
	}
}
