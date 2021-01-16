package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/philippecery/maths/webapp/database/dataaccess"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/i18n"
)

type socket struct {
	message           map[string]interface{}
	homeworkSessionID string
	language          string
	userID            string
	conn              *websocket.Conn
}

func (s *socket) get(key string) (interface{}, error) {
	if val, exists := s.message[key]; exists {
		return val, nil
	}
	return nil, fmt.Errorf("Entry %s does not exist", key)
}

func (s *socket) getString(key string) (string, error) {
	if val, isString := s.message[key].(string); isString {
		return val, nil
	}
	return "", fmt.Errorf("Entry %s does not exist or is not a string", key)
}

func (s *socket) getBool(key string) (bool, error) {
	if val, isBool := s.message[key].(bool); isBool {
		return val, nil
	}
	return false, fmt.Errorf("Entry %s does not exist or is not a boolean", key)
}

func (s *socket) getNumber(key string) (float64, error) {
	if val, isFloat64 := s.message[key].(float64); isFloat64 {
		return val, nil
	}
	return 0, fmt.Errorf("Entry %s does not exist or is not a number", key)
}

func (s *socket) getInt(key string) (int, error) {
	if val, err := s.getNumber(key); err == nil {
		return int(val), nil
	}
	return 0, fmt.Errorf("Entry %s does not exist or is not a number", key)
}

func (s *socket) toInt(key string) (int, error) {
	var err error
	var val string
	if val, err = s.getString(key); err == nil {
		var number int
		if number, err = strconv.Atoi(val); err == nil {
			return number, nil
		}
	}
	return 0, err
}

func (s *socket) getCurrentLanguage() string {
	return s.language
}

func (s *socket) addOperation(currentOperation *document.Operation) {
	if session := s.getHomeworkSession(); session != nil {
		session.Operations = append(session.Operations, currentOperation)
		s.saveHomeworkSession(session)
	}
}

func (s *socket) getHomeworkSession() *document.HomeworkSession {
	return dataaccess.GetSessionByID(s.homeworkSessionID)
}

func (s *socket) saveHomeworkSession(homeworkSession *document.HomeworkSession) {
	homeworkSession.EndTime = time.Now()
	dataaccess.UpdateHomeworkSession(homeworkSession)
}

func (s *socket) getLocalizedMessage(messageID string, data ...interface{}) string {
	return i18n.GetLocalizedMessage(s.getCurrentLanguage(), messageID, data...)
}

func (s *socket) emit(messageType int, message interface{}) error {
	var responseMessage []byte
	var err error
	if responseMessage, err = json.Marshal(message); err == nil {
		if err = s.conn.WriteMessage(messageType, responseMessage); err == nil {
			log.Printf("/websocket: Response sent: %s\n", string(responseMessage))
		}
	}
	return err
}

func (s *socket) emitTextMessage(message interface{}) error {
	return s.emit(websocket.TextMessage, message)
}

func (s *socket) emitErrorMessage(messageID string) error {
	return s.emit(websocket.TextMessage, map[string]interface{}{
		"response": "error",
		"message":  s.getLocalizedMessage(messageID),
	})
}
