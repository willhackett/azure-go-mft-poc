package daemon

import (
	"context"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/willhackett/azure-mft/pkg/logger"
)

type QueueMessage struct {
	context    context.Context
	text       string
	popReceipt azqueue.PopReceipt
	URL        azqueue.MessageIDURL
}

func (qm *QueueMessage) Delete() {
	_, err := qm.URL.Delete(qm.context, qm.popReceipt)
	if err != nil {
		logger.Get().Trace(err)
	}
}

func (qm *QueueMessage) IncreaseLease() {
	log := logger.Get()
	update, err := qm.URL.Update(qm.context, qm.popReceipt, time.Second*120, qm.text)
	if err != nil {
		log.Debug("Failed to increase lease", err)
	} else {
		qm.popReceipt = update.PopReceipt
	}
}
