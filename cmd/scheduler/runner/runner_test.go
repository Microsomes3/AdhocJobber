package runner

import (
	"sync"
	"testing"
)

func TestHeartbeat(t *testing.T) {
	heartbeat := NewHeartbeat()

	var wg sync.WaitGroup
	wg.Add(1)

	heartbeat.StartTaskPollers(&wg)

	wg.Wait()
}
