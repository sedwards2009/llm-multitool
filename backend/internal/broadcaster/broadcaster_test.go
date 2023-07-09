package broadcaster

import "testing"

func TestBroadcasterHappy(t *testing.T) {
	broadcaster := NewBroadcaster()
	id := "1234567890"

	listenerChan := make(chan string, 16)
	broadcaster.Register(id, listenerChan)
	broadcaster.Send(id, "changed")
	sentMessage := <-listenerChan
	if sentMessage != "changed" {
		t.Errorf("Didn't receive 'changed' message.")
	}

	broadcaster.Unregister(listenerChan)
	broadcaster.Send(id, "changed")

	broadcaster.Quit()
}
func TestBroadcasterDifferentId(t *testing.T) {
	broadcaster := NewBroadcaster()
	id := "1234567890"
	otherId := "abcdefg"

	listenerChan := make(chan string, 16)
	broadcaster.Register(id, listenerChan)
	broadcaster.Send(otherId, "changed")

	if len(listenerChan) != 0 {
		t.Errorf("listenerChan should be empty.")
	}

	broadcaster.Unregister(listenerChan)
	broadcaster.Send(id, "changed")
	if len(listenerChan) != 0 {
		t.Errorf("listenerChan should be empty after Unregister.")
	}

	broadcaster.Quit()
}
