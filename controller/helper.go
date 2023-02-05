package controller

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func getSuccessResponse() *response {
	return &response{
		Status: success,
	}
}

func GetFailResponse(err string) *response {
	return &response{
		Status: fail,
		Error:  err,
	}
}
