package pkgthing

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/pkg/errors"
)

type Ubuntu struct {
}

func (ubuntu *Ubuntu) GetInstalledPackages() ([]PackageInfo, error) {
	const errMsg = "Ubuntu.GetInstalledPackages failed"

	cmd := exec.Command(__DPKG_COMMAND, __DPKG_LIST_ARG)
	infoText, err := cmd.Output()

	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	list, err := ubuntu.parseDpkgList(infoText)

	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	return list, nil
}

func (ubuntu *Ubuntu) Get(info PackageInfo) (Package, error) {
	return Package{}, nil
}

func (ubuntu *Ubuntu) parseDpkgList(infoText []byte) ([]PackageInfo, error) {
	lines := bytes.Split(infoText, []byte("\n"))

	info := make([]PackageInfo, 0, len(lines))

	for _, l := range lines {
		fields := bytes.Fields(l)
		if len(fields) < __DPKG_FIELD_SIZE {
			// TODO return error?
			log.Print("Missing fields from dpkg line: %s", string(l))
			continue
		}

		if len(fields[0]) != len(__DPKG_IS_INSTALLED) {
			continue
		}
		for i, fieldChr := range fields[0] {
			if fieldChr != __DPKG_IS_INSTALLED[i] {
				continue
			}
		}

		name := string(fields[1])
		version := string(fields[2])
		architecture := string(fields[3])

		metadata := []MetaDataEntry{
			MetaDataEntry{
				MetaDataKey:   VERSION_KEY,
				MetaDataValue: version,
			},
		}

		dpkgInfo := PackageInfo{
			Name:     name,
			MetaData: metadata,
			System:   ubuntu.SystemName() + architecture,
		}

		info = append(info, dpkgInfo)
	}

	return info, nil
}

// TODO implement
func (ubuntu *Ubuntu) SystemName() string {
	return "ubuntu16.04"
}

const VERSION_KEY = "package_version"
const __DPKG_IS_INSTALLED = "ii"
const __DPKG_FIELD_SIZE = 4
const __DPKG_COMMAND = "dpkg"
const __DPKG_LIST_ARG = "-l"
