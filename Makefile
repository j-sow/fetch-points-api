CC=gcc

build:
	go build -o RunRewardsAPI

test_store:
	CC=$(CC) go test -v api/pkg

