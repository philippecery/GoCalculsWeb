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
				http.Redirect(w, r, "/student/dashboard", http.StatusFound)
				return
			}
			log.Println("/websocket: Client Connected")
			request := &request{session: httpsession}
			for {
				var messageType int
				var requestMessage, responseMessage []byte
				var err error
				if messageType, requestMessage, err = conn.ReadMessage(); err == nil {
					log.Printf("/websocket: Request received: %s\n", string(requestMessage))
					if err = json.Unmarshal(requestMessage, &request.message); err == nil {
						if csrfToken, isString := request.message["token"].(string); isString && csrfToken == token {
							var response interface{}
							if requestType, isString := request.message["request"].(string); isString {
								switch requestType {
								case "operation":
									response = request.operation()
								case "answer":
									response = request.answer()
								case "toggle":
									response = request.toggle()
								case "results":
									response = request.results()
								default:
									err = fmt.Errorf("/websocket: Invalid message type: %s", requestType)
								}
								if response != nil {
									if responseMessage, err = json.Marshal(response); err == nil {
										if err = conn.WriteMessage(messageType, responseMessage); err == nil {
											log.Printf("/websocket: Response sent: %s\n", string(responseMessage))
										}
									}
								}
							} else {
								err = fmt.Errorf("/websocket: Request type is not a string")
							}
						} else {
							err = fmt.Errorf("/websocket: Invalid CSRF token")
						}
					}
				}
				if err != nil {
					log.Println(err)
					break
				}
			}
			conn.Close()
		} else {
			log.Println("/websocket: CSRF token not found in session")
		}
	} else {
		log.Println("/websocket: User is not authenticated or does not have Student role")
	}
	log.Println("/websocket: Redirecting to Login page")
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
