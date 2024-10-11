# Microservices

This is a monorepo containing code for microservices pertaining to an ecommerce webapp

- [Installation](#installation)
- [Description](#description)
  - [GraphQL API](#1-graphql-api)
  - [Users Service](#2-users-service)
  - [Products Service](#3-products-service)
  - [Orders Service](#4-orders-service)
- [Architecture](#architecture)
- [Todos](#todos)

## Installation

```bash
# 1. Clone and cd into the repo
git clone https://github.com/ary82/microservices.git
cd microservices

# 2. Start the dependencies
docker compose up -d

# 3. Copy the .env
cp .env.example .env

# 4. Start each service
make run-api
make run-user
make run-product
make run-order
```

## Description

This repo contains the following microservices with careful design practices. Their entrypoints are in `cmd/` and their logic is in their independent directories in `internal/`. They use two common packages:

- `mq`(location: `internal/mq`), contains types needed for data exchange in Message Queues
- `proto`(location: `internal/proto`), contains `.proto` files and generated code for grpc

### 1. GraphQL API

This is the GraphQL API that will be exposed to the client, fetching and aggregating data from all three microservices.

### 2. Users Service

Handles User registration and Authentication.

### 3. Products Service

Handles products catalogue and Inventory.

### 4. Orders Service

Handles order placing and fetching.

## Architecture

Here's how everything interacts with each other:

<p align="center">
    <img  src="./docs/graphviz/arch.png">
</p>

Here's how the microservices interact with each other and with RabbitMQ:

<p align="center">
    <img  src="./docs/graphviz/rabbitmq.png">
</p>

## Todos

- [ ] Dockerize microservices
- [ ] Cache & dupe some data at various locations
- [ ] Observability
- [ ] Tests
- [ ] Docs
