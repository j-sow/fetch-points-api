build:
	go mod download github.com/google/btree
	go build -o RunRewardsAPI cmd/main.go 

run:
	go run cmd/main.go

run-docker:
	docker run --rm -it -v $(PWD):/api -p 8080:8080 golang make -C /api run

test:
	CGO_ENABLED=0 go test -v api/pkg

test-docker:
	docker run --rm -it -v $(PWD):/api golang make -C /api test
	
