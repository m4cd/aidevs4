package types

// type AnswerActionS01E05 struct {
// 	Action string `json:"action"`
// }

type AnswerS01E05 struct {
	Task   string `json:"task"`
	ApiKey string `json:"apikey"`
	Answer any    `json:"answer"`
}

type HelpInput struct {
	Action string `json:"action"` // "help"
}

type ReconfigureInput struct {
	Action string `json:"action"` // "reconfigure"
	Route  string `json:"route"`
}

type GetStatusInput struct {
	Action string `json:"action"` // "getstatus"
	Route  string `json:"route"`
}

type SaveInput struct {
	Action string `json:"action"` // "save"
	Route  string `json:"route"`
}
type SetStatusInput struct {
	Action string `json:"action"` // "setstatus"
	Route  string `json:"route"`
	Value  string `json:"value"` // "RTOPEN" | "RTCLOSE"
}


type SuccessInput struct {
	Flag string `json:"flag"` // "help"
}