package ctapi

import (
	"golang.org/x/sys/windows"
	"time"
)

type Subscription struct {
	cancel chan bool
}

type SubscriptionUpdate struct {
}

func (this *CtApi) Subscribe(hConnection windows.Handle, updateRate time.Duration, tags []string) (chan SubscriptionUpdate, error) {
	ch := make(chan SubscriptionUpdate)

	return ch, nil
}
