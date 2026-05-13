package types

type AnswerRotateS02E02 struct {
	Rotate string `json:"rotate"`
}

type AnswerS02E02 struct {
	Task   string             `json:"task"`
	ApiKey string             `json:"apikey"`
	Answer AnswerRotateS02E02 `json:"answer"`
}
