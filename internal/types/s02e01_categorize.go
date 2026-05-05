package types

type AnswerActionS02E01 struct {
	Prompt string `json:"prompt"`
}

type AnswerS02E01 struct {
	Task   string             `json:"task"`
	ApiKey string             `json:"apikey"`
	Answer AnswerActionS02E01 `json:"answer"`
}
