package main

import (
	"database/sql"
	"embed"
	"html/template"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

var db, err = sql.Open("sqlite3", "metrics.db")

func checkErr(err error) {
	if err != nil {
			panic(err)
	}
}

func readFromDatabase() {
	rows, err := db.Query("SELECT * FROM hello")
	checkErr(err)
	var uid int
	var name string

	for rows.Next() {
			err = rows.Scan(&uid, &name)
			checkErr(err)
			fmt.Println(uid)
			fmt.Println(name)
	}

	rows.Close()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		readFromDatabase();
		data := map[string]string{
			"Region": os.Getenv("FLY_REGION"),
		}

		t.ExecuteTemplate(w, "index.html.tmpl", data)
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}