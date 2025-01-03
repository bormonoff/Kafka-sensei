.PHONY: prep-producer-conf
prep-producer-conf:
	cp .env.example .env

.PHONY: run-producer
run-producer:
	bash -c 'set -a; . .env; set +a; go run ./cmd/producer/main.go'

.PHONY: consumer-help
consumer-help:
	go run ./cmd/consumer/main.go --help

.PHONY: run-consumer-1
run-consumer-1:
	go run ./cmd/consumer/main.go \
		-broker=localhost:9094,localhost:9095,localhost:9096 -topic=shapes \
		-group=group1 -poll-period=1000 -auto-commit=false

.PHONY: run-consumer-2
run-consumer-2:
	go run ./cmd/consumer/main.go \
		-broker=localhost:9094,localhost:9095,localhost:9096 -topic=shapes \
		-group=group2 -poll-period=1000 -auto-commit=true
		