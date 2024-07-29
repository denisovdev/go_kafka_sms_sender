package models

type Message struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type MessageRequest struct {
	Phone string `json:"phone" binding:"required,min=11"`
}

type MessageStorage struct {
	ID      int               `json:"id"`
	Payload map[string]string `json:"payload"`
}

func (msgReq *MessageRequest) ConvertToMessage(code string) *Message {
	return &Message{
		Phone: msgReq.Phone,
		Code:  code,
	}
}
