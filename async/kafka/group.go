package kafka

type group struct {
	baseConsumer
}

func (g *group) createInfo() {
	g.baseConsumer.createInfo()
	g.info["type"] = "kafka-consumer-group"
}
