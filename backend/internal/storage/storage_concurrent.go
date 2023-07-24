package storage

import (
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/data/responsestatus"
	"sedwards2009/llm-workbench/internal/data/role"
)

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
	messageType_DeleteSession
	messageType_WriteSession
	messageType_NewResponse
	messageType_DeleteResponse
	messageType_AppendMessage
	messageType_AppendToLastMessage
	messageType_SetResponseStatus
)

type message struct {
	messageType messageType
	out         chan *response
	payload     any
}

type readSessionPayload struct {
	sessionId string
}

type deleteSessionPayload struct {
	sessionId string
}

type writeSessionPayload struct {
	session *data.Session
}

type newResponsePayload struct {
	sessionId string
}

type deleteResponsePayload struct {
	sessionId  string
	responseId string
}

type appendMessagePayload struct {
	sessionId  string
	responseId string
	role       role.Role
	text       string
}

type appendToLastMessagePayload struct {
	sessionId  string
	responseId string
	text       string
}

type setResponseStatusPayload struct {
	sessionId  string
	responseId string
	status     responsestatus.ResponseStatus
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

		case messageType_DeleteSession:
			payload := message.payload.(*deleteSessionPayload)
			sessionStorage.DeleteSession(payload.sessionId)
			message.out <- &response{}

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

		case messageType_DeleteResponse:
			payload := message.payload.(*deleteResponsePayload)
			err := sessionStorage.DeleteResponse(payload.sessionId, payload.responseId)
			message.out <- &response{
				err: &err,
			}

		case messageType_AppendMessage:
			payload := message.payload.(*appendMessagePayload)
			err := sessionStorage.AppendMessage(payload.sessionId, payload.responseId, payload.role, payload.text)
			message.out <- &response{
				err: &err,
			}

		case messageType_AppendToLastMessage:
			payload := message.payload.(*appendToLastMessagePayload)
			err := sessionStorage.AppendToLastMessage(payload.sessionId, payload.responseId, payload.text)
			message.out <- &response{
				err: &err,
			}

		case messageType_SetResponseStatus:
			payload := message.payload.(*setResponseStatusPayload)
			err := sessionStorage.SetResponseStatus(payload.sessionId, payload.responseId, payload.status)
			message.out <- &response{
				err: &err,
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

func (this *ConcurrentSessionStorage) DeleteSession(id string) {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_DeleteSession,
		out:         returnChannel,
		payload:     &deleteSessionPayload{sessionId: id},
	}
	<-returnChannel
	close(returnChannel)
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

func (this *ConcurrentSessionStorage) DeleteResponse(sessionId string, responseId string) error {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_DeleteResponse,
		out:         returnChannel,
		payload:     &deleteResponsePayload{sessionId: sessionId, responseId: responseId},
	}
	response := <-returnChannel
	close(returnChannel)
	return *(response.err)
}

func (this *ConcurrentSessionStorage) AppendMessage(sessionId string, responseId string, messageRole role.Role, text string) error {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_AppendMessage,
		out:         returnChannel,
		payload:     &appendMessagePayload{sessionId: sessionId, responseId: responseId, role: messageRole, text: text},
	}
	response := <-returnChannel
	close(returnChannel)
	return *(response.err)
}

func (this *ConcurrentSessionStorage) AppendToLastMessage(sessionId string, responseId string, text string) error {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_AppendToLastMessage,
		out:         returnChannel,
		payload:     &appendToLastMessagePayload{sessionId: sessionId, responseId: responseId, text: text},
	}
	response := <-returnChannel
	close(returnChannel)
	return *(response.err)
}

func (this *ConcurrentSessionStorage) SetResponseStatus(sessionId string, responseId string, status responsestatus.ResponseStatus) error {
	returnChannel := make(chan *response)
	this.toWorker <- &message{
		messageType: messageType_SetResponseStatus,
		out:         returnChannel,
		payload:     &setResponseStatusPayload{sessionId: sessionId, responseId: responseId, status: status},
	}
	response := <-returnChannel
	close(returnChannel)
	return *(response.err)
}
