package pkgthing

import (
	"io"
)

type ContentAddressableStorage interface {
	Add(data io.Reader) string
	Cat(hash string) io.Reader
}
