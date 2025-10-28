package model

type LoteriaInvalidException struct {
	Message string
}

func (e *LoteriaInvalidException) Error() string {
	return e.Message
}

type ResourceNotFoundException struct {
	Message string
}

func (e *ResourceNotFoundException) Error() string {
	return e.Message
}
