package matcher

import (
    "container/heap"
    "sync"
    "github.com/redis/go-redis/v9"
    "github.com/twmb/murmur3"
)

type Order struct {
    ID        string
    Side      Side
    Price     float64
    Quantity  float64
    UserID    string
    Timestamp int64
}

type OrderBook struct {
    mu      sync.RWMutex
    bids    *PriceLevelHeap
    asks    *PriceLevelHeap
    redis   *redis.Client // For persistence/pubsub
}

func (ob *OrderBook) AddOrder(o *Order) ([]Trade, error) {
    ob.mu.Lock()
    defer ob.mu.Unlock()
    
    // Shard by symbol hash for concurrency
    shardID := murmur3.Sum32([]byte(o.Symbol)) % 64
    
    var trades []Trade
    switch o.Side {
    case Buy:
        trades = ob.match(o, ob.asks)
    case Sell:
        trades = ob.match(o, ob.bids)
    }
    if o.RemainingQuantity > 0 {
        heap.Push(ob.bidsOrAsks(o.Side), o)
        ob.redis.ZAdd(ctx, ob.key(o.Side), &redis.Z{Score: o.Price, Member: o.ID})
    }
    // Publish Kafka event
    publishEvent("order_update", o)
    return trades, nil
}

// Risk check integration before match
func (ob *OrderBook) PreMatchRisk(userID string, posSize float64) bool {
    pos, _ := getPosition(userID) // From Redis/Postgres
    return calculateIMR(pos, posSize) <= getAvailableMargin(userID)
}

## This uses priority heaps for price-time priority and Redis for fast L2 snapshots.


