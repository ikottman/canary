package main

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

var db, err = sql.Open("sqlite3", "measurements.db")

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readFromDatabase() {
	rows, err := db.Query("SELECT * FROM measurements")
	checkErr(err)
	var temperature float32
	var pressure float32
	var humidity float32
	var gas_resistance int
	var iaq float32
	var accuracy int
	var co2_equivalent float32
	var voc_estimate float32
	var created_at string
	for rows.Next() {
		err = rows.Scan(
			&temperature,
			&pressure,
			&humidity,
			&gas_resistance,
			&iaq,
			&accuracy,
			&co2_equivalent,
			&voc_estimate,
			&created_at,
		)
		checkErr(err)
		fmt.Println(temperature)
	}

	rows.Close()
}

func recordMeasurement() {
	const query = `
		INSERT INTO measurements(
			temperature,
			pressure,
			humidity,
			gas_resistance,
			iaq,
			accuracy,
			co2_equivalent,
			voc_estimate,
			created_at
		)
		VALUES(?,?,?,?,?,?,?,?,?)
	`
	statement, err := db.Prepare(query)
	checkErr(err)
	var _, error = statement.Exec(
		24.65,
		1019.42,
		46.23,
		597617,
		26.9,
		2,
		506.7,
		0.51,
		time.Now().Format("2006-01-02T15:04:05Z"),
	)
	checkErr(error)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		readFromDatabase()
		data := map[string]string{
			"Region": os.Getenv("FLY_REGION"),
		}

		t.ExecuteTemplate(w, "index.html.tmpl", data)
	})

	http.HandleFunc("/measurement", func(w http.ResponseWriter, r *http.Request) {
		recordMeasurement()
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
