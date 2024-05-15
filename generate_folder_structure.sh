#!/bin/bash

# Root folder
mkdir -p multi-channel-framework

# Root-level files
touch multi-channel-framework/docker-compose.yml
touch multi-channel-framework/db_init.sql
touch multi-channel-framework/main.go

# Models folder
mkdir -p multi-channel-framework/models
touch multi-channel-framework/models/models.go

# Repositories folder
mkdir -p multi-channel-framework/repositories
touch multi-channel-framework/repositories/repositories.go

# Services folder
mkdir -p multi-channel-framework/services
touch multi-channel-framework/services/services.go

# Handlers folder
mkdir -p multi-channel-framework/handlers
touch multi-channel-framework/handlers/http_handlers.go

# MQTT folder
mkdir -p multi-channel-framework/mqtt
touch multi-channel-framework/mqtt/factory.go
touch multi-channel-framework/mqtt/mqtt_handler.go
touch multi-channel-framework/mqtt/mqtt_client.go

# RabbitMQ folder
mkdir -p multi-channel-framework/rabbitmq
touch multi-channel-framework/rabbitmq/consumer.go
touch multi-channel-framework/rabbitmq/producer.go

# Middleware folder
mkdir -p multi-channel-framework/middleware
touch multi-channel-framework/middleware/auth_middleware.go

# Mosquitto folder
mkdir -p multi-channel-framework/mosquitto/certs
touch multi-channel-framework/mosquitto/mosquitto.conf
touch multi-channel-framework/mosquitto/certs/ca.crt
touch multi-channel-framework/mosquitto/certs/server.crt
touch multi-channel-framework/mosquitto/certs/server.key

echo "Folder structure and placeholder files have been created."
