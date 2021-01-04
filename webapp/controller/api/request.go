package api

import (
	"fmt"
	"strconv"

	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/session"
)

type request struct {
	message map[string]interface{}
	session *session.HTTPSession
}

func (r *request) get(key string) (interface{}, error) {
	if val, exists := r.message[key]; exists {
		return val, nil
	}
	return nil, fmt.Errorf("Entry %s does not exist", key)
}

func (r *request) getString(key string) (string, error) {
	if val, isString := r.message[key].(string); isString {
		return val, nil
	}
	return "", fmt.Errorf("Entry %s does not exist or is nor a string", key)
}

func (r *request) getBool(key string) (bool, error) {
	if val, isBool := r.message[key].(bool); isBool {
		return val, nil
	}
	return false, fmt.Errorf("Entry %s does not exist or is nor a boolean", key)
}

func (r *request) toInt(key string) (int, error) {
	var err error
	var val string
	if val, err = r.getString(key); err == nil {
		var number int
		if number, err = strconv.Atoi(val); err == nil {
			return number, nil
		}
	}
	return 0, err
}

func (r *request) getCurrentLanguage() string {
	if lang, isString := r.session.GetAttribute("Lang").(string); isString {
		return lang
	}
	return "en"
}

func (r *request) addOperation(currentOperation *document.Operation) {
	if session := r.getHomeworkSession(); session != nil {
		session.Operations = append(session.Operations, currentOperation)
		r.saveHomeworkSession(session)
	}
}

func (r *request) getHomeworkSession() *document.HomeworkSession {
	if homeworkSession, isHomeworkSession := r.session.GetAttribute("HomeworkSession").(*document.HomeworkSession); isHomeworkSession {
		return homeworkSession
	}
	return nil
}

func (r *request) saveHomeworkSession(homeworkSession *document.HomeworkSession) {
	r.session.SetAttribute("HomeworkSession", homeworkSession)
}

func (r *request) getLocalizedMessage(messageID string) string {
	return i18n.GetLocalizedMessage(r.getCurrentLanguage(), messageID)
}
