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

func GenerateNewMessage(text string) Message {
	return Message{UUID: uuid.New(), Text: text}
}

func (m *Message) ToJSON() []byte {
	data, err := json.Marshal(m)

	if err != nil {
		logging.ErrorLog.Fatal("Failed to marshal message json")
	}

	return data
}

func MessageFromString(s string) (Message, error) {
	return MessageFromBytes([]byte(s))
}

func MessageFromBytes(b []byte) (Message, error) {
	message := Message{}
	err := json.Unmarshal([]byte(b), &message)

	if err != nil {
		return Message{}, err
	}

	return message, nil
}

func (m *Message) String() string {
	return m.Text
}
