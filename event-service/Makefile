.SILENT:

build:
	go build -gcflags="all=-N -l" -o dist/event-service main.go

docker-build:
	docker build -t event-service:local .

docker-run:
	docker run --rm -it \
	--network=host \
	-p 8081:8081 \
	-p 8082:8082 \
	event-service:local

run:
	go run main.go

generate:
	oa3 go -w --param timetype=chrono -o oa3gen swagger.yml
