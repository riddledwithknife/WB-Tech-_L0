package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/nats-io/stan.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func subscriptionHandler(db *gorm.DB) stan.MsgHandler {
	return func(msg *stan.Msg) {
		var order Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Panicln("failed to unmarshal order: ", err)
		}

		lock.Lock()
		if err := db.Create(&order).Error; err != nil {
			defer func() {
				if r := recover(); r != nil {
					log.Panicln("failed to create order: ", err)
				}
			}()
		}
		lock.Unlock()
	}
}

type OrderCache struct {
	cache map[string]string
	mu    sync.RWMutex
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		cache: make(map[string]string),
	}
}

func (o *OrderCache) Set(id string, data string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.cache[id] = data
}

func (o *OrderCache) Get(id string) (string, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	data, ok := o.cache[id]
	return data, ok
}

func getOrderHandler(cache *OrderCache, db *gorm.DB) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")

		if jsonData, ok := cache.Get(id); ok {
			wr.Header().Set("Content-Type", "application/json")
			wr.Write([]byte(jsonData))
			return
		}

		var result Order
		lock.Lock()
		if err := db.Preload("Delivery").Preload("Payment").Preload("Items").First(&result, "order_uid = ?", id).Error; err != nil {
			http.Error(wr, "Order not found", http.StatusNotFound)
		}
		lock.Unlock()

		jsonBytes, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
		}

		cache.Set(id, string(jsonBytes))

		wr.Header().Set("Content-Type", "application/json")
		wr.Write(jsonBytes)
	}
}

var (
	lock sync.Mutex
)

func main() {
	db, err := gorm.Open(postgres.Open("host=localhost dbname=wb_service port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Panicln("failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&Order{}, &Delivery{}, &Payment{}, &Item{})
	if err != nil {
		log.Panicln("failed to migrate db: ", err)
	}

	sc, err := stan.Connect("test-cluster", "order-service")
	if err != nil {
		log.Panicln("Can't connect to cluster: ", err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe("orders", subscriptionHandler(db), stan.DurableName("order-service"))
	if err != nil {
		log.Panicln("failed to subscribe to order: ", err)
	}
	defer sub.Close()

	cache := NewOrderCache()
	http.HandleFunc("/order", getOrderHandler(cache, db))
	log.Panicln(http.ListenAndServe(":8080", nil))
}
