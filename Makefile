.PHONY: proto build run-server run-client clean

proto:
	docker-compose run --rm protoc

build:
	docker-compose build

run-server:
	docker-compose up server

run-client:
	docker-compose up client

clean:
	docker-compose down
	rm -f server/server client/client protoc