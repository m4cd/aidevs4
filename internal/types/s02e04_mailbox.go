package types

import "encoding/json"

type MailboxHelp struct {
	ApiKey string `json:"apikey"`
	Action string `json:"action"`
	Page   int    `json:"page"`
}

type MailboxGetInbox struct {
	ApiKey  string `json:"apikey"`
	Action  string `json:"action"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

type MailboxSearch struct {
	ApiKey  string `json:"apikey"`
	Action  string `json:"action"`
	Query   string `json:"query"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

type MailboxGetMessages struct {
	ApiKey string `json:"apikey"`
	Action string `json:"action"`
	Ids    IDList `json:"ids"`
}

type IDList []string

func (l *IDList) UnmarshalJSON(data []byte) error {
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		arr = []json.RawMessage{data} // bare scalar -> single-element
	}
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		s := string(e)
		if len(s) >= 2 && s[0] == '"' { // strip quotes if it was a JSON string
			var str string
			if err := json.Unmarshal(e, &str); err != nil {
				return err
			}
			s = str
		}
		out = append(out, s)
	}
	*l = out
	return nil
}

type MailboxGetThread struct {
	ApiKey string `json:"apikey"`
	Action string `json:"action"`
	ThreadId    IDList `json:"threadID"`
}

type AnswerMessageS02E04 struct {
	Password         string `json:"password"`
	Date             string `json:"date"`
	ConfirmationCode string `json:"confirmation_code"`
}

type AnswerS02E04 struct {
	Task   string              `json:"task"`
	ApiKey string              `json:"apikey"`
	Answer AnswerMessageS02E04 `json:"answer"`
}
