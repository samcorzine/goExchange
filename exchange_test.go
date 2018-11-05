package main

import (
	"testing"
)

func TestMoveChecker(t *testing.T) {
	testOrderBook := initOrderBook()
	testOrderBook.addOrder(Order{UUID: "sdfaD", Price: 23545.0, ContractType: "Bid", OrderID: 1})
	testOrderBook.addOrder(Order{UUID: "sdfafsD", Price: 23545.0, ContractType: "Bid", OrderID: 3})
	testOrderBook.addOrder(Order{UUID: "sdfaD", Price: 23545.0, ContractType: "Bid", OrderID: 5})
	testOrderBook.addOrder(Order{UUID: "sdffdsaD", Price: 2345.0, ContractType: "Bid", OrderID: 7})
	// spew.Dump(testOrderBook)

	if testOrderBook.numBids() != 4 {
		t.Errorf("bid wasn't added")
	}
	testOrderBook.addOrder(Order{UUID: "ssadfasdf", Price: 342, ContractType: "Ask", OrderID: 2})
	testOrderBook.addOrder(Order{UUID: "ssaddf", Price: 342, ContractType: "Ask", OrderID: 4})
	testOrderBook.addOrder(Order{UUID: "ssaddf", Price: 34, ContractType: "Ask", OrderID: 4})
	if testOrderBook.numAsks() != 3 {
		t.Errorf("ask wasn't added")
	}
	if testOrderBook.clear() != true {
		t.Errorf("clear didn't occur when it should have")
	}
	if testOrderBook.clear() != false {
		t.Errorf("clear occured when it shouldn't have")
	}
	if testOrderBook.numOrders() != 1 {
		t.Errorf("orders weren't removed after clear")
	}
}
