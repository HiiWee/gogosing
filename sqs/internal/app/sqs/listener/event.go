package listener

type ConsumedEvent struct {
	ReceiptHandle *string
	body          ConsumedMessage
}

type ConsumedMessage struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func (c *ConsumedMessage) GetMessage() string {
	return c.Message
}

func (c *ConsumedMessage) GetFrom() string {
	return c.From
}
