package engine

import (
	"log"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/engine/oobabooga"
	"sedwards2009/llm-workbench/internal/engine/openai"
	"sedwards2009/llm-workbench/internal/engine/request"
)

type Engine struct {
	toWorkerChan      chan *message
	workQueue         []*request.Request
	engineDoneChan    chan bool
	isComputing       bool
	computeWorkerChan chan *request.Request
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

type listModelsPayload struct {
	out chan *data.ModelOverview
}

func NewEngine() *Engine {
	engine := &Engine{
		toWorkerChan:      make(chan *message, 16),
		workQueue:         make([]*request.Request, 0),
		engineDoneChan:    make(chan bool, 16),
		isComputing:       false,
		computeWorkerChan: make(chan *request.Request, 2),
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
				payload := message.payload.(*request.Request)
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

func (this *Engine) computeWorker(in chan *request.Request, done chan bool) {
	for work := range in {
		openai.Process(work)
		done <- true
	}
}

func (this *Engine) scanModels() {
	allModels := []*data.Model{}
	allModels = append(allModels, openai.ScanModels()...)
	allModels = append(allModels, oobabooga.ScanModels()...)

	this.models = allModels
}

func (this *Engine) Enqueue(prompt string, appendFunc func(string), completeFunc func(),
	setStatusFunc func(data.ResponseStatus)) {

	payload := &request.Request{
		Prompt:        prompt,
		AppendFunc:    appendFunc,
		CompleteFunc:  completeFunc,
		SetStatusFunc: setStatusFunc,
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
