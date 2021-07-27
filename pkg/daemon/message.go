package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

type Message struct {
	context context.Context
	text string
	popReceipt *azqueue.PopReceipt
	URL azqueue.MessageIDURL
}

func (m *Message) Delete() {
		_, err := m.URL.Delete(m.context, *m.popReceipt)
	if err != nil {
		fmt.Println("Something went wrong", err)
	}
}

func (m *Message) IncreaseLease() {
	update, err := m.URL.Update(m.context, *m.popReceipt, time.Second * 60, m.text)
	if err != nil {
		fmt.Println("Something went wrong", err)
	} else {
		*m.popReceipt = update.PopReceipt
	}
}