package pkgthing

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"

	"github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/query"
)

type Package struct {
	PackageInfo
	Data []byte
}

type PackageInfo struct {
	Name       string
	System     string
	IpfsPath   string
	MetaData   []MetaDataEntry
	Signatures []Signature
}

type MetaDataEntry struct {
	MetaDataKey   string
	MetaDataValue string
}

type SearchMethod uint8

const (
	SEARCH_WILDCARD = SearchMethod(iota)
)

type SearchKey uint8

const (
	SEARCH_NAME = SearchKey(iota)
	SEARCH_SYSTEM
)

func ParseSearchKey(key string) (SearchKey, error) {
	switch key {
	case "name":
		return SEARCH_NAME, nil
	case "system":
		return SEARCH_SYSTEM, nil
	default:
		return 0, fmt.Errorf("Unknown SearchKey: %s", key)
	}
}

type PackageSearchTerm struct {
	SearchKey    SearchKey
	SearchMethod SearchMethod
	SearchTerm   string
	System       string
	Keys         []KeyReference
}

type PackageManager interface {
	Get(info PackageInfo) (Package, error)
	Add(pack Package) (PackageInfo, error)
	Search(term PackageSearchTerm) ([]PackageInfo, error)
}

type Options struct {
	Store   ContentAddressableStorage
	Godless api.Client
}

func New(options Options) PackageManager {
	return &pkgthing{
		Options: options,
	}
}

type pkgthing struct {
	Options
}

func (thing *pkgthing) Get(info PackageInfo) (Package, error) {
	const failMsg = "Get failed"

	builder := &getBuilder{}
	builder.setPackageInfo(info)

	resp, err := thing.sendQueryWithBuilder(builder)

	thing.logResponse(resp)

	if err != nil {
		return Package{}, errors.Wrap(err, failMsg)
	}

	pack, err := readPackage(resp)

	if err != nil {
		return Package{}, errors.Wrap(err, failMsg)
	}

	log.Printf("Found package: %v", pack)

	err = thing.loadPackageData(&pack)

	if err != nil {
		return Package{}, errors.Wrap(err, failMsg)
	}

	return pack, nil
}

func (thing *pkgthing) Add(pack Package) (PackageInfo, error) {
	const failMsg = "Add failed"

	path, err := thing.addIpfsBlob(pack.Data)

	if err != nil {
		return PackageInfo{}, errors.Wrap(err, failMsg)
	}

	pack.IpfsPath = path
	builder := &addBuilder{}
	builder.setPackage(pack)

	resp, err := thing.sendQueryWithBuilder(builder)

	thing.logResponse(resp)

	if err != nil {
		return PackageInfo{}, errors.Wrap(err, failMsg)
	}

	info := pack.PackageInfo
	info.IpfsPath = path

	return info, nil
}

func (thing *pkgthing) Search(term PackageSearchTerm) ([]PackageInfo, error) {
	const failMsg = "Search failed"

	builder := &searchBuilder{}
	builder.setSearchTerm(term)

	resp, err := thing.sendQueryWithBuilder(builder)

	thing.logResponse(resp)

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

	info, err := readPackageInfo(resp)

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

	return info, nil
}

func (thing *pkgthing) loadPackageData(pack *Package) error {
	reader, err := thing.Store.Cat(pack.IpfsPath)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	go func() {
		err := reader.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	pack.Data = data
	return nil
}

func (thing *pkgthing) sendQueryWithBuilder(builder queryBuilder) (api.Response, error) {
	query, err := builder.buildQuery()

	if err != nil {
		return api.RESPONSE_FAIL, err
	}

	return thing.sendQuery(query)
}

func (thing *pkgthing) sendQuery(query *query.Query) (api.Response, error) {
	request := api.MakeQueryRequest(query)
	return thing.Godless.Send(request)
}

func (thing *pkgthing) addIpfsBlob(blob []byte) (string, error) {
	reader := bytes.NewReader(blob)
	return thing.Store.Add(reader)
}

func (thing *pkgthing) logResponse(resp api.Response) {
	log.Println(resp)
}
