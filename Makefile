all: url-shortener

url-shortener: clean
	go build -o $@ cmd/url-shortener/main.go

clean:
	rm -rf url-shortener