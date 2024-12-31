.PHONY: prep-producer-conf
prep-producer-conf:
	cp ./producer/.env.example ./producer/.env

.PHONY: run-producer
run-producer:
	bash -c 'cd producer; set -a; . .env; set +a; go run ./cmd/producer/main.go'

.PHONY: consumer-help
consumer-help:
	cd consumer && \
	go run ./cmd/consumer/main.go --help

.PHONY: run-consumer-1
run-consumer-1:
	cd consumer && \
	go run ./cmd/consumer/main.go \
		-broker=localhost:9094,localhost:9095,localhost:9096 -topic=orders \
		-group=group1 -poll-period=1000 -auto-commit=false

.PHONY: run-consumer-2
run-consumer-2:
	cd consumer && \
	go run ./cmd/consumer/main.go \
		-broker=localhost:9094,localhost:9095,localhost:9096 -topic=orders \
		-group=group2 -poll-period=1000 -auto-commit=true
		