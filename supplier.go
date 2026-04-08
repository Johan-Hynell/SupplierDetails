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

var supplierGlobal Supplier
var productDatabase *sql.DB
/*var testProduct Product
func testProductInit() {
	testProduct.ProductName = "Example Product"
	testProduct.Details = "An example product to order"
	testProduct.EAN = "0"
	testProduct.PricePerUnit = "123.45"
	testProduct.Currency = "SEK"
	testProduct.ISO4217 = 752
	testProduct.Unit = "EA"
}*/

type ServerConfig struct {
	SupplierInfo Supplier
	Port int64
	AllowAdd bool
	FormatJSON bool
}
var serverConfig ServerConfig

func main() {
	configHandler()
	//supplierGlobal.ProductList = append(supplierGlobal.ProductList, p)
	productDB, err, closeDB := openDB()
	if err != nil {
		log.Fatalf("Error opening the database: %v\n", err)
	}
	defer closeDB()
	productDatabase = productDB
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/add", addProductHandler)
	http.HandleFunc("/addForm", addProductFormHandler)
	fmt.Printf("Using port: %d\n",serverConfig.Port)
	fmt.Printf("Allow adding products through http: %t\n", serverConfig.AllowAdd)
	fmt.Printf("Format json output (tabs and newlines): %t\n", serverConfig.FormatJSON)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",serverConfig.Port), nil))
}

//HTTP HANDLERS----------------------------------------------------------------
func infoHandler(w http.ResponseWriter, r *http.Request) {
	UpdateList(productDatabase,&supplierGlobal)
	var b []byte
	if serverConfig.FormatJSON {
		b, _ = json.MarshalIndent(supplierGlobal,"","\t")
	} else {
		b, _ = json.Marshal(supplierGlobal)
	}
	
	fmt.Fprint(w, string(b))
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
	if !serverConfig.AllowAdd {
		fmt.Fprint(w, "Adding products through http is disabled")
		return
	}
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

//CONFIG HANDLING--------------------------------------------------------------
func configHandler() error{
	
	// Check if config exists
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		fmt.Println("config.json does not exist, attempting to create it...")
		createConfig()
	} else if err != nil {
		fmt.Errorf("error checking config.json: %w", err)
		return err
	}
	// Get config contents
	var b []byte
	b, err = os.ReadFile("config.json")
	if err != nil {
		fmt.Errorf("error reading config.json: %w", err)
		return err
	}
	err = json.Unmarshal(b, &serverConfig)
	if err != nil {
		fmt.Errorf("error parsing config.json: %w", err)
		return err
	}
	supplierGlobal = serverConfig.SupplierInfo
	return nil
}

func createConfig() error {
	//Defaults
	supplierGlobal.SupplierName = "Example Supplier"
	supplierGlobal.SupplierDetails = "An example supplier of products"
	supplierGlobal.PEPPOLEndpointID = "0"
	supplierGlobal.Country = "Sweden"
	supplierGlobal.City = "Lulea"
	supplierGlobal.Street = "Luleå University of Technology"
	supplierGlobal.Postcode = "SE-97187"
	serverConfig.SupplierInfo = supplierGlobal
	serverConfig.AllowAdd = false
	serverConfig.Port = 934
	serverConfig.FormatJSON = false
	//Make json
	b, berr := json.MarshalIndent(serverConfig,"","\t")
	fmt.Println(string(b))
	if berr != nil {
		fmt.Errorf("error making json: %w", berr)
		return berr
	}
	//Create config with defaults
	err := os.WriteFile("config.json", b, 0644)
	if err != nil {
		fmt.Errorf("error creating config.json: %w", err)
		return err
	}
	return nil
}

//DATABASE---------------------------------------------------------------------
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

func UpdateList(db *sql.DB, supplierGlobal *Supplier) error {
	query := `SELECT * FROM Products;`
	result, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error fetching products: %w", err)
	}
	supplierGlobal.ProductList = nil
	for result.Next() {
		var prod Product
		err := result.Scan(&prod.ProductName, &prod.ProductID,
			&prod.EAN, &prod.Details, &prod.PricePerUnit,
			&prod.Currency, &prod.ISO4217, &prod.Unit)
		if(err != nil) {
			return fmt.Errorf("error fetching products: %w", err)
		}
		supplierGlobal.ProductList = append(supplierGlobal.ProductList, prod)
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