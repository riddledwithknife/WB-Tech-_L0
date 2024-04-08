package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

func main() {
	sc, err := stan.Connect("test-cluster", "order-producer")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	data, err := os.ReadFile("./order.json")
	if err != nil {
		log.Panicln("Error reading file:", err)
	}

	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Panicln("Error parsing json:", err)
	}

	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Panicln("Error converting json into bytes:", err)
	}

	err = sc.Publish("orders", jsonDataBytes)
	if err != nil {
		log.Panicln("Error publishing:", err)
	}

	fmt.Println("Published order successfully")
}
