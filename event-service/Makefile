.SILENT:

build:
	go build -gcflags="all=-N -l" -o dist/event-service main.go

run:
	go run main.go

generate:
	oa3 go -w --param timetype=chrono -o oa3gen swagger.yml
