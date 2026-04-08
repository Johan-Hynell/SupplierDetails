package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"


	"database/sql"
	_ "modernc.org/sqlite"
)

type Product struct {
	ProductName	 string
	ProductID	 int64
	EAN       	 string
	Details      string
	PricePerUnit string
	Currency     string
	ISO4217      int64
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
var sup Supplier
var productDatabase *sql.DB

var testProduct Product

func testProductInit()
{
	testProduct.ProductName = "Example Product"
	testProduct.Details = "An example product to order"
	testProduct.EAN = "0"
	testProduct.PricePerUnit = "123.45"
	testProduct.Currency = "SEK"
	testProduct.ISO4217 = 752
	testProduct.Unit = "EA"
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	UpdateList(productDatabase,&sup)
	fmt.Fprint(w, supplierJSON(sup))
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "error parsing request")
		return
	}
	fmt.Println(r.PostForm.Get("pid"))
	pid, pidErr := strconv.ParseInt(r.PostForm.Get("pid"),10,64)
	if pidErr != nil {
		fmt.Println(pidErr)
		fmt.Fprint(w, "error parsing Product ID, must be an integer")
		return
	}
	iso4217,  iso4217err := strconv.ParseInt(r.PostForm.Get("iso4217"),10,64)
	if iso4217err != nil {
		fmt.Fprint(w, "error parsing ISO4217, must be an integer")
		return
	}
	productToAdd := Product{
		ProductName: r.PostForm.Get("name"),
		ProductID: pid,
		EAN: r.PostForm.Get("ean"),
		PricePerUnit: r.PostForm.Get("ppu"),
		Details: r.PostForm.Get("details"),
		Currency: r.PostForm.Get("currency"),
		ISO4217: iso4217,
		Unit: r.PostForm.Get("unit"),
	}

	err = addProduct(productDatabase, productToAdd)
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w,"error adding product")
	} else {
		fmt.Fprint(w,"added")
	}
}

func addProductFormHandler(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("addProductForm.html")
	if err != nil {
		fmt.Println("error gettig form: ", err)
		fmt.Fprint(w, "error getting form")
	} else {
		html := string(b)
		fmt.Fprint(w, html)
	}
}

func main() {
	sup.SupplierName = "Example Supplier"
	sup.SupplierDetails = "An example supplier of products"
	sup.PEPPOLEndpointID = "0"
	sup.Country = "Sweden"
	sup.City = "Lulea"
	sup.Street = "Luleå University of Technology"
	sup.Postcode = "SE-97187"
	
	

	//sup.ProductList = append(sup.ProductList, p)
	productDB, err, closeDB := openDB()
	if err != nil {
		log.Fatalf("Error opening the database: %v\n", err)
	}
	defer closeDB()
	productDatabase = productDB
	http.HandleFunc("/info", testHandler)
	http.HandleFunc("/add", addProductHandler)
	http.HandleFunc("/addForm", addProductFormHandler)
	log.Fatal(http.ListenAndServe(":934", nil))
	
}

func openDB() (*sql.DB, error, func()) {
	dbPath := "products.db"

	// Check if the file exists
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		fmt.Println("database does not exist, attempting to create it...")
	} else if err != nil {
		fmt.Println("error checking database file:", err)
		return nil, err, nil
	}

	// Try to open (or create) the database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println("error opening/creating database:", err)
		return nil, err, nil
	}

	// Create table if it does not exist
	if err := CreateTableIfNotExists(db); err != nil {
		fmt.Println("error creating the product table:", err)
		return nil, err, nil
	}
	fmt.Println("Database Ready")

	// Return the database and a cleanup function to close it
	return db, nil, func() {
		db.Close()
		log.Println("closing the service registry database connection")
	}
}

// CreateTableIfNotExists checks if the table exists and creates it if it does not.
func CreateTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS Products (
		Name TEXT NOT NULL,
		ID INTEGER PRIMARY KEY,
		EAN TEXT,
		Details TEXT,
		PricePerUnit TEXT,
		Currency TEXT,
		ISO4217 INTEGER,
		Unit TEXT
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func UpdateList(db *sql.DB, sup *Supplier) error {
	query := `SELECT * FROM Products;`
	result, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error fetching products: %w", err)
	}
	sup.ProductList = nil
	for result.Next() {
		var prod Product
		err := result.Scan(&prod.ProductName, &prod.ProductID,
			&prod.EAN, &prod.Details, &prod.PricePerUnit,
			&prod.Currency, &prod.ISO4217, &prod.Unit)
		if(err != nil) {
			return fmt.Errorf("error fetching products: %w", err)
		}
		sup.ProductList = append(sup.ProductList, prod)
	}	
	return nil
}

func addProduct(db *sql.DB, prod Product) error {
	query := `
	INSERT INTO Products (
		Name, ID, EAN, Details, PricePerUnit, Currency, ISO4217, Unit
	) VALUES (?,?,?,?,?,?,?,?);
	`
	_, err := db.Exec(query, prod.ProductName, prod.ProductID,
			prod.EAN, prod.Details, prod.PricePerUnit,
			prod.Currency, prod.ISO4217, prod.Unit)
	if err != nil {
		return fmt.Errorf("error adding product: %w", err)
	}
	return nil
}