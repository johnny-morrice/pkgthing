package pkgthing

import (
	"fmt"
	"log"

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
		__DATAPATH_KEY: crdt.PointText(builder.pack.IpfsPath),
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

func (builder *getBuilder) buildQuery() (*query.Query, error) {
	info := builder.info
	return query.Compile("select ? where str_eq(@key, ?)", systemTable(info.System), info.Name)
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
	case SEARCH_NAME:
		return builder.nameWildcardQuery()
	default:
		return nil, fmt.Errorf("Unknown SearchKey: %v", builder.term.SearchKey)
	}
}

func (builder *searchBuilder) systemQuery() (*query.Query, error) {
	return query.Compile("select ?", systemTable(builder.term.System))
}

// FIXME Sprintf is security hole. We need parametrized queries from godless 0.19.0.
func (builder *searchBuilder) nameWildcardQuery() (*query.Query, error) {
	return query.Compile("select ? where str_glob(@key, ?)", systemTable(builder.term.System), builder.term.SearchTerm)
}

func readPackage(resp api.Response) (Package, error) {
	info, err := readPackageInfo(resp)

	log.Print("Read package info")

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

		dataentry, err := row.GetEntry(__DATAPATH_KEY)

		if err != nil {
			return
		}

		points := dataentry.GetValues()

		if len(points) != 1 {
			return
		}

		datapath := points[0].Text()

		info := PackageInfo{
			System:   system,
			Name:     string(r),
			IpfsPath: string(datapath),
		}

		allInfo = append(allInfo, info)
	})

	return allInfo, nil
}

func readSystemTableName(tableName crdt.TableName) (string, error) {
	var system string
	_, err := fmt.Sscanf(string(tableName), __SYSTEM_TABLE_PREFIX+"%s", &system)

	if err != nil {
		return "", err
	}

	return system, nil
}

func systemTable(system string) crdt.TableName {
	if system == "" {
		panic("BUG system was empty")
	}

	return crdt.TableName(__SYSTEM_TABLE_PREFIX + system)
}

// TODO should be a method probably.
func metaKey(metaDataKey string) crdt.EntryName {
	if metaDataKey == "" {
		panic("BUG metaDataKey was empty")
	}

	return crdt.EntryName(__META_DATA_PREFIX + metaDataKey)
}

const __DATAPATH_KEY = "datapath"
const __SYSTEM_TABLE_PREFIX = "system_"
const __META_DATA_PREFIX = "meta_"
