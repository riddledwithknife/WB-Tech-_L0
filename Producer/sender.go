package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
)

type Order struct {
	OrderUID          string    `json:"order_uid" gorm:"primarykey"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery" gorm:"foreignKey:OrderUID"`
	Payment           Payment   `json:"payment" gorm:"foreignKey:OrderUID"`
	Items             []Item    `json:"items" gorm:"foreignKey:OrderUID"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	ID       uint   `json:"-" gorm:"primarykey"`
	OrderUID string `json:"-"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
	ID           uint   `json:"-" gorm:"primarykey"`
	OrderUID     string `json:"-"`
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ID          uint   `json:"-" gorm:"primarykey"`
	OrderUID    string `json:"-"`
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	sc, err := stan.Connect("test-cluster", "order-Producer", stan.NatsURL("nats://nats-streaming:4222"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to nats cluster")
	defer sc.Close()

	data, _ := os.ReadFile("./order_example.json")

	var order Order
	json.Unmarshal(data, &order)

	for i := range 20 {
		order.OrderUID = strconv.Itoa(i)
		jsonBytes, _ := json.Marshal(order)

		sc.Publish("orders", jsonBytes)

		log.Println("Published order: ", order.OrderUID)
	}

	fmt.Println("Done Publishing")
}
