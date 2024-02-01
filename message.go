package ServiceBus

type Message interface {
	GetRoutingKey() string
}

type TestMessage struct {
	Key  string
	Body string
}

func (tm *TestMessage) GetRoutingKey() string {
	return "test_routing_key"
}
