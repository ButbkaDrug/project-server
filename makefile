build:
	npm run build
	GOOS=linux GOARCH=386 go build -o bin/ main.go
	GOOS=windows GOARCH=386 go build -o bin/ main.go


run:
	npm run build
	go build -o ./tmp/main ./cmd/.

