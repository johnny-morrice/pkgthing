package pkgthing

import (
	"io"
	"log"

	ipfs "github.com/ipfs/go-ipfs-api"
	godless "github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/http"
)

func MakeIpfsStorage(url string) ContentAddressableStorage {
	return ipfsShell{
		ipfs: ipfs.NewShell(url),
	}
}

func MakeRemoteGodlessClient(url string) (godless.Client, error) {
	options := http.ClientOptions{
		ServerAddr: url,
	}

	return http.MakeClient(options)
}

type ipfsShell struct {
	ipfs *ipfs.Shell
}

func (shell ipfsShell) Cat(hash string) (io.ReadCloser, error) {
	log.Printf("Catting data from IPFS: '%s'", hash)

	return shell.ipfs.Cat(hash)
}

func (shell ipfsShell) Add(r io.Reader) (string, error) {
	log.Printf("Adding data to IPFS...")
	hash, err := shell.ipfs.Add(r)

	if err != nil {
		return "", err
	}

	log.Printf("Added data to IPFS at: '%s'", hash)

	return hash, nil
}
