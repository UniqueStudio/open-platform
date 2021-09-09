package pkg

type UniformResource struct {
	State   bool        `json:"state,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ErrorResponse(err error) *UniformResource {
	return &UniformResource{
		State:   false,
		Message: err.Error(),
		Data:    err,
	}
}

func SuccessResponse(data interface{}) *UniformResource {
	return &UniformResource{
		State:   true,
		Message: "ok",
		Data:    data,
	}
}
