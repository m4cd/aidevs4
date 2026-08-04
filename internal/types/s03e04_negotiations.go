package types

type AnswerS03E04 struct {
	Task   string            `json:"task"`
	ApiKey string            `json:"apikey"`
	Answer AnswerToolsS03E04 `json:"answer"`
}

type AnswerToolsS03E04 struct {
	Tools []ToolS0304 `json:"tools"`
}

type ToolS0304 struct {
	Url         string `json:"URL"`
	Description string `json:"description"`
}

type AnswerS03E04Check struct {
	Task   string      `json:"task"`
	ApiKey string      `json:"apikey"`
	Answer CheckS03E04 `json:"answer"`
}

type CheckS03E04 struct {
	Action string `json:"action"`
}

type InputFromAgent struct {
	Params string `json:"params"`
}

type OutputToAgent struct {
	Output string `json:"output"`
}

type ItemS03E04 struct {
	Name string `json:"name"`
	Parameters []ItemParamS03E04 `json:"parameters"`
}

type ItemParamS03E04 struct{
	Value string `json:"value"`
}