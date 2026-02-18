# TOPDEX

TOPDEX is a high-performance, hybrid DEX backend emphasizing low-latency event-driven trading, margining, liquidations, funding rates, oracle pricing, and risk management. It uses Go for the core matching engine and services, TypeScript for the frontend/API gateway, PostgreSQL for persistent storage, Kafka for event streaming, and Redis for caching/order books to achieve sub-millisecond latencies in off-chain execution with on-chain settlement hooks.

Architecture Overview

TOPDEX employs a microservices architecture with event-driven design for scalability. Core components include a Go-based matching engine (low-latency order matching), Kafka producers/consumers 
for trade events/liquidations, Redis for in-memory order books and risk caches, and PostgreSQL for user positions/margin history. Off-chain logic handles real-time pricing/oracle feeds 
(e.g., via Chainlink), funding calculations (hourly rates), and risk checks; on-chain integration uses Ethereum/Solana SDKs for custody/settlement via smart contract calls 
(e.g., position openings via relayers for gasless user experience).

Services communicate via gRPC internally and REST/WS externally. Docker containers orchestrate via Kubernetes; CI/CD uses GitHub Actions for testing/linting/builds/deployments.

Key Features Implemented

Perpetual Futures Trading: Supports isolated/cross margin, maker/taker fees, leverage up to 50x, non-expiring contracts on crypto/stocks/FX.

Low-Latency Matching: Go engine with lock-free data structures; Redis sorted sets for L2/L3 order books; <1ms match times via sharding.

Risk & Liquidation: Real-time PnL monitoring, auto-liquidation via Kafka events if maintenance margin breached (e.g., equity < 5% collateral).

Funding & Pricing: Hourly funding payments (longs pay shorts if premium >0); oracle-based mark prices to prevent manipulation.

On/Off-Chain Hybrid: Off-chain execution for speed; batch settlements to L1 via Merkle proofs for self-custody/verifiability.

# Deployment & Ops

Dockerfile (Go service):

FROM golang:1.22-alpine AS builder

COPY . /app

RUN go build -o matcher ./cmd/matcher

FROM alpine:latest

COPY --from=builder /app/matcher /matcher

CMD ["/matcher"]

docker-compose.yml (dev stack):

services:

  postgres:
  
    image: postgres:16
    
    environment:
    
      POSTGRES_DB: topdex
  
  redis:
  
    image: redis:7-alpine
  
  kafka:
  
    image: confluentinc/cp-kafka:7.5.0
  
  matcher:
  
    build: .
   
    depends_on: [redis, kafka]

CI/CD (GitHub Actions): Lints (golangci-lint), tests (95% coverage), builds Docker images, deploys to EKS. Helm charts for prod scaling (HPA on CPU<70%).

Scaling & Patterns

Cache: Redis Cluster for order books (TTL 1s snapshots); LRU for positions.

Event-Driven: Kafka topics: trades, liquidations, funding, risk_alerts.

On-Chain Bridge: Go relayer batches 1000 txs/epoch to Solana/Eth for settlements, using custom state channels via Merkle trees.

Monitoring: Prometheus/Grafana for latency/PnL; circuit breakers for oracles.

This forms a production-ready skeleton.

