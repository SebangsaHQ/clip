package clip

// Handler has responsibility to gather information from various sources
// For example HTML, image, video, pdf, doc, etc
// Every type of handler should implement Handler interface

type Handler interface {
	Test(grab *Grab, contentType string) bool
	Handle(grab *Grab) error
}

// List of handler
// Put your own handler in this list
var HandlerList = []Handler {
	NewImageHandler(),
	NewHtmlHandler(),
}

// Figure out how to choose right handler for content type
// SelectHandler returns error that will be delegated to Grabber
func SelectHandler(grab *Grab, contentType string) error {
	for _, handler := range HandlerList {
		if (handler.Test(grab, contentType)) {
			return handler.Handle(grab)
		}
	}

	return nil
}