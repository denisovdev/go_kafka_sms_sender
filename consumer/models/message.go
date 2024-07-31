package models

import "encoding/json"

type Message struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type QueueMessage struct {
	Offset  int
	Payload []byte
}

func (qm *QueueMessage) ConvertToMessage() (*Message, error) {
	message := new(Message)
	err := json.Unmarshal(qm.Payload, &message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
