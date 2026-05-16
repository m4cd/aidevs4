package types

type AnswerLogsS02E03 struct {
	Logs string `json:"logs"`
}

type AnswerS02E03 struct {
	Task   string           `json:"task"`
	ApiKey string           `json:"apikey"`
	Answer AnswerLogsS02E03 `json:"answer"`
}
