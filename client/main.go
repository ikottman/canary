package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ikottman/canary/auth"
	"github.com/ikottman/canary/cat"
)

type Measurement struct {
	Temperature   float32 `json:"temperature"`
	Pressure      float32 `json:"pressure"`
	Humidity      float32 `json:"humidity"`
	GasResistance int     `json:"gasResistance"`
	IAQ           float32 `json:"IAQ"`
	Accuracy      int     `json:"iaqAccuracy"`
	CO2           float32 `json:"eqCO2"`
	VOC           float32 `json:"eqBreathVOC"`
}

func recordMeasurement(measurement string, token string) {
	if len(measurement) == 0 {
		return
	}
	// parse measurement to struct
	body := Measurement{}
	json.Unmarshal([]byte(measurement), &body)

	// convert struct to buffer
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

	cat.Cat("/dev/ttyACM0", token, recordMeasurement)
}
