package types

type AnswerDeclarationS01E04 struct {
	Declaration string `json:"declaration"`
}

type AnswerS01E04 struct {
	Task   string                  `json:"task"`
	ApiKey string                  `json:"apikey"`
	Answer AnswerDeclarationS01E04 `json:"answer"`
}
