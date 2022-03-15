package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

var db, err = sql.Open("sqlite3", "measurements.db")

type Measurement struct {
	Temperature   float32
	Pressure      float32
	Humidity      float32
	GasResistance int
	IAQ           float32
	Accuracy      int
	CO2           float32
	VOC           float32
}

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

func recordMeasurement(measurement Measurement) {

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
		measurement.Temperature,
		measurement.Pressure,
		measurement.Humidity,
		measurement.GasResistance,
		measurement.IAQ,
		measurement.Accuracy,
		measurement.CO2,
		measurement.VOC,
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
		requestBody, _ := ioutil.ReadAll(r.Body)
		var measurement Measurement
		json.Unmarshal(requestBody, &measurement)
		recordMeasurement(measurement)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		stmt, err := db.Prepare("delete from measurements where 1=1")
		checkErr(err)

		stmt.Exec()
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
