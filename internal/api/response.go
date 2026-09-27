package api

type APIResponse[T any] struct {
	Result  *T    `json:"result,omitempty"`
	Message string `json:"message,omitempty"`
}

func Success[T any](result T) APIResponse[T] {
	return APIResponse[T]{
		Result: &result,
	}
}

func Error(message string) APIResponse[any] {
	return APIResponse[any]{
		Message: message,
	}
}