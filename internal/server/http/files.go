package http

import (
	"fmt"

	"github.com/bnkamalesh/errors"
	"github.com/h2non/filetype"
	svg "github.com/h2non/go-is-svg"
)

func detectFileType(content []byte) (string, error) {
	if svg.Is(content) {
		return "image/svg+xml", nil
	}

	kind, err := filetype.Match(content)
	if kind == filetype.Unknown || err != nil {
		if err != nil {
			return "", err
		}
		return "", errors.Validation("unknown file type")
	}

	return fmt.Sprintf("%s/%s", kind.MIME.Type, kind.MIME.Subtype), nil
}
