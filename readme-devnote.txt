#Next / TODO 

- Pass json struct [] for multiple emails




#start fake smtp server
https://github.com/ReachFive/fake-smtp-server
su admin
npm install -g fake-smtp-server
exit
fake-smtp-server

# test 
go build gomail.go
./gomail


#Edit config in config.yml





# To learn : CHeck arrow json schema

// if err := readEftJSON("eft-test1.json", &eftInfos); err != nil {
// 	log.Fatal(err)
// }
// schema := arrow.NewSchema(
// 	[]arrow.Field{
// 		{Name: "email", Type: arrow.BinaryTypes.String},
// 		{Name: "customerName", Type: arrow.BinaryTypes.String},
// 		{Name: "bankAccountNumber", Type: arrow.BinaryTypes.String},
// 		{Name: "invoiceNumber", Type: arrow.BinaryTypes.String},
// 		{Name: "amount", Type: arrow.PrimitiveTypes.Float32},
// 		{Name: "paymentDate", Type: arrow.BinaryTypes.String},
// 	},
// 	nil,
// )