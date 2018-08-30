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

type BidsAndAsks struct {
	Bids  []Order
	Asks  []Order
	Price float64
	sync.Mutex
}

type OrderBook struct {
	// TODO: rewrite everything as a slice of BidsAndAsks with the price inside those structs and the book sorted with that as a key
	pricePoints []BidsAndAsks
	sync.Mutex
}

func (bidsasks *BidsAndAsks) addOrder(ord Order) {
	if ord.ContractType == "Bid" {
		if len(bidsasks.Bids) == 0 {
			bidsasks.Bids = []Order{ord}
		} else {
			bidsasks.Bids = append(bidsasks.Bids, ord)
		}
	}
	if ord.ContractType == "Ask" {
		if len(bidsasks.Asks) == 0 {
			bidsasks.Asks = []Order{ord}
			return
		} else {
			bidsasks.Asks = append(bidsasks.Asks, ord)
			return
		}
	}
}

func (book *OrderBook) addOrder(ord Order) {
	if len(book.pricePoints) == 0 {
		var newPricePoint BidsAndAsks
		newPricePoint.addOrder(ord)
		newPricePoint.Price = ord.Price

		book.pricePoints = []BidsAndAsks{newPricePoint}
		return
	}
	for i, x := range book.pricePoints {
		if x.Price == ord.Price {
			book.pricePoints[i].addOrder(ord)
			return
		}

		if ord.Price < x.Price {
			var newPricePoint BidsAndAsks
			newPricePoint.addOrder(ord)
			newPricePoint.Price = ord.Price
			if i == 0 {
				tmp := []BidsAndAsks{newPricePoint}
				for _, x := range book.pricePoints {
					tmp = append(tmp, x)
				}
				book.pricePoints = tmp
				return
			}
			if len(book.pricePoints) > 1 {
				var tmp []BidsAndAsks
				copy(tmp, book.pricePoints[i:])
				book.pricePoints = append(book.pricePoints[:i-1], newPricePoint)
				book.pricePoints = append(book.pricePoints, tmp...)
			} else {
				newStart := []BidsAndAsks{newPricePoint}
				newStart = append(newStart, book.pricePoints[0])
				book.pricePoints = newStart
			}
		}
	}
}

func (book *OrderBook) clear() bool {
	lastSucceeded := true
	didAClear := false
	for lastSucceeded == true {
		foundAsk := false
		foundBid := false
		counter := 0
		bidPlace := -1
		askPlace := -1
		for foundAsk == false && counter < len(book.pricePoints) {
			if len(book.pricePoints[counter].Asks) != 0 {
				foundAsk = true
				askPlace = counter
				for foundBid == false {
					if len(book.pricePoints[counter].Bids) != 0 {
						foundBid = true
						bidPlace = counter
					} else {
						counter += 1
					}
				}
			} else {
				counter += 1
			}
		}
		if foundAsk && foundBid {
			didAClear = true
			if len(book.pricePoints[bidPlace].Bids) > 1 {
				book.pricePoints[bidPlace].Bids = book.pricePoints[bidPlace].Bids[:1]
			} else {
				book.pricePoints[bidPlace].Bids = []Order{}
			}
			if len(book.pricePoints[bidPlace].Asks) > 1 {
				book.pricePoints[askPlace].Asks = book.pricePoints[askPlace].Asks[:1]
			} else {
				book.pricePoints[askPlace].Asks = []Order{}
			}
		} else {
			lastSucceeded = false
		}
	}
	return didAClear
}

func (book OrderBook) numBids() int {
	var numBids int
	for _, v := range book.pricePoints {
		for range v.Bids {
			numBids += 1
		}
	}
	return numBids
}
func (book OrderBook) numAsks() int {
	var numAsks int
	for _, v := range book.pricePoints {
		for range v.Asks {
			numAsks += 1
		}
	}
	return numAsks
}
func (book OrderBook) numOrders() int {
	var numOrders int
	for _, v := range book.pricePoints {
		for range v.Asks {
			numOrders += 1
		}
		for range v.Bids {
			numOrders += 1
		}

	}
	return numOrders
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

func initOrderBook() *OrderBook {
	return &OrderBook{pricePoints: make([]BidsAndAsks, 0)}
}

func main() {
	testOrderBook := initOrderBook()
	launchHTTPAPI(testOrderBook)
	for {
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Println("Before Clear:")
		spew.Dump(testOrderBook)
		// testOrderBook.clear()
		fmt.Println("After Clear:")
		spew.Dump(testOrderBook)
	}

}
