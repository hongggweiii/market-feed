# Market Feed Ingestor

A high-throughput data pipeline built in Go. This service streams real-time cryptocurrency trade data from the Binance WebSocket API, buffers the events through Apache Kafka, and utilises an in-memory batching engine to write financial data into a ClickHouse OLAP database.

## Architecture Overview

1. **Producer:** Connects to the Binance WebSocket (`BTCUSDT` by default) and streams raw trades.
2. **Message Broker:** Pushes trades to a Kafka topic (`crypto.trades.raw`) to decouple ingestion from storage and handle sudden traffic spikes.
3. **Batching Consumer:** Reads from Kafka and accumulates trades in memory.
4. **Storage:** Flushes batches to ClickHouse every 2 seconds or 1,000 trades

## Setup & Configuration

**1. Clone the repository**
```bash
git clone [https://github.com/yourusername/market-feed.git](https://github.com/yourusername/market-feed.git)
cd market-feed
```
**2. Configure Env Variables**
Create a .env file from .env.example
```bash
cp .env.example .env
```
Ensure your .env has the following values to run locally:
```bash
RUN_MIGRATIONS=true
KAFKA_BROKER=localhost:9092
CLICKHOUSE_ADDR=localhost:9000
```

## Running the Infrastructure
Start local Docker Compose:
```bash
docker compose up -d
```

## Database Migrations
By setting `RUN_MIGRATIONS=true`, app will execute database schema SQL files to build Clickhouse tables. 

If you want to run it manually using CLI:
```bash
migrate -path db/migrations -database "clickhouse://default:@localhost:9000/default" up
```

## Running Application
```bash
go run cmd/ingestor/main.go
```

## Accessing ClickHouse terminal
```bash
docker exec -it clickhouse-db clickhouse-client
```