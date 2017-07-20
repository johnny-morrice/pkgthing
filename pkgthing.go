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
	MetaDataKey    string
	MetaDataValues []string
}

type SearchKey uint8

const (
	NameExact = SearchKey(iota)
)

type PackageSearchTerm struct {
	SearchKey  SearchKey
	SearchTerm string
	Keys       []KeyReference
}

type PkgThing struct {
}

// TODO could get directly from IPFS if info.IpfsPath set.
func (thing *PkgThing) Get(info PackageInfo) (Package, error) {
	const failMsg = "Get failed"

	builder := &getBuilder{}
	err := builder.setPackageInfo(info)

	if err != nil {
		return Package{}, errors.Wrap(err, failMsg)
	}

	resp, err := thing.sendQueryWithBuilder(builder)

	thing.logResponse(resp)

	if err != nil {
		return Package{}, errors.Wrap(err, failMsg)
	}

	pkg, err := readPackage(resp)

	if err != nil {
		return Package{}, errors.Wrap(err, failMsg)
	}

	return pkg, nil
}

func (thing *PkgThing) Add(pkg Package) (PackageInfo, error) {
	const failMsg = "Add failed"

	builder := &addBuilder{}
	err := builder.setPackage(pkg)

	if err != nil {
		return PackageInfo{}, errors.Wrap(err, failMsg)
	}

	resp, err := thing.sendQueryWithBuilder(builder)

	thing.logResponse(resp)

	if err != nil {
		return PackageInfo{}, errors.Wrap(err, failMsg)
	}

	path, err := thing.addIpfsBlob(pkg.Blob)

	if err != nil {
		return PackageInfo{}, errors.Wrap(err, failMsg)
	}

	info := pkg.PackageInfo
	info.IpfsPath = path

	return info, nil
}

func (thing *PkgThing) Search(term PackageSearchTerm) ([]PackageInfo, error) {
	const failMsg = "Search failed"

	builder := &searchBuilder{}
	err := builder.setSearchTerm(term)

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

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

func (thing *PkgThing) sendQueryWithBuilder(builder queryBuilder) (api.Response, error) {
	query := builder.buildQuery()
	return thing.sendQuery(query)
}

func (thing *PkgThing) sendQuery(query *query.Query) (api.Response, error) {
	panic("not implemented")
}

func (thing *PkgThing) addIpfsBlob(blob []byte) (string, error) {
	panic("not implemented")
}

func (thing *PkgThing) logResponse(resp api.Response) {
	log.Println(resp)
}

type KeyType uint16

const (
	GODLESS_KEY = KeyType(iota)
)

type SignatureBlob []byte

type KeyFingerprint []byte

type KeyReference struct {
	Type        KeyType
	Fingerprint KeyFingerprint
}

type Signature struct {
	Fingerprint KeyReference
	Data        SignatureBlob
}

type queryBuilder interface {
	buildQuery() *query.Query
}

type addBuilder struct {
}

func (builder *addBuilder) setPackage(pkg Package) error {
	panic("not implemented")
}

func (builder *addBuilder) buildQuery() *query.Query {
	panic("not implemented")
}

type getBuilder struct {
}

func (builder *getBuilder) setPackageInfo(info PackageInfo) error {
	panic("not implemented")
}

func (builder *getBuilder) buildQuery() *query.Query {
	panic("not implemented")
}

type searchBuilder struct {
}

func (builder *searchBuilder) setSearchTerm(term PackageSearchTerm) error {
	panic("not implemented")
}

func (builder *searchBuilder) buildQuery() *query.Query {
	panic("not implemented")
}

func readPackageInfo(resp api.Response) ([]PackageInfo, error) {
	panic("not implemented")
}

func readPackage(resp api.Response) (Package, error) {
	panic("not implemented")
}
