all: build-linux build-mac build-freebsd
	if [ ! -d "./bin" ]; then \
        echo mkdir "./bin"; \
    fi

build-web:
	GOOS=linux GOARCH=amd64 go build -o ./bin/web.linux ./web/web.go

build-linux:
	# 编译linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/quickws.linux ./lib/quickws/quickws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/greatws.linux ./lib/greatws/greatws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/greatws.race.linux ./lib/greatws/greatws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gws.linux ./lib/gws/gws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gws-std.linux ./lib/gws-std/gws-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gorilla.linux ./lib/gorilla/gorilla.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gobwas.linux  ./lib/gobwas/gobwas.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nettyws.linux ./lib/nettyws/nettyws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-std.linux ./lib/nbio-std/nbio-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-nonblocking.linux ./lib/nbio-nonblocking/nbio-nonblocking.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-blocking.linux ./lib/nbio-blocking/nbio-blocking.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-mixed.linux ./lib/nbio-mixed/nbio-mixed.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/hertz-std.linux ./lib/hertz-std/hertz-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/hertz.linux ./lib/hertz/hertz.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/fasthttp-ws-std.linux ./lib/fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/web.linux ./web/web.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/bench-ws.linux ./cmd/bench-ws/bench-ws.go

build-freebsd:
	# 编译freebsd
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/quickws.freebsd ./lib/quickws/quickws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/greatws.freebsd ./lib/greatws/greatws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/greatws.race.freebsd ./lib/greatws/greatws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gws.freebsd ./lib/gws/gws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gws-std.freebsd ./lib/gws-std/gws-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gorilla.freebsd ./lib/gorilla/gorilla.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gobwas.freebsd  ./lib/gobwas/gobwas.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nettyws.freebsd ./lib/nettyws/nettyws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-std.freebsd ./lib/nbio-std/nbio-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-nonblocking.freebsd ./lib/nbio-nonblocking/nbio-nonblocking.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-blocking.freebsd ./lib/nbio-blocking/nbio-blocking.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-mixed.freebsd ./lib/nbio-mixed/nbio-mixed.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/hertz-std.freebsd ./lib/hertz-std/hertz-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/hertz.freebsd ./lib/hertz/hertz.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/fasthttp-ws-std.freebsd ./lib/fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/web.freebsd ./web/web.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/bench-ws.freebsd ./cmd/bench-ws/bench-ws.go

build-mac:
	# 编译mac
	GOOS=darwin GOARCH=arm64 go build -o ./bin/quickws.mac ./lib/quickws/quickws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/greatws.mac ./lib/greatws/greatws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gws.mac ./lib/gws/gws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gws-std.mac ./lib/gws-std/gws-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gorilla.mac ./lib/gorilla/gorilla.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gobwas.mac ./lib/gobwas/gobwas.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nettyws.mac ./lib/nettyws/nettyws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-std.mac ./lib/nbio-std/nbio-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-nonblocking.mac ./lib/nbio-nonblocking/nbio-nonblocking.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-blocking.mac ./lib/nbio-blocking/nbio-blocking.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-mixed.mac ./lib/nbio-mixed/nbio-mixed.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/hertz-std.mac ./lib/hertz-std/hertz-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/hertz.mac ./lib/hertz/hertz.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/fasthttp-ws-std.mac ./lib/fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/web.mac ./web/web.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/bench-ws.mac ./cmd/bench-ws/bench-ws.go
clean:
	rm *.linux
	rm *.mac

