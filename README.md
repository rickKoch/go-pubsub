# go-pubsub

Simple queue manager inspired by google pubsub.

To start the queue just run `go run .`. This will start the grpc server and
the queue manager.

List servives:
```
grpcurl -plaintext localhost:8080 list
```

describe service
```
grpcurl -plaintext localhost:8080 describe go_pubsub.PubSubService
```


describe method
```
grpcurl -plaintext localhost:8080 describe go_pubsub.PubSubService.CreateTopic
```

## Methods

Create topic:
```
grpcurl -d '{"name": "test-topic"}' -plaintext localhost:8080 go_pubsub.PubSubService/CreateTopic
```

Create subscription:
```
grpcurl -d '{"topicName": "test-topic", "subscriptionName": "test-sub"}' -plaintext localhost:8080 go_pubsub.PubSubService/CreateSubscription
```

Publish message:
```
grpcurl -d '{"topic": "test-topic", "message": "test message"}' -plaintext localhost:8080 go_pubsub.PubSubService/PublishMessage

```

Pull messages:
```
grpcurl -d '{"subscription": "test-sub"}' -plaintext localhost:8080 go_pubsub.PubSubService/PullMessages
```

Pull messages in stream:
```
grpcurl -d '{"subscription": "test-sub"}' -plaintext localhost:8080 go_pubsub.PubSubService/PullMessagesStreaming
```

Acknowledge message:
```
grpcurl -d '{"subscription": "test-sub", "ackId": "1"}' -plaintext localhost:8080 go_pubsub.PubSubService/AcknowledgeMessage
```
