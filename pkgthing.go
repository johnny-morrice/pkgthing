package pkgthing

import (
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/crdt"
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

type PkgThing struct {
}

func (thing *PkgThing) Get(info PackageInfo) (Package, error) {
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

func (thing *PkgThing) Add(pack Package) (PackageInfo, error) {
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

func (thing *PkgThing) loadPackageData(pack *Package) error {
	panic("not implemented")
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

func (thing *PkgThing) logIpfsPath(path string) {
	log.Printf("Added to IPFS at: %s", path)
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
	pack Package
}

func (builder *addBuilder) setPackage(pack Package) {
	builder.pack = pack
}

// FIXME working directly with query structure is awful.
func (builder *addBuilder) buildQuery() (*query.Query, error) {
	table := systemTable(builder.pack.System)
	rowKey := builder.pack.Name

	entries := map[crdt.EntryName]crdt.PointText{
		"datapath": crdt.PointText(builder.pack.IpfsPath),
	}
	row := query.QueryRowJoin{
		RowKey:  crdt.RowName(rowKey),
		Entries: entries,
	}

	for _, meta := range builder.pack.MetaData {
		key := metaKey(meta.MetaDataKey)
		entries[key] = crdt.PointText(meta.MetaDataValue)
	}

	q := &query.Query{
		OpCode:   query.JOIN,
		TableKey: crdt.TableName(table),
		Join: query.QueryJoin{
			Rows: []query.QueryRowJoin{row},
		},
	}
	return q, nil
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
	queryText := fmt.Sprintf(queryFormat, systemTable(info.System), info.Name)
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
	queryText := fmt.Sprintf(queryFormat, systemTable(builder.term.System))
	return query.Compile(queryText)
}

// FIXME Sprintf is security hole. We need parametrized queries from godless 0.19.0.
func (builder *searchBuilder) exactNameQuery() (*query.Query, error) {
	queryFormat := "select %s where str_wildcard(@key, \"%s\")"
	queryText := fmt.Sprintf(queryFormat, systemTable(builder.term.System), builder.term.SearchTerm)
	return query.Compile(queryText)
}

func readPackage(resp api.Response) (Package, error) {
	panic("not implemented")
}

func readPackageInfo(resp api.Response) ([]PackageInfo, error) {
	panic("not implemented")
}

func systemTable(system string) crdt.TableName {
	return crdt.TableName("system" + system)
}

// TODO should be a method probably.
func metaKey(metaDataKey string) crdt.EntryName {
	return crdt.EntryName("meta_" + metaDataKey)
}
