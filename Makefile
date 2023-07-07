all:
	go build ./quickws/autobahn-quickws.go
	go build -o test-client ./client/client.go

	GOOS=linux GOARCH=amd64 go build -o autobahn-quickws-linux ./quickws/autobahn-quickws.go
	GOOS=linux GOARCH=amd64 go build -o test-client-linux ./client/client.go
