package runner

import "testing"

func TestTimeoutServer(t *testing.T) {
	timeoutServerJob := NewTimeoutServerJob()
	timeoutServerJob.Run()
}
