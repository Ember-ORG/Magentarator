kill () {
        pkill -e gopherjs
        rm ./public/client.js
        rm ./public/client.js.map
}
trap kill EXIT
GOOS=linux gopherjs build client.go -w -o ./public/client.js &
go run server.go
