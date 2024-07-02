package v1

type PostHandler interface {
}

type postHandler struct {
}

func NewPostHandler() PostHandler {
	return &postHandler{}
}
