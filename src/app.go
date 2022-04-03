package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ikottman/canary/auth"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

var db, _ = sql.Open("sqlite3", "measurements.db")

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

// check Authorization header contains a valid RS256 encoded JWT, signed with a private key which matches our CANARY_PUBLIC_KEY
// returns an empty 401 if authentication fails
func authenticated(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.ValidateJwt(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		f(w, r)
	}
}

func readFromDatabase() Measurement {
	rows, err := db.Query(`
		SELECT
			temperature,
			pressure,
			humidity,
			gas_resistance,
			iaq,
			accuracy,
			co2_equivalent,
			voc_estimate
		FROM measurements
		LIMIT 1
	`)
	checkErr(err)
	var measurement Measurement
	for rows.Next() {
		err = rows.Scan(
			&measurement.Temperature,
			&measurement.Pressure,
			&measurement.Humidity,
			&measurement.GasResistance,
			&measurement.IAQ,
			&measurement.Accuracy,
			&measurement.CO2,
			&measurement.VOC,
		)
		checkErr(err)
	}
	rows.Close()
	return measurement
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

// route handlers

func index(w http.ResponseWriter, r *http.Request) {
	var measurement = readFromDatabase()
	t.ExecuteTemplate(w, "index.html.tmpl", measurement)
}

func measurement(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var measurement Measurement
	json.Unmarshal(requestBody, &measurement)
	recordMeasurement(measurement)
}

func reset(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("delete from measurements where 1=1")
	checkErr(err)
	stmt.Exec()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	http.HandleFunc("/", authenticated(index))
	http.HandleFunc("/measurement", authenticated(measurement))
	http.HandleFunc("/reset", authenticated(reset))

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
