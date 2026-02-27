package manager

import (
	"context"
	"sync"

	datafuncs "github.com/Egot3/Dialyn/internal/dataFuncs"
	"github.com/Egot3/Zhao/queues"
	"github.com/Egot3/Zhao/sub"
)

type SubscriberManager struct {
	cancelFuncs map[string]context.CancelFunc
	mu          sync.Mutex
	subscriber  *sub.Subscriber
}

func NewSubscriberManager(sub *sub.Subscriber) *SubscriberManager {
	return &SubscriberManager{
		subscriber: sub,
	}
}

func (sm *SubscriberManager) Reconcile(newQsChan <-chan []*queues.QueueStruct) {
	newQs := <-newQsChan

	sm.mu.Lock()
	defer sm.mu.Unlock()

	var toAdd []*queues.QueueStruct

	existingQueues := datafuncs.SliceFromKeys(sm.cancelFuncs)

	//diffEN := datafuncs.Difference(existingQueues, newQs)
}
