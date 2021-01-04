package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/philippecery/maths/webapp/session"
)

// Endpoints controller
func Endpoints(w http.ResponseWriter, r *http.Request) {
	httpsession := session.GetSession(w, r)
	if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsStudent() {
		if token := httpsession.GetCSRFToken(); token != "" {
			upgrader := websocket.Upgrader{}
			upgrader.CheckOrigin = func(r *http.Request) bool {
				origin := r.Header["Origin"]
				if len(origin) == 0 {
					return false
				}
				u, err := url.Parse(origin[0])
				if err != nil {
					return false
				}
				return equalASCIIFold(u.Host, r.Host)
			}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
			}
			log.Println("Client Connected")
			request := &request{session: httpsession}
			for {
				var messageType int
				var requestMessage, responseMessage []byte
				var err error
				if messageType, requestMessage, err = conn.ReadMessage(); err == nil {
					log.Printf("Request received: %s\n", string(requestMessage))
					if err = json.Unmarshal(requestMessage, &request.message); err == nil {
						var response interface{}
						if requestType, isString := request.message["request"].(string); isString {
							switch requestType {
							case "operation":
								fmt.Println("generating next operation")
								response = request.generateNextOperation()
							case "answer":
								fmt.Println("processing answer")
								response = request.processAnswer()
							case "toggle":
								fmt.Println("toggling answer/result")
								response = request.toggleResult()
							case "results":
								fmt.Println("displaying final results")
								response = request.results()
							default:
								log.Printf("Invalid message type: %s", requestType)
							}
							if responseMessage, err = json.Marshal(response); err == nil {
								if err = conn.WriteMessage(messageType, responseMessage); err == nil {
									log.Printf("Response sent: %s\n", string(responseMessage))
								}
							}
						} else {
							err = fmt.Errorf("Request type is not a string")
						}
					}
				}
				if err != nil {
					log.Println(err)
					return
				}
			}
		} else {
			log.Println("CSRF token not found in session")
		}
	} else {
		log.Println("User is not authenticated or does not have Student role")
	}
	log.Println("Redirecting to Login page")
	http.Redirect(w, r, "/logout", http.StatusFound)
}

func equalASCIIFold(s, t string) bool {
	for s != "" && t != "" {
		sr, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		tr, size := utf8.DecodeRuneInString(t)
		t = t[size:]
		if sr == tr {
			continue
		}
		if 'A' <= sr && sr <= 'Z' {
			sr = sr + 'a' - 'A'
		}
		if 'A' <= tr && tr <= 'Z' {
			tr = tr + 'a' - 'A'
		}
		if sr != tr {
			return false
		}
	}
	return s == t
}
