.PHONY: prep-conf
prep-conf:
	cp .env.example .env

.PHONY: run-message-emitter
run-message-emitter:
	bash -c 'set -a; . .env; set +a; go run ./cmd/messageemitter/main.go'

.PHONY: run-ban-emitter
run-ban-emitter:
	bash -c 'set -a; . .env; set +a; go run ./cmd/blockemitter/main.go -user=$(user) -block=$(block-user)'

.PHONY: run-processor
run-processor:
	bash -c 'set -a; . .env; set +a; go run ./cmd/processor/main.go'
