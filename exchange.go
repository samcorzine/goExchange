package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

type Order struct {
	UUID         string
	Price        float64
	ContractType string
	OrderID      int
}

type orderRequest struct {
	UUID         string  `json:"uuid"`
	Price        float64 `json:"price"`
	ContractType string  `json:"contracttype"`
}

type BidBook struct {
	orders []Order
	sync.Mutex
}
type AskBook struct {
	orders []Order
	sync.Mutex
}
type OrderBook struct {
	Bids BidBook
	Asks AskBook
}

func (book *BidBook) addOrder(ord Order) {
	if len(book.orders) == 0 {
		book.orders = []Order{ord}
		return
	}
	for i := 0; i < len(book.orders); i = i + 1 {
		currentEntry := book.orders[i]
		if ord.Price > currentEntry.Price {
			if i == 0 {
				book.orders = append([]Order{ord}, book.orders...)
			} else {
				var tmp []Order
				copy(tmp, book.orders[i:])
				book.orders = append(book.orders[:i-1], ord)
				book.orders = append(book.orders, tmp...)
			}
			return
		} else if i == len(book.orders)-1 {
			book.orders = append(book.orders, ord)
			return
		}
	}
}
func (book *AskBook) addOrder(ord Order) {
	if len(book.orders) == 0 {
		book.orders = []Order{ord}
		return
	}
	for i := 0; i < len(book.orders); i = i + 1 {
		currentEntry := book.orders[i]
		if ord.Price < currentEntry.Price {
			if i == 0 {
				book.orders = append([]Order{ord}, book.orders...)
			} else {
				var tmp []Order
				copy(tmp, book.orders[i:])
				book.orders = append(book.orders[:i-1], ord)
				book.orders = append(book.orders, tmp...)
			}
			return
		} else if i == len(book.orders)-1 {
			book.orders = append(book.orders, ord)
			return
		}
	}
}

func (book *OrderBook) addOrder(ord Order) {
	if ord.ContractType == "Ask" {
		book.Asks.Lock()
		defer book.Asks.Unlock()
		book.Asks.addOrder(ord)
	} else if ord.ContractType == "Bid" {
		book.Bids.Lock()
		defer book.Bids.Unlock()
		book.Bids.addOrder(ord)
	} else {
		log.Fatal("Order is of invalid type, cannot add to OrderBook")
	}
}

func (book *OrderBook) clear() bool {
	book.Asks.Lock()
	defer book.Asks.Unlock()
	book.Bids.Lock()
	defer book.Bids.Unlock()
	if len(book.Bids.orders) != 0 && len(book.Bids.orders) != 0 {
		highestBid := book.Bids.orders[0]
		lowestAsk := book.Asks.orders[0]
		if highestBid.Price > lowestAsk.Price {
			book.Bids.orders = book.Bids.orders[1:]
			book.Asks.orders = book.Asks.orders[1:]
			return true
		}
	}
	return false
}

func (book OrderBook) numBids() int {
	return len(book.Bids.orders)
}
func (book OrderBook) numAsks() int {
	return len(book.Asks.orders)
}
func (book OrderBook) numOrders() int {
	return book.numAsks() + book.numBids()
}

func launchHTTPAPI(book *OrderBook) {
	go func() {
		m := mux.NewRouter()
		m.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
			// vars := mux.Vars(r)
			// newVal, err := strconv.Atoi(vars["newVal"])
			var orderDetails orderRequest
			err := json.NewDecoder(r.Body).Decode(&orderDetails)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var order = Order{UUID: orderDetails.UUID, Price: orderDetails.Price, ContractType: orderDetails.ContractType, OrderID: rand.Int()}
			book.addOrder(order)
		}).Methods("POST")
		log.Fatal(http.ListenAndServe(":8080", m))
	}()
}

func main() {
	testOrderBook := OrderBook{Bids: BidBook{orders: []Order{}}, Asks: AskBook{orders: []Order{}}}
	launchHTTPAPI(&testOrderBook)
	for {
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Println("Before Clear:")
		spew.Dump(testOrderBook)
		// testOrderBook.clear()
		fmt.Println("After Clear:")
		spew.Dump(testOrderBook)
	}

}
