all:
	# 编译mac
	GOOS=darwin	GOARCH=arm64 go build -o quickws-std.mac ./quickws/quickws-std.go
	GOOS=darwin	GOARCH=arm64 go build -o gws.mac ./gws/gws.go
	GOOS=darwin	GOARCH=arm64 go build -o gws-std.mac ./gws-std/gws-std.go
	GOOS=darwin	GOARCH=arm64 go build -o test-client.mac ./client/client.go

	# 编译linux
	GOOS=linux GOARCH=amd64 go build -o gws.linux ./gws/gws.go
	GOOS=linux GOARCH=amd64 go build -o gws-std.linux ./gws-std/gws-std.go
	GOOS=linux GOARCH=amd64 go build -o quickws-std.linux ./quickws/quickws-std.go
	GOOS=linux GOARCH=amd64 go build -o test-client.linux ./client/client.go

clean:
	rm *.linux
	rm *.mac

