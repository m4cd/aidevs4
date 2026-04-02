package types

type InputMessage struct {
	SessionID string `json:"sessionID"`
	Msg       string `json:"msg"`
}

type ResponseMessage struct {
	Msg string `json:"msg"`
}

type PackageCheck struct {
	ApiKey    string `json:"apikey"`
	Action    string `json:"action"`
	PackageID string `json:"packageid"`
}

func NewPackageCheck() PackageCheck {
	return PackageCheck{
		Action: "check",
	}
}

type PackageCheckResponse struct {
	OK        bool   `json:"ok"`
	PackageID string `json:"packageid"`
	Location  string `json:"location"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type PackageRedirect struct {
	ApiKey      string `json:"apikey"`
	Action      string `json:"action"`
	PackageID   string `json:"packageid"`
	Destination string `json:"destination"`
	Code        string `json:"code"`
}

func NewPackageRedirect() PackageRedirect {
	return PackageRedirect{
		Action: "redirect",
	}
}

type PackageRedirectResponse struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	Confirmation string `json:"confirmation"`
}

type AnswerEndpointS01E03 struct {
	Url       string `json:"url"`
	SessionID string `json:"sessionID"`
}

type AnswerS01E03 struct {
	Task   string               `json:"task"`
	ApiKey string               `json:"apikey"`
	Answer AnswerEndpointS01E03 `json:"answer"`
}
