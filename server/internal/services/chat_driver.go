package services

import (
	"errors"
	"github.com/CryptoSingh1337/serpent-clash/server/internal/utils"
	"github.com/labstack/echo/v4"
)

type ChatDriver struct {
	ChatBuffer chan *utils.ChatMessage
}

func NewChatDriver() *ChatDriver {
	return &ChatDriver{
		ChatBuffer: make(chan *utils.ChatMessage, 100),
	}
}

func (cd *ChatDriver) publishMessage(message *utils.ChatMessage) error {
	if len(cd.ChatBuffer) == 100 {
		return errors.New("chat buffer is full")
	}
	cd.ChatBuffer <- message
	return nil
}

func (cd *ChatDriver) broadcastMessage(message *utils.ChatMessage, rmap map[string]*echo.Response) error {
	return nil
}

func (cd *ChatDriver) close() {
	close(cd.ChatBuffer)
}
