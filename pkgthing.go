package pkgthing

import (
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

func (thing *PkgThing) Get(info PackageInfo) (Package, error) {
	panic("not implemented")
}

func (thing *PkgThing) Add(pkg Package) (PackageInfo, error) {
	panic("not implemented")
}

func (thing *PkgThing) Search(term PackageSearchTerm) ([]PackageInfo, error) {
	const failMsg = "Search failed"

	builder := &searchBuilder{}
	err := builder.setSearchTerm(term)

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

	query := builder.buildQuery()

	resp, err := thing.sendQuery(query)

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

	info, err := readPackageInfo(resp)

	if err != nil {
		return nil, errors.Wrap(err, failMsg)
	}

	return info, nil
}

func (thing *PkgThing) sendQuery(query *query.Query) (api.Response, error) {
	panic("not implemented")
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
