### Next / TODO 
- Decide: should I use infra for logs, or just a file?
- Decide: should I go for orchestrator?


## Quick start
### start fake smtp server
https://github.com/ReachFive/fake-smtp-server
```
su admin
npm install -g fake-smtp-server
exit
fake-smtp-server  --http-ip 127.0.0.1
```
http://localhost:1080

# test 
```
go build
go-mail -configFile /my/secure/location/eftconfig/config-prod.yml
```

### Edit config in config.yml
See sample config.yml





#### To learn : CHeck arrow json schema
```go
 if err := readEftJSON("eft-test1.json", &eftInfos); err != nil {
 	log.Fatal(err)
 }
 schema := arrow.NewSchema(
 	[]arrow.Field{
 		{Name: "email", Type: arrow.BinaryTypes.String},
 		{Name: "customerName", Type: arrow.BinaryTypes.String},
 		{Name: "bankAccountNumber", Type: arrow.BinaryTypes.String},
 		{Name: "invoiceNumber", Type: arrow.BinaryTypes.String},
 		{Name: "amount", Type: arrow.PrimitiveTypes.Float32},
 		{Name: "paymentDate", Type: arrow.BinaryTypes.String},
 	},
 	nil,
 )
```