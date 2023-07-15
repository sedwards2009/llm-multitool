package engine

import (
	"log"
	"sedwards2009/llm-workbench/internal/data"

	openai "github.com/sashabaranov/go-openai"
)

type Engine struct {
	toWorkerChan      chan *message
	workQueue         []*enqueueWorkPayload
	engineDoneChan    chan bool
	isComputing       bool
	computeWorkerChan chan *enqueueWorkPayload
	models            []*data.Model
}

type messageType uint8

const (
	messageType_Enqueue messageType = iota
	messageType_ListModels
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

type listModelsPayload struct {
	out chan *data.ModelOverview
}

func NewEngine() *Engine {
	engine := &Engine{
		toWorkerChan:      make(chan *message, 16),
		workQueue:         make([]*enqueueWorkPayload, 0),
		engineDoneChan:    make(chan bool, 16),
		isComputing:       false,
		computeWorkerChan: make(chan *enqueueWorkPayload, 2),
		models:            make([]*data.Model, 0),
	}
	go engine.worker(engine.toWorkerChan)
	return engine
}

func (this *Engine) worker(in chan *message) {
	log.Printf("Engine worker")

	this.scanModels()

	go this.computeWorker(this.computeWorkerChan, this.engineDoneChan)

	for {
		select {
		case message := <-in:
			switch message.messageType {
			case messageType_Enqueue:
				payload := message.payload.(*enqueueWorkPayload)
				log.Printf("engine worker: enqueue %p", payload)
				this.workQueue = append(this.workQueue, payload)
				this.tryNextCompute()

			case messageType_ListModels:
				payload := message.payload.(*listModelsPayload)
				payload.out <- &data.ModelOverview{
					Models: this.models[:],
				}
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

func (this *Engine) scanModels() {
	this.models = []*data.Model{
		{
			ID:              "openai.com_chatgpt3.5turbo",
			Name:            "OpenAI - ChatGPT 3.5 Turbo",
			InternalModelID: openai.GPT3Dot5Turbo,
		},
		{
			ID:              "openai.com_gpt4",
			Name:            "OpenAI - GPT 4",
			InternalModelID: openai.GPT4,
		},
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

func (this *Engine) ModelOverview() *data.ModelOverview {
	returnChannel := make(chan *data.ModelOverview)
	this.toWorkerChan <- &message{
		messageType: messageType_ListModels,
		payload:     &listModelsPayload{out: returnChannel},
	}
	return <-returnChannel
}

func (this *Engine) ValidateModelSettings(modelSettings *data.ModelSettings) bool {
	return this.validateModelID(modelSettings.ModelID)
}

func (this *Engine) validateModelID(modelID string) bool {
	models := this.ModelOverview()
	for _, m := range models.Models {
		if m.ID == modelID {
			return true
		}
	}
	return false
}
