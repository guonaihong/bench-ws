all:
	# 编译mac
	GOOS=darwin	GOARCH=arm64 go build -o quickws.mac ./quickws/quickws.go
	GOOS=darwin	GOARCH=arm64 go build -o gws.mac ./gws/gws.go
	GOOS=darwin	GOARCH=arm64 go build -o gws-std.mac ./gws-std/gws-std.go
	GOOS=darwin	GOARCH=arm64 go build -o gorilla.mac ./gorilla/gorilla.go
	GOOS=darwin	GOARCH=arm64 go build -o gobwas.mac ./gobwas/gobwas.go
	GOOS=darwin	GOARCH=arm64 go build -o nettyws.mac ./nettyws/nettyws.go
	GOOS=darwin	GOARCH=arm64 go build -o nbio-std.mac ./nbio-std/nbio-std.go
	GOOS=darwin	GOARCH=arm64 go build -o nbio-nonblocking.mac ./nbio-nonblocking/nbio-nonblocking.go
	GOOS=darwin	GOARCH=arm64 go build -o test-client.mac ./client/client.go

	# 编译linux
	GOOS=linux GOARCH=amd64 go build -o gws.linux ./gws/gws.go
	GOOS=linux GOARCH=amd64 go build -o gws-std.linux ./gws-std/gws-std.go
	GOOS=linux GOARCH=amd64 go build -o quickws.linux ./quickws/quickws.go
	GOOS=linux GOARCH=amd64 go build -o gorilla.linux ./gorilla/gorilla.go
	GOOS=linux GOARCH=amd64 go build -o gobwas.linux  ./gobwas/gobwas.go
	GOOS=linux GOARCH=amd64 go build -o nettyws.linux ./nettyws/nettyws.go
	GOOS=linux GOARCH=amd64 go build -o nbio-std.linux ./nbio-std/nbio-std.go
	GOOS=linux GOARCH=amd64 go build -o nbio-nonblocking.linux ./nbio-nonblocking/nbio-nonblocking.go
	GOOS=linux GOARCH=amd64 go build -o test-client.linux ./client/client.go

clean:
	rm *.linux
	rm *.mac

