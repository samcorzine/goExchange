package main

import(
  "testing";
  "fmt"
)

func TestMoveChecker(t *testing.T) {
    testOrderBook := OrderBook{Bids: BidBook{orders: []Order{}}, Asks: AskBook{orders: []Order{}}}
    testOrderBook.addOrder(Order{UUID: "sdfaD", Price: 2352345, ContractType: "Bid", OrderID: 1})
    fmt.Println(testOrderBook.Bids.orders)
    if testOrderBook.numBids() != 1 {
      t.Errorf("bid wasn't added")
    }
    testOrderBook.addOrder(Order{UUID: "ssadfasdf", Price: 342, ContractType: "Ask", OrderID: 2})
    if testOrderBook.numAsks() != 1 {
      t.Errorf("ask wasn't added")
    }
    if testOrderBook.clear() != true {
      t.Errorf("clear didn't occur when it should have")
    }
    if testOrderBook.clear() != false {
      t.Errorf("clear occured when it shouldn't have")
    }
}
