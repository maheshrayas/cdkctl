hello:

build:
	go build -o cdkctl ./cmd
	mv cdkctl /usr/local/bin

windows:
	env GOOS=windows GOARCH=386  go build -o cdkctl ./cmd

run:
	go run ./cmd/main.go