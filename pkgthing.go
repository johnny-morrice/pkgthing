package pkgthing

import (
	"log"

	"github.com/pkg/errors"

	"github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/query"
)

type Package struct {
	PackageInfo
	Blob []byte
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

type SearchKey uint8

const (
	SEARCH_NAME_WILDCARD = SearchKey(iota)
	SEARCH_SYSTEM
)

type PackageSearchTerm struct {
	SearchKey  SearchKey
	SearchTerm string
	System     string
	Keys       []KeyReference
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

	err = thing.loadPackageData(&pack)

	if err != nil {
		return Package{}, errors.Wrap(err, failMsg)
	}

	return pack, nil
}

func (thing *pkgthing) Add(pack Package) (PackageInfo, error) {
	const failMsg = "Add failed"

	path, err := thing.addIpfsBlob(pack.Blob)

	if err != nil {
		return PackageInfo{}, errors.Wrap(err, failMsg)
	}

	thing.logIpfsPath(path)

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
	panic("not implemented")
}

func (thing *pkgthing) sendQueryWithBuilder(builder queryBuilder) (api.Response, error) {
	query, err := builder.buildQuery()

	if err != nil {
		return api.RESPONSE_FAIL, err
	}

	return thing.sendQuery(query)
}

func (thing *pkgthing) sendQuery(query *query.Query) (api.Response, error) {
	panic("not implemented")
}

func (thing *pkgthing) addIpfsBlob(blob []byte) (string, error) {
	panic("not implemented")
}

func (thing *pkgthing) logResponse(resp api.Response) {
	log.Println(resp)
}

func (thing *pkgthing) logIpfsPath(path string) {
	log.Printf("Added to IPFS at: %s", path)
}
