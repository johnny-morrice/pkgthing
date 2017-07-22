package pkgthing

import (
	"io"
)

type ContentAddressableStorage interface {
	Add(data io.Reader) (string, error)
	Cat(hash string) (io.ReadCloser, error)
}
