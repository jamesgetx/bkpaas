package runtime

// dindHandler is build handler that runs in docker-in-docker mode
type dindHandler struct{}

func (h *dindHandler) GetWorkspace() string {
	return "/tmp/smart-app"
}

func NewDindHandler() (*dindHandler, error) {
	// TODO
	h := &dindHandler{}
	return h, nil
}
