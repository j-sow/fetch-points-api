CC=clang #required on ubuntu for http modules?

build:
	go build -o RunRewardsAPI cmd/main.go 

run:
	go run cmd/main.go

test:
	CC=$(CC) go test -v api/pkg

