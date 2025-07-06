all: build-linux build-mac 
	if [ ! -d "./bin" ]; then \
        echo mkdir "./bin"; \
    fi

build-web:
	GOOS=linux GOARCH=amd64 go build -o ./bin/web.linux ./web/web.go

build-linux:
	# 编译linux
	#GOOS=linux GOARCH=amd64 go build -o ./bin/nhooyr.linux ./lib/nhooyr/nhooyr.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/quickws.linux ./lib/quickws/quickws.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/greatws-io.linux ./lib/greatws-io/greatws-io.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/great-elastic.linux ./lib/great-elastic/great-elastic.go
	GOOS=linux GOARCH=amd64 go build -o ./bin/greatws-onebyone.linux ./lib/greatws-onebyone/greatws-onebyone.go
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


build-mac:
	# 编译mac
	GOOS=darwin GOARCH=arm64 go build -o ./bin/nhooyr.mac ./lib/nhooyr/nhooyr.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/quickws.mac ./lib/quickws/quickws.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/greatws-io.mac ./lib/greatws-io/greatws-io.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/greatws-onebyone.mac ./lib/greatws-onebyone/greatws-onebyone.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/great-elastic.mac ./lib/great-elastic/great-elastic.go
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

