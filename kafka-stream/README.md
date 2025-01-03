# Introductioin

~~Here we go again.~~ I've recently gotten deep into Kafka Streams. 

As a project I've wrote a real-time chat, with broadcast messages, censoring, and the ability to ban users – everything as we like.

## Top-level architecture
For a better understanding of the data sources, I've wrote a top-level architecture diagram:

![](./docs/kafka_v2.svg)

### Explanation

**System Components**:
- Message Emitter: An external service (e.g., a user) responsible for publishing messages to the mp-messages-stream Kafka topic. This service is located at /cmd/messageemitter.
- Ban Emitter: An external service (e.g., a client or admin interface) that emits events containing user ban information to the mp-blocked-users-stream Kafka topic. This service is located at /cmd/blockemitter.
- Processor Microservice: This microservice consumes data from both the mp-messages-stream and mp-blocked-users-stream Kafka topics. It censors incoming messages and aggregates user ban events into a Kafka table for subsequent join operations.

**MP** is the **message processor** abbreviation. 



## Launch
1. Lauch kafka broker and create topics using:
```
- docker compose -f docker_compose.yml 
```
2. Prepare the .env file. One file is used by all services:
```
- make prep-conf
```
3. Run services:
```
Message emitter:
- make run-message-emitter
Processor:
- make run-processor
```
Ban emitter has differ logic. It emits one msg during it's work. In the example below user-1 blocks user-3 and will not gain messages from the user-3
```
Message emitter:
- make run-ban-emitter user=user-1 block-user=user-3 
```