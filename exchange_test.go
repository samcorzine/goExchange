package main

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMoveChecker(t *testing.T) {
	testOrderBook := initOrderBook()
	testOrderBook.addOrder(Order{UUID: "sdfaD", Price: 2352345.0, ContractType: "Bid", OrderID: 1})
	spew.Dump(testOrderBook)
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
