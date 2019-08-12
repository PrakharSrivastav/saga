## Redis channels

From Orchestrator to services
- OrderChannel
- PaymentChannel
- RestaurantChannel
- DeliveryChannel

From all services to Orchestrator
- ReplyChannel

## Message Types

From Orchestrator to services (Order,Payment,Restaurant,Delivery)
- Start
- Rollback

From services to orchestrator (ReplyChannel)
- Done
- Error
