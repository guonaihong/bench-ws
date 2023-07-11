all:
	go build ./quickws/autobahn-quickws.go
	go build ./gws-async-write/gws-async-write-server.go
	go build ./gws-sync-write/gws-sync-write-server.go
	go build -o test-client ./client/client.go

	GOOS=linux GOARCH=amd64 go build -o autobahn-quickws-linux ./quickws/autobahn-quickws.go
	GOOS=linux GOARCH=amd64 go build -o test-client-linux ./client/client.go
