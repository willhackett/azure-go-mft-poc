package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

type QueueMessage struct {
	context    context.Context
	text       string
	popReceipt *azqueue.PopReceipt
	URL        azqueue.MessageIDURL
}

func (qm *QueueMessage) Delete() {
	_, err := qm.URL.Delete(qm.context, *qm.popReceipt)
	if err != nil {
		fmt.Println("Something went wrong", err)
	}
}

func (qm *QueueMessage) IncreaseLease() {
	update, err := qm.URL.Update(qm.context, *qm.popReceipt, time.Second*120, qm.text)
	if err != nil {
		fmt.Println("Something went wrong", err)
	} else {
		*qm.popReceipt = update.PopReceipt
	}
}
