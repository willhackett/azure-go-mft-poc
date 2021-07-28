package registry

import (
	"fmt"
	"time"

	"github.com/willhackett/azure-mft/pkg/constant"
)

const (
	IntervalDuration = time.Second * 60
)

type Transfer struct {
	ID         string
	Details    constant.FileRequestMessage
	Expiration int64
}

var (
	transfers = make(map[string]Transfer)
)

func AddTransfer(id string, obj constant.FileRequestMessage, expiresIn int64) {
	transfers[id] = Transfer{
		ID:         id,
		Details:    obj,
		Expiration: time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
	}
}

func DeleteTransfer(id string) {
	delete(transfers, id)
}

func GetTransfer(id string) (Transfer, bool) {
	t, ok := transfers[id]
	fmt.Println(transfers)
	return t, ok
}

func DeleteExpired() {
	for id, t := range transfers {
		if t.Expiration > 0 && t.Expiration < time.Now().Unix() {
			delete(transfers, id)
		}
	}
}

func init() {
	ticker := time.NewTicker(IntervalDuration)

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				DeleteExpired()
			}
		}
	}()
}
