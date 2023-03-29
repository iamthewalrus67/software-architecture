package common

import (
	"app/internal/logging"
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	UUID uuid.UUID `json:"uuid"`
	Text string    `json:"text"`
}

func NewMessage(uuid uuid.UUID, text string) Message {
	return Message{UUID: uuid, Text: text}

}

func (m *Message) ToJSON() []byte {
	data, err := json.Marshal(m)

	if err != nil {
		logging.ErrorLog.Fatal("Failed to marshal message json")
	}

	return data
}

func (m *Message) String() string {
	return m.Text
}
