# Defaults
size=1000
workers=8

server:
	go run ./cmd/server

client:
	go run ./cmd/client -n $(size) -w $(workers)

clean:
	rm -rf bin/ report.html
