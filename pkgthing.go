package pkgthing

import (
	"fmt"
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
	SEARCH_NAME_WILDCARD = SearchKey(iota)
	SEARCH_SYSTEM
)

type PackageSearchTerm struct {
	SearchKey  SearchKey
	SearchTerm string
	System     string
	Keys       []KeyReference
}

type PkgThing struct {
}

// TODO could get directly from IPFS if info.IpfsPath set.
func (thing *PkgThing) Get(info PackageInfo) (Package, error) {
	const failMsg = "Get failed"

	builder := &getBuilder{}
	builder.setPackageInfo(info)

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
	builder.setPackage(pkg)

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

func (thing *PkgThing) sendQueryWithBuilder(builder queryBuilder) (api.Response, error) {
	query, err := builder.buildQuery()

	if err != nil {
		return api.RESPONSE_FAIL, err
	}

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
	buildQuery() (*query.Query, error)
}

type addBuilder struct {
	pkg Package
}

func (builder *addBuilder) setPackage(pkg Package) {
	builder.pkg = pkg
}

// FIXME Sprintf is security hole. We need parametrized queries from godless 0.19.0.
func (builder *addBuilder) buildQuery() (*query.Query, error) {
	pkg := builder.pkg
	queryFormat := "join %s rows (@key=%s, metadata=\"%s\")"
	metadata := encodeMetaDataAsText(pkg.MetaData)
	queryText := fmt.Sprintf(queryFormat, packageTable(pkg.System), pkg.Name, metadata)
	return query.Compile(queryText)
}

type getBuilder struct {
	info PackageInfo
}

func (builder *getBuilder) setPackageInfo(info PackageInfo) {
	builder.info = info
}

// FIXME Sprintf is security hole. We need parametrized queries from godless 0.19.0.
func (builder *getBuilder) buildQuery() (*query.Query, error) {
	info := builder.info
	queryFormat := "select %s where str_eq(@key, \"%s\")"
	queryText := fmt.Sprintf(queryFormat, packageTable(info.System), info.Name)
	return query.Compile(queryText)
}

type searchBuilder struct {
	term PackageSearchTerm
}

func (builder *searchBuilder) setSearchTerm(term PackageSearchTerm) {
	builder.term = term
}

func (builder *searchBuilder) buildQuery() (*query.Query, error) {
	switch builder.term.SearchKey {
	case SEARCH_SYSTEM:
		return builder.systemQuery()
	case SEARCH_NAME_WILDCARD:
		return builder.exactNameQuery()
	default:
		return nil, fmt.Errorf("Unknown SearchKey: %v", builder.term.SearchKey)
	}
}

// FIXME Sprintf is security hole. We need parametrized queries from godless 0.19.0.
func (builder *searchBuilder) systemQuery() (*query.Query, error) {
	queryFormat := "select %s"
	queryText := fmt.Sprintf(queryFormat, packageTable(builder.term.System))
	return query.Compile(queryText)
}

// FIXME Sprintf is security hole. We need parametrized queries from godless 0.19.0.
func (builder *searchBuilder) exactNameQuery() (*query.Query, error) {
	queryFormat := "select %s where str_wildcard(@key, \"%s\")"
	queryText := fmt.Sprintf(queryFormat, packageTable(builder.term.System), builder.term.SearchTerm)
	return query.Compile(queryText)
}

func readPackageInfo(resp api.Response) ([]PackageInfo, error) {
	panic("not implemented")
}

func readPackage(resp api.Response) (Package, error) {
	panic("not implemented")
}

func encodeMetaDataAsText(metaData []MetaDataEntry) string {
	panic("not implemented")
}

func packageTable(system string) string {
	return "package_" + system
}
