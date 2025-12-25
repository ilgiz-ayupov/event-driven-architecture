package publisher

// type natsPublisher struct {
// 	conn *nats.Conn
// }

// func NewNATSPublisher(
// 	conn *nats.Conn,
// ) usecase.EventPublisher {
// 	return &natsPublisher{
// 		conn: conn,
// 	}
// }

// func (p *natsPublisher) Publish(event usecase.Event) error {
// 	data, err := json.Marshal(event)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal event: %w", err)
// 	}

// 	return p.conn.Publish(event.EventType(), data)
// }
