size=1000
workers=8

server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go -n $(size) -w $(workers)

clean:
	rm -rf bin/
