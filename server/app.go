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
	Temperature   float32 `json:"temperature"`
	Pressure      float32 `json:"pressure"`
	Humidity      float32 `json:"humidity"`
	GasResistance int     `json:"gasResistance"`
	IAQ           float32 `json:"IAQ"`
	Accuracy      int     `json:"iaqAccuracy"`
	CO2           float32 `json:"eqCO2"`
	VOC           float32 `json:"eqBreathVOC"`
	CreatedAt     string  `json:createdAt`
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

func readFromDatabase(limit int) []Measurement {
	rows, err := db.Query(`
		SELECT
			temperature,
			pressure,
			humidity,
			gas_resistance,
			iaq,
			accuracy,
			co2_equivalent,
			voc_estimate,
			created_at
		FROM measurements
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	checkErr(err)
	var measurements []Measurement
	for rows.Next() {
		var m Measurement
		err = rows.Scan(
			&m.Temperature,
			&m.Pressure,
			&m.Humidity,
			&m.GasResistance,
			&m.IAQ,
			&m.Accuracy,
			&m.CO2,
			&m.VOC,
			&m.CreatedAt,
		)
		measurements = append(measurements, m)
		checkErr(err)
	}
	rows.Close()
	return measurements
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
func getAccuracy(accuracy int) string {
	if accuracy == 0 {
		return "invalid"
	} else if accuracy == 1 {
		return "good"
	} else if accuracy == 2 {
		return "better"
	} else {
		return "best"
	}
}

type TemplateData struct {
	Measurements []Measurement
	Accuracy     string
}

// route handlers
func index(w http.ResponseWriter, r *http.Request) {
	var measurements = readFromDatabase(120)
	data := TemplateData{
		Measurements: measurements,
		Accuracy:     getAccuracy(measurements[len(measurements)-1].Accuracy),
	}
	t.ExecuteTemplate(w, "index.html.tmpl", data)
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

	http.HandleFunc("/", index)
	http.HandleFunc("/measurement", authenticated(measurement))
	http.HandleFunc("/reset", authenticated(reset))

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
