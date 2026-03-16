# Prototype Supplier information and product list provider
A simple service that provides information about you (the supplier) to buyers for easier order placement

## Notes
- Currently a randomly picked port (934) is used for testing before a standard port is set.
- Currently the path used to get the suppliers information is ```/info```

## Running
```go mod init github.com/Johan-Hynell/SupplierDetails```

```go run .``` may need to run as administrator/superuser

## Todo
- Database storing supplier information
- Read from config file

## Example output
```json
{
"SupplierName":"Example Supplier",
"SupplierDetails":"An example supplier of products",
"PEPPOLEndpointID":"0",
"Country":"Sweden","City":"Lulea","Street":"Luleå University of Technology","Postcode":"SE-97187","ProductList":[{"ProductName":"Example Product","EAN":"0","Details":"An example product to order","PricePerUnit":"123.45","Currency":"SEK","ISO4217":752,"Unit":"EA"}]}
```