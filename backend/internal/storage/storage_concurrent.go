package storage

import "sedwards2009/llm-workbench/internal/data"

type ConcurrentSessionStorage struct {
	toWorker chan *message
}

func NewConcurrentSessionStorage(storagePath string) *ConcurrentSessionStorage {
	sessionStorage := NewSessionStorage(storagePath)
	instance := &ConcurrentSessionStorage{
		toWorker: make(chan *message, 16),
	}
	sessionStorage.Scan()

	go worker(sessionStorage, instance.toWorker)

	return instance
}

type message struct {
	readSessionOverview bool
	readSession         string
	newSession          bool
	writeSession        *data.Session
	out                 chan *response
}

type response struct {
	sessionOverview *data.SessionOverview
	session         *data.Session
}

func worker(sessionStorage *SessionStorage, in chan *message) {
	for message := range in {
		if message.readSessionOverview {
			message.out <- &response{
				sessionOverview: sessionStorage.SessionOverview(),
			}
		}
		if message.newSession {
			message.out <- &response{
				session: sessionStorage.NewSession(),
			}
		}
		if message.readSession != "" {
			message.out <- &response{
				session: sessionStorage.ReadSession(message.readSession),
			}
		}
		if message.writeSession != nil {
			sessionStorage.WriteSession(message.writeSession)
			message.out <- &response{}
		}
	}
}

func (this *ConcurrentSessionStorage) SessionOverview() *data.SessionOverview {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		readSessionOverview: true,
		out:                 returnChannel,
	}
	response := <-returnChannel
	return response.sessionOverview
}

func (this *ConcurrentSessionStorage) NewSession() *data.Session {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		newSession: true,
		out:        returnChannel,
	}
	response := <-returnChannel
	return response.session
}

func (this *ConcurrentSessionStorage) ReadSession(id string) *data.Session {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		readSession: id,
		out:         returnChannel,
	}
	response := <-returnChannel
	return response.session
}

func (this *ConcurrentSessionStorage) WriteSession(session *data.Session) {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		writeSession: session,
		out:          returnChannel,
	}
	<-returnChannel
	return
}
