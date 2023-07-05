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

type messageType uint8

const (
	messageType_ReadSessionOverview messageType = iota
	messageType_ReadSession
	messageType_NewSession
	messageType_WriteSession
	messageType_NewResponse
)

type message struct {
	messageType messageType
	out         chan *response
	payload     any
}

type readSessionPayload struct {
	sessionId string
}

type writeSessionPayload struct {
	session *data.Session
}

type newResponsePayload struct {
	sessionId string
}

type response struct {
	sessionOverview *data.SessionOverview
	session         *data.Session
	response        *data.Response
	err             *error
}

func worker(sessionStorage *SessionStorage, in chan *message) {
	for message := range in {
		switch message.messageType {
		case messageType_ReadSessionOverview:
			message.out <- &response{
				sessionOverview: sessionStorage.SessionOverview(),
			}

		case messageType_ReadSession:
			payload := message.payload.(*readSessionPayload)
			message.out <- &response{
				session: sessionStorage.ReadSession(payload.sessionId),
			}

		case messageType_NewSession:
			message.out <- &response{
				session: sessionStorage.NewSession(),
			}

		case messageType_WriteSession:
			payload := message.payload.(*writeSessionPayload)
			sessionStorage.WriteSession(payload.session)
			message.out <- &response{}

		case messageType_NewResponse:
			payload := message.payload.(*newResponsePayload)
			newResponse, err := sessionStorage.NewResponse(payload.sessionId)
			message.out <- &response{
				response: newResponse,
				err:      &err,
			}
		}
	}
}

func (this *ConcurrentSessionStorage) SessionOverview() *data.SessionOverview {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_ReadSessionOverview,
		out:         returnChannel,
	}
	response := <-returnChannel
	close(returnChannel)
	return response.sessionOverview
}

func (this *ConcurrentSessionStorage) NewSession() *data.Session {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_NewSession,
		out:         returnChannel,
	}
	response := <-returnChannel
	close(returnChannel)
	return response.session
}

func (this *ConcurrentSessionStorage) ReadSession(id string) *data.Session {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_ReadSession,
		out:         returnChannel,
		payload:     &readSessionPayload{sessionId: id},
	}
	response := <-returnChannel
	close(returnChannel)
	return response.session
}

func (this *ConcurrentSessionStorage) WriteSession(session *data.Session) {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_WriteSession,
		out:         returnChannel,
		payload:     &writeSessionPayload{session: session},
	}
	<-returnChannel
	close(returnChannel)
	return
}

func (this *ConcurrentSessionStorage) NewResponse(sessionId string) (*data.Response, error) {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_NewResponse,
		out:         returnChannel,
		payload:     &newResponsePayload{sessionId: sessionId},
	}
	response := <-returnChannel
	close(returnChannel)
	return response.response, *(response.err)
}
