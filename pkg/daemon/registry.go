package daemon

import (
	"time"
)

const (
	IntervalDuration = time.Second * 60
)

type Transfer struct {
	ID					string
	Object	interface{}
	Expiration int64
}

var (
	transfers = make(map[string]Transfer)
)

func Add(id string, obj interface{}, expiresIn int64) {
	transfers[id] = Transfer{
		ID:				id,
		Object:			obj,
		Expiration:	time.Now().Unix() + expiresIn,
	}
}

func Delete(id string) {
	delete(transfers, id)
}

func Get(id string) (Transfer, bool) {
	t, ok := transfers[id]
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