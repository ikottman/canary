package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ikottman/canary/auth"
)

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

func recordMeasurement(token string) {
	body := &Measurement{
		Temperature:   3.15,
		Pressure:      1019.38,
		Humidity:      45.64,
		GasResistance: 597617,
		IAQ:           28.3,
		Accuracy:      3,
		CO2:           511.21,
		VOC:           0.52,
	}
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(body)

	var req, _ = http.NewRequest("POST", "https://canary.ikottman.com/measurement", payloadBuffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	var _, err = http.DefaultClient.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

func main() {
	var token, err = auth.CreateJwt()
	if err != nil {
		log.Fatalln(err)
	}

	recordMeasurement(token)
}
