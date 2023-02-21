package responses

type Response struct {
	MessageType string      `json:"message_type" `
	Message     string      `json:"message"`
	Data        interface{} `json:"data" extensions:"x-nullable"`
}

func NewResponse(messageType string, message string, data interface{}) Response {
	return Response{messageType, message, data}
}
