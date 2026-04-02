package main

import "fmt"

type Session struct {
	ID string
	Messages []string
}

func (s *Session) Print() {
	for _, msg := range s.Messages {
		fmt.Println(msg)
	}
}
func SessionsToString(sessions map[string]*Session) string {
	s := ""
	for _, session := range sessions {
		s = s + fmt.Sprintf("SessionID: %v\n", session.ID)
		for _, msg := range session.Messages {
			s = s + "[+] " + msg + "\n"
		}
	}
	return s
}
