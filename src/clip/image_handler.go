package clip

import (
	"strings"
)

type ImageHandler struct {
	// Put field definition here
	Handler
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{}
}

func (h ImageHandler) Test(grab *Grab, contentType string) bool {
	return strings.HasPrefix(contentType, "image")
}

func (h ImageHandler) Handle(me *Grab) (err error) {
	Logr.Info("Serving from image handler")
	width, height := GetImageSize(me.Url)
	if width != "0" && height != "0" {
		me.ImageWidth       = width
		me.ImageHeight      = height
		me.MediaType        = "photo"
		me.ImageThumbUrl    = me.Url
		me.wrapURL(&me.ImageThumbUrl, true, &me.ImageThumbUrlStatus)

	}
	return nil
}