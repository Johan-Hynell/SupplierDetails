package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Product struct {
	ProductName  string
	EAN          string
	Details      string
	PricePerUnit string
	Currency     string
	ISO4217      int16
	Unit         string
}

type Supplier struct {
	SupplierName     string
	SupplierDetails  string
	PEPPOLEndpointID string
	Country          string
	City             string
	Street           string
	Postcode         string
	ProductList      []Product
}

func supplierJSON(s Supplier) string {
	var b, _ = json.Marshal(s)
	return string(b)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	var sup Supplier
	sup.SupplierName = "Example Supplier"
	sup.SupplierDetails = "An example supplier of products"
	sup.PEPPOLEndpointID = "0"
	sup.Country = "Sweden"
	sup.City = "Lulea"
	sup.Street = "Luleå University of Technology"
	sup.Postcode = "SE-97187"
	var p Product
	p.ProductName = "Example Product"
	p.Details = "An example product to order"
	p.EAN = "0"
	p.PricePerUnit = "123.45"
	p.Currency = "SEK"
	p.ISO4217 = 752
	p.Unit = "EA"
	sup.ProductList = append(sup.ProductList, p)

	fmt.Fprint(w, supplierJSON(sup))
}

func main() {
	http.HandleFunc("/info", testHandler)
	log.Fatal(http.ListenAndServe(":934", nil))
}
