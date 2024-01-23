all:
	# 编译mac
	GOOS=darwin	GOARCH=arm64 go build -o quickws.mac ./quickws/quickws.go
	GOOS=darwin GOARCH=arm64 go build -o greatws.mac ./greatws/greatws.go
	GOOS=darwin	GOARCH=arm64 go build -o gws.mac ./gws/gws.go
	GOOS=darwin	GOARCH=arm64 go build -o gws-std.mac ./gws-std/gws-std.go
	GOOS=darwin	GOARCH=arm64 go build -o gorilla.mac ./gorilla/gorilla.go
	GOOS=darwin	GOARCH=arm64 go build -o gobwas.mac ./gobwas/gobwas.go
	GOOS=darwin	GOARCH=arm64 go build -o nettyws.mac ./nettyws/nettyws.go
	GOOS=darwin	GOARCH=arm64 go build -o nbio-std.mac ./nbio-std/nbio-std.go
	GOOS=darwin	GOARCH=arm64 go build -o nbio-nonblocking.mac ./nbio-nonblocking/nbio-nonblocking.go
	GOOS=darwin	GOARCH=arm64 go build -o nbio-blocking.mac ./nbio-blocking/nbio-blocking.go
	GOOS=darwin	GOARCH=arm64 go build -o nbio-mixed.mac ./nbio-mixed/nbio-mixed.go
	GOOS=darwin	GOARCH=arm64 go build -o hertz-std.mac ./hertz-std/hertz-std.go
	GOOS=darwin	GOARCH=arm64 go build -o hertz.mac ./hertz/hertz.go
	GOOS=darwin GOARCH=arm64 go build -o fasthttp-ws-std.mac ./fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=darwin	GOARCH=arm64 go build -o test-client.mac ./client/client.go

	# 编译linux
	GOOS=linux GOARCH=amd64 go build -o quickws.linux ./quickws/quickws.go
	GOOS=linux GOARCH=amd64 go build -o greatws.linux ./greatws/greatws.go
	GOOS=linux GOARCH=amd64 go build -o greatws.race.linux ./greatws/greatws.go
	GOOS=linux GOARCH=amd64 go build -o gws.linux ./gws/gws.go
	GOOS=linux GOARCH=amd64 go build -o gws-std.linux ./gws-std/gws-std.go
	GOOS=linux GOARCH=amd64 go build -o gorilla.linux ./gorilla/gorilla.go
	GOOS=linux GOARCH=amd64 go build -o gobwas.linux  ./gobwas/gobwas.go
	GOOS=linux GOARCH=amd64 go build -o nettyws.linux ./nettyws/nettyws.go
	GOOS=linux GOARCH=amd64 go build -o nbio-std.linux ./nbio-std/nbio-std.go
	GOOS=linux GOARCH=amd64 go build -o nbio-nonblocking.linux ./nbio-nonblocking/nbio-nonblocking.go
	GOOS=linux GOARCH=amd64 go build -o nbio-blocking.linux ./nbio-blocking/nbio-blocking.go
	GOOS=linux GOARCH=amd64 go build -o nbio-mixed.linux ./nbio-mixed/nbio-mixed.go
	GOOS=linux GOARCH=amd64 go build -o hertz-std.linux ./hertz-std/hertz-std.go
	GOOS=linux GOARCH=amd64 go build -o hertz.linux ./hertz/hertz.go
	GOOS=linux GOARCH=amd64 go build -o fasthttp-ws-std.linux ./fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=linux GOARCH=amd64 go build -o test-client.linux ./client/client.go

	# 编译freebsd
	GOOS=freebsd GOARCH=amd64 go build -o quickws.freebsd ./quickws/quickws.go
	GOOS=freebsd GOARCH=amd64 go build -o greatws.freebsd ./greatws/greatws.go
	GOOS=freebsd GOARCH=amd64 go build -o greatws.race.freebsd ./greatws/greatws.go
	GOOS=freebsd GOARCH=amd64 go build -o gws.freebsd ./gws/gws.go
	GOOS=freebsd GOARCH=amd64 go build -o gws-std.freebsd ./gws-std/gws-std.go
	GOOS=freebsd GOARCH=amd64 go build -o gorilla.freebsd ./gorilla/gorilla.go
	GOOS=freebsd GOARCH=amd64 go build -o gobwas.freebsd  ./gobwas/gobwas.go
	GOOS=freebsd GOARCH=amd64 go build -o nettyws.freebsd ./nettyws/nettyws.go
	GOOS=freebsd GOARCH=amd64 go build -o nbio-std.freebsd ./nbio-std/nbio-std.go
	GOOS=freebsd GOARCH=amd64 go build -o nbio-nonblocking.freebsd ./nbio-nonblocking/nbio-nonblocking.go
	GOOS=freebsd GOARCH=amd64 go build -o nbio-blocking.freebsd ./nbio-blocking/nbio-blocking.go
	GOOS=freebsd GOARCH=amd64 go build -o nbio-mixed.freebsd ./nbio-mixed/nbio-mixed.go
	GOOS=freebsd GOARCH=amd64 go build -o hertz-std.freebsd ./hertz-std/hertz-std.go
	GOOS=freebsd GOARCH=amd64 go build -o hertz.freebsd ./hertz/hertz.go
	GOOS=freebsd GOARCH=amd64 go build -o fasthttp-ws-std.freebsd ./fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=freebsd GOARCH=amd64 go build -o test-client.freebsd ./client/client.go
clean:
	rm *.linux
	rm *.mac

