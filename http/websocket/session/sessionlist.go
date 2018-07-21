

package session

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

const MAX_SESSION_COUNT = 3000

type SessionList struct {
	sync.RWMutex
	mapOnlineList map[string]*Session //key is SessionId
}

func NewSessionList() *SessionList {
	return &SessionList{
		mapOnlineList: make(map[string]*Session),
	}
}
func (self *SessionList) NewSession(wsConn *websocket.Conn) (session *Session, err error) {
	if self.GetSessionCount() > MAX_SESSION_COUNT {
		return nil, errors.New("over MAX_SESSION_COUNT")
	}
	session = newSession(wsConn)

	self.Lock()
	self.mapOnlineList[session.GetSessionId()] = session
	self.Unlock()

	return session, nil
}
func (self *SessionList) CloseSession(session *Session) {
	if session == nil {
		return
	}
	self.removeSession(session)
	session.Close()
}

func (self *SessionList) removeSession(session *Session) {
	self.Lock()
	defer self.Unlock()
	delete(self.mapOnlineList, session.GetSessionId())
}

func (self *SessionList) GetSessionById(sessionId string) *Session {
	self.RLock()
	defer self.RUnlock()
	return self.mapOnlineList[sessionId]

}

func (self *SessionList) GetSessionCount() int {
	self.RLock()
	defer self.RUnlock()
	return len(self.mapOnlineList)
}

func (self *SessionList) ForEachSession(visit func(*Session)) {
	self.RLock()
	defer self.RUnlock()
	for _, v := range self.mapOnlineList {
		visit(v)
	}
}
