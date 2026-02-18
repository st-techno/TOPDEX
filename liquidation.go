func LiquidationConsumer() {
    for {
        msg, _ := kafkaConsumer.ReadMessage(100 * time.Millisecond)
        var event RiskEvent
        json.Unmarshal(msg.Value, &event)
        if event.Equity < event.MaintenanceMargin {
            liquidatePosition(event.UserID, event.Symbol)
            // Batch tx to on-chain via relayer
            submitSettlement(event.PositionProof)
        }
    }
}
