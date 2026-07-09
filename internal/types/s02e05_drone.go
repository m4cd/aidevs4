package types

type AnswerInstructionsS02E05 struct {
	Instructions []string `json:"instructions"`
}

type AnswerS02E05 struct {
	Task   string              `json:"task"`
	ApiKey string              `json:"apikey"`
	Answer AnswerInstructionsS02E05 `json:"answer"`
}