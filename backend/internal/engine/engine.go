package engine

import (
	"log"
	"sedwards2009/llm-workbench/internal/data"
)

type Engine struct {
	toWorkerChan      chan *message
	workQueue         []*enqueueWorkPayload
	engineDoneChan    chan bool
	isComputing       bool
	computeWorkerChan chan *enqueueWorkPayload
}

type messageType uint8

const (
	messageType_Enqueue messageType = iota
)

type message struct {
	messageType messageType
	payload     any
}

type enqueueWorkPayload struct {
	prompt        string
	appendFunc    func(string)
	completeFunc  func()
	setStatusFunc func(status data.ResponseStatus)
}

func NewEngine() *Engine {
	engine := &Engine{
		toWorkerChan:      make(chan *message, 16),
		workQueue:         make([]*enqueueWorkPayload, 0),
		engineDoneChan:    make(chan bool, 16),
		isComputing:       false,
		computeWorkerChan: make(chan *enqueueWorkPayload, 2),
	}
	go engine.worker(engine.toWorkerChan)
	return engine
}

func (this *Engine) worker(in chan *message) {
	go this.computeWorker(this.computeWorkerChan, this.engineDoneChan)
	log.Printf("engine worker")

	for {
		select {
		case message := <-in:
			switch message.messageType {
			case messageType_Enqueue:
				payload := message.payload.(*enqueueWorkPayload)
				log.Printf("engine worker: enqueue %p", payload)
				this.workQueue = append(this.workQueue, payload)
				this.tryNextCompute()
			}

		case <-this.engineDoneChan:
			log.Printf("engine worker: compute done")
			this.isComputing = false
			this.tryNextCompute()
		}
	}
}

func (this *Engine) tryNextCompute() {
	if this.isComputing || len(this.workQueue) == 0 {
		return
	}
	nextWork := this.workQueue[0]

	this.workQueue = this.workQueue[1:]
	this.computeWorkerChan <- nextWork
	this.isComputing = true
}

func (this *Engine) computeWorker(in chan *enqueueWorkPayload, done chan bool) {
	for work := range in {
		processOpenAI(work)
		done <- true
	}
}

func (this *Engine) Enqueue(prompt string, appendFunc func(string), completeFunc func(),
	setStatusFunc func(data.ResponseStatus)) {

	payload := &enqueueWorkPayload{
		prompt:        prompt,
		appendFunc:    appendFunc,
		completeFunc:  completeFunc,
		setStatusFunc: setStatusFunc,
	}
	message := &message{
		messageType: messageType_Enqueue,
		payload:     payload,
	}
	this.toWorkerChan <- message
}
