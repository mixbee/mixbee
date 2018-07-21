

package vbft

import "testing"

func constructEventTimer() *EventTimer {
	server := constructServer()
	return NewEventTimer(server)
}

func TestStartTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	err := eventtimer.StartTimer(1, 10)
	t.Logf("TestStartTimer: %v", err)
}

func TestCancelTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	err := eventtimer.StartTimer(1, 10)
	t.Logf("TestStartTimer: %v", err)
	err = eventtimer.CancelTimer(1)
	t.Logf("TestCancelTimer: %v", err)
}

func TestStartEventTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	err := eventtimer.startEventTimer(EventProposeBlockTimeout, 1)
	t.Logf("TestStartEventTimer: %v", err)
}

func TestCancelEventTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	err := eventtimer.startEventTimer(EventProposeBlockTimeout, 1)
	t.Logf("startEventTimer: %v", err)
	err = eventtimer.cancelEventTimer(EventProposeBlockTimeout, 1)
	t.Logf("cancelEventTimer: %v", err)
}
