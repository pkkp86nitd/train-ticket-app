# train-ticket-app


References- :\
 1.https://grpc.io/docs/protoc-installation/\
 2.https://grpc.io/docs/languages/go/quickstart/\

Step before running apps (client,server,server_test)\
   -> protoc --go_out=. --go-grpc_out=. proto/train.proto (FROM ROOT DIRECTORY to auto generate proto files ) \

Commands to run the client,server,server_test \
  -> go run  client/main.go   \
  -> go run server/main.go \
  -> go test testing/server_test.go   
       

  
