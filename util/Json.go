package util

type ResultMessage struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Version int         `json:"version"`
	Data    interface{} `json:"data"`
}

func GetSuccessJson(message string) (result *ResultMessage) {
	result = &ResultMessage{
		Status:  "success",
		Message: message,
		Version: 1,
	}
	return
}
func GetSuccessReJson(data interface{}) (result *ResultMessage) {
	result = &ResultMessage{
		Status:  "success",
		Data:    data,
		Version: 1,
	}
	return
}

func GetErrorJson(message string) (result *ResultMessage) {
	result = &ResultMessage{
		Status:  "error",
		Message: message,
		Version: 1,
	}
	return
}

func GetFailJson(message string) (result *ResultMessage) {
	result = &ResultMessage{
		Status:  "fail",
		Message: message,
		Version: 1,
	}
	return
}
