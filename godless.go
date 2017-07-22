package pkgthing

import (
	"fmt"

	"github.com/johnny-morrice/godless/api"
	"github.com/johnny-morrice/godless/crdt"
	"github.com/johnny-morrice/godless/query"
)

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
	info, err := readPackageInfo(resp)

	if err != nil {
		return Package{}, err
	}

	if len(info) != 1 {
		return Package{}, fmt.Errorf("Expected exactly 1 PackageInfo but received: %d", len(info))
	}

	pack := Package{
		PackageInfo: info[0],
	}

	return pack, nil
}

func readPackageInfo(resp api.Response) ([]PackageInfo, error) {
	allInfo := []PackageInfo{}

	resp.Namespace.ForeachRow(func(t crdt.TableName, r crdt.RowName, row crdt.Row) {
		// TODO return the error
		system, err := readSystemTableName(t)

		if err != nil {
			return
		}

		info := PackageInfo{
			System: system,
			Name:   string(r),
		}

		allInfo = append(allInfo, info)
	})

	return allInfo, nil
}

func readSystemTableName(tableName crdt.TableName) (string, error) {
	var system string
	_, err := fmt.Scanf(systemTablePrefix+"%s", &system)

	if err != nil {
		return "", err
	}

	return system, nil
}

func systemTable(system string) crdt.TableName {
	return crdt.TableName(systemTablePrefix + system)
}

// TODO should be a method probably.
func metaKey(metaDataKey string) crdt.EntryName {
	return crdt.EntryName(metaDataPrefix + metaDataKey)
}

const systemTablePrefix = "system_"
const metaDataPrefix = "meta_"
