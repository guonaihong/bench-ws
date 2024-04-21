all: build-linux build-mac build-freebsd
	if [ ! -d "./bin" ]; then \
        echo mkdir "./bin"; \
    fi

build-linux:
	# 编译linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/quickws.linux ./wslib/quickws/quickws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/greatws.linux ./wslib/greatws/greatws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/greatws.race.linux ./wslib/greatws/greatws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gws.linux ./wslib/gws/gws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gws-std.linux ./wslib/gws-std/gws-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gorilla.linux ./wslib/gorilla/gorilla.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/gobwas.linux  ./wslib/gobwas/gobwas.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nettyws.linux ./wslib/nettyws/nettyws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-std.linux ./wslib/nbio-std/nbio-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-nonblocking.linux ./wslib/nbio-nonblocking/nbio-nonblocking.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-blocking.linux ./wslib/nbio-blocking/nbio-blocking.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/nbio-mixed.linux ./wslib/nbio-mixed/nbio-mixed.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/hertz-std.linux ./wslib/hertz-std/hertz-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/hertz.linux ./wslib/hertz/hertz.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/fasthttp-ws-std.linux ./wslib/fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/web.linux ./web/web.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/bench-ws.linux ./cmd/bench-ws/bench-ws.go

build-freebsd:
	# 编译freebsd
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/quickws.freebsd ./wslib/quickws/quickws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/greatws.freebsd ./wslib/greatws/greatws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/greatws.race.freebsd ./wslib/greatws/greatws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gws.freebsd ./wslib/gws/gws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gws-std.freebsd ./wslib/gws-std/gws-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gorilla.freebsd ./wslib/gorilla/gorilla.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/gobwas.freebsd  ./wslib/gobwas/gobwas.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nettyws.freebsd ./wslib/nettyws/nettyws.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-std.freebsd ./wslib/nbio-std/nbio-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-nonblocking.freebsd ./wslib/nbio-nonblocking/nbio-nonblocking.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-blocking.freebsd ./wslib/nbio-blocking/nbio-blocking.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/nbio-mixed.freebsd ./wslib/nbio-mixed/nbio-mixed.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/hertz-std.freebsd ./wslib/hertz-std/hertz-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/hertz.freebsd ./wslib/hertz/hertz.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/fasthttp-ws-std.freebsd ./wslib/fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/web.freebsd ./web/web.go
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/bench-ws.freebsd ./cmd/bench-ws/bench-ws.go

build-mac:
	# 编译mac
	GOOS=darwin GOARCH=arm64 go build -o ./bin/quickws.mac ./wslib/quickws/quickws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/greatws.mac ./wslib/greatws/greatws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gws.mac ./wslib/gws/gws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gws-std.mac ./wslib/gws-std/gws-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gorilla.mac ./wslib/gorilla/gorilla.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gobwas.mac ./wslib/gobwas/gobwas.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nettyws.mac ./wslib/nettyws/nettyws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-std.mac ./wslib/nbio-std/nbio-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-nonblocking.mac ./wslib/nbio-nonblocking/nbio-nonblocking.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-blocking.mac ./wslib/nbio-blocking/nbio-blocking.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nbio-mixed.mac ./wslib/nbio-mixed/nbio-mixed.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/hertz-std.mac ./wslib/hertz-std/hertz-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/hertz.mac ./wslib/hertz/hertz.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/fasthttp-ws-std.mac ./wslib/fasthttp-ws-std/fasthttp-ws-std.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/web.mac ./web/web.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/bench-ws.mac ./cmd/bench-ws/bench-ws.go
clean:
	rm *.linux
	rm *.mac

