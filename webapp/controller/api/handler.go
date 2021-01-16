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
	if httpsession := session.GetSession(w, r); httpsession != nil {
		if user := httpsession.GetAuthenticatedUser(); user != nil && user.IsStudent() {
			if token := httpsession.GetCSWHToken(); token != "" {
				if r.FormValue("token") == token {
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
					socket := &socket{
						userID:            user.UserID,
						language:          httpsession.GetStringAttribute("Lang"),
						homeworkSessionID: httpsession.GetStringAttribute("HomeworkSessionID"),
						conn:              conn,
					}
					for {
						var messageType int
						var requestMessage []byte
						var err error
						if messageType, requestMessage, err = conn.ReadMessage(); err == nil {
							switch messageType {
							case websocket.TextMessage:
								log.Printf("/websocket: Request received: %s\n", string(requestMessage))
								if err = json.Unmarshal(requestMessage, &socket.message); err == nil {
									if cswhToken, isString := socket.message["token"].(string); isString && cswhToken == token {
										if requestType, isString := socket.message["request"].(string); isString {
											switch requestType {
											case "operation":
												err = socket.operation()
											case "answer":
												err = socket.answer()
											case "toggle":
												err = socket.toggle()
											case "end":
												err = socket.end()
											case "results":
												var homeworkType, status, page int
												if homeworkType, err = socket.getInt("type"); err != nil {
													homeworkType = -1
												}
												if status, err = socket.getInt("status"); err != nil {
													status = -1
												}
												page, _ = socket.getInt("page")
												err = socket.results(homeworkType, status, page)
											case "details":
												if sessionID, isString := socket.message["sessionID"].(string); isString {
													err = socket.details(sessionID)
												}
											default:
												err = fmt.Errorf("/websocket: Invalid request type: %s", requestType)
											}
										} else {
											err = fmt.Errorf("/websocket: Request type is not a string")
										}
									} else {
										err = fmt.Errorf("/websocket: Invalid CSWH token in message")
									}
								}
							case websocket.CloseMessage:
								log.Printf("/websocket: Close message received\n")
								break
							default:
								err = fmt.Errorf("/websocket: Expected a text message, got type %d", messageType)
							}
						}
						if err != nil {
							log.Println(err)
							break
						}
					}
					conn.Close()
					return
				}
				log.Println("/websocket: Invalid CSWH token in URL")
			} else {
				log.Println("/websocket: CSWH token not found in session")
			}
		} else {
			log.Println("/websocket: User is not authenticated or does not have Student role")
		}
	} else {
		log.Println("/websocket: User session not found")
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
