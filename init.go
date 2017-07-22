package pkgthing

import (
	ipfs "github.com/ipfs/go-ipfs-api"
	godless "github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/http"
)

func MakeIpfsStorage(url string) ContentAddressableStorage {
	return ipfs.NewShell(url)
}

func MakeRemoteGodlessClient(url string) (godless.Client, error) {
	options := http.ClientOptions{
		ServerAddr: url,
	}

	return http.MakeClient(options)
}
