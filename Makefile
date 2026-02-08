s=1000
w=8

server:
	go run ./cmd/server

client:
	go run ./cmd/client -n $(s) -w $(w)

clean:
	rm -rf bin/ report.html
