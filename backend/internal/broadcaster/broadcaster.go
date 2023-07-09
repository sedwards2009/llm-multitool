package broadcaster

import (
	"github.com/bobg/go-generics/slices"
)

type listener struct {
	id           string
	listenerChan chan string
}

type messageType uint8

const (
	messageType_Register messageType = iota
	messageType_Unregister
	messageType_Send
	messageType_Quit
)

type message struct {
	messageType  messageType
	listenerChan chan string
	id           string
	message      string
}

type Broadcaster struct {
	listeners    []listener
	toWorkerChan chan message
	doneChan     chan bool
}

func NewBroadcaster() *Broadcaster {
	broadcaster := &Broadcaster{
		listeners:    []listener{},
		toWorkerChan: make(chan message, 16),
		doneChan:     make(chan bool, 1),
	}
	go broadcaster.worker(broadcaster.toWorkerChan, broadcaster.doneChan)
	return broadcaster
}

func (this *Broadcaster) worker(in chan message, done chan bool) {
	for message := range in {
		switch message.messageType {
		case messageType_Register:
			this.listeners = append(this.listeners, listener{message.id, message.listenerChan})

		case messageType_Unregister:
			targetChan := message.listenerChan
			this.listeners, _ = slices.Filter(this.listeners,
				func(l listener) (bool, error) {
					return l.listenerChan != targetChan, nil
				})

		case messageType_Send:
			id := message.id
			for _, listener := range this.listeners {
				if listener.id == id {
					listener.listenerChan <- message.message
				}
			}
		case messageType_Quit:
			this.listeners = []listener{}
			close(in)
			done <- true
			return
		}
	}
}

func (this *Broadcaster) Register(id string, listenerChan chan string) {
	this.toWorkerChan <- message{messageType: messageType_Register, id: id, listenerChan: listenerChan}
}

func (this *Broadcaster) Unregister(listenerChan chan string) {
	this.toWorkerChan <- message{messageType: messageType_Unregister, listenerChan: listenerChan}

}

func (this *Broadcaster) Send(id string, messageText string) {
	this.toWorkerChan <- message{messageType: messageType_Send, id: id, message: messageText}
}

func (this *Broadcaster) Quit() {
	this.toWorkerChan <- message{messageType: messageType_Quit}
	<-this.doneChan
	close(this.doneChan)
}
