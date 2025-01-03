## Introduction

This folder includes a Docker Compose file that sets up a Kafka cluster, along with examples of consumer and producer applications configured with various settings.

## Launch
All commands should be launched from the "basics" folder
### Producer 
The producer is configured using env variables. See example configurations in the ./producer/.env.example file.


Launch steps:
1. Lauch kafka broker using:
```
- docker compose -f docker_compose.yml up -d
```
### Producer
2. Launch the producer using:
 ```
 - make prep-producer-conf
 - make run-producer
 ```

 ### Consumer
The consumer takes cli arguments for configuration. To explore available arguments and get help, use the make command:
  ```
 - make consumer-help
 ```
2. Launch the cosumers using:
 ```
 - make run-consumer-1 
 - make run-consumer-2
 ```
