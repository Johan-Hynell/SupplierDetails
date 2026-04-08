# Prototype Supplier information and product list provider
A simple service that provides information about you (the supplier) to buyers for easier order placement

## Notes
- Currently a randomly picked port (934) is used for testing before a standard port is set.
- The default port and supplier details can be changed in config.json
- Currently the path used to get the suppliers information is ```/info```
- Modifying the product list can be done through http if ```AllowAdd``` is set to true in config or through accessing the database locally
- Server must be restarted to update config if config.json is changed

## Paths
- ```/info``` to get supplier details and product list
- ```/add``` send form in POST here to add product to list if ```AllowAdd``` is true
- ```/addForm``` A simple form to input product details, on submit sends request to ```/add```

## Running
```go mod init github.com/Johan-Hynell/SupplierDetails```

```go get modernc.org/sqlite``` (dependency)

```go run .``` may need to run as administrator/superuser

## Todo
- Code cleanup

## Example output
Real output does not contain linebreaks or tabs unless ```FormatJSON``` is set true in config
```json
{
    "SupplierName":"Example Supplier",
    "SupplierDetails":"An example supplier of products",
    "PEPPOLEndpointID":"0",
    "Country":"Sweden",
    "City":"Lulea",
    "Street":"Luleå University of Technology",
    "Postcode":"SE-97187",
    "ProductList":
    [{
        "ProductName":"Example Product",
        "EAN":"0",
        "Details":"An example product to order",
        "PricePerUnit":"123.45",
        "Currency":"SEK",
        "ISO4217":752,
        "Unit":"EA"
    }]
}
```