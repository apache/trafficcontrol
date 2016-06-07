package adapter

import (
	"io"
)

type Adapter interface {
	Transform(io.Reader) (interface{}, error)
}
