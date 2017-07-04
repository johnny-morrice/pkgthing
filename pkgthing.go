package pkgthing

type Package struct {
}

type PackageInfo struct {
}

type PackageSearchTerm struct {
}

func Get(deets PackageInfo) (Package, error) {
	return Package{}, nil
}

func Add(pkg Package) (PackageInfo, error) {
	return PackageInfo{}, nil
}

func Search(term PackageSearchTerm) ([]PackageInfo, error) {
	return nil, nil
}
