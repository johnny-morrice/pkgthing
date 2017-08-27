package pkgthing

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type Ubuntu struct {
	tempDir string
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
	const errMsg = "Ubuntu.Get failed"

	ubuntu.chdirTemp()

	cmd := exec.Command(__FAKEROOT_COMMAND, __FAKEROOT_ARG, __REPACK_COMMAND, info.Name)
	err := cmd.Run()

	if err != nil {
		return Package{}, errors.Wrap(err, errMsg)
	}

	debFileName := ubuntu.debFileName(info)
	contents, err := ioutil.ReadFile(debFileName)

	if err != nil {
		return Package{}, errors.Wrap(err, errMsg)
	}

	pkg := Package{
		PackageInfo: info,
		Data:        contents,
	}

	return pkg, nil
}

func (ubuntu *Ubuntu) parseDpkgList(infoText []byte) ([]PackageInfo, error) {
	lines := bytes.Split(infoText, []byte("\n"))

	info := make([]PackageInfo, 0, len(lines))

	for _, l := range lines {
		fields := bytes.Fields(l)
		if len(fields) < __DPKG_FIELD_SIZE {
			// TODO return error?
			log.Printf("Missing fields from dpkg line: %s", string(l))
			continue
		}

		if len(fields[__DPKG_STATUS_FIELD]) != len(__DPKG_IS_INSTALLED) {
			continue
		}
		for i, fieldChr := range fields[__DPKG_STATUS_FIELD] {
			if fieldChr != __DPKG_IS_INSTALLED[i] {
				continue
			}
		}

		name := string(fields[__DPKG_NAME_FIELD])
		version := string(fields[__DPKG_VERSION_FIELD])
		architecture := string(fields[__DPKG_ARCH_FIELD])

		metadata := []MetaDataEntry{
			MetaDataEntry{
				MetaDataKey:   VERSION_KEY,
				MetaDataValue: version,
			},
			MetaDataEntry{
				MetaDataKey:   ARCHITECTURE_KEY,
				MetaDataValue: architecture,
			},
		}

		dpkgInfo := PackageInfo{
			Name:     name,
			MetaData: metadata,
			System:   ubuntu.SystemName(),
		}

		info = append(info, dpkgInfo)
	}

	return info, nil
}

func (ubuntu *Ubuntu) debFileName(info PackageInfo) string {
	parts := []string{
		info.Name,
		info.GetMetaData(VERSION_KEY),
		info.GetMetaData(ARCHITECTURE_KEY),
	}
	return strings.Join(parts, "_") + ".deb"
}

// TODO implement
func (ubuntu *Ubuntu) SystemName() string {
	return "ubuntu16.04"
}

func (ubuntu *Ubuntu) chdirTemp() {
	if ubuntu.tempDir != "" {
		return
	}

	dirname, err := ioutil.TempDir(__TEMP_ROOT, __TEMP_PREFIX)

	if err == nil {
		err := os.Chdir(dirname)

		if err != nil {
			log.Print("Failed to change directory prior to repacking deb file")
			return
		}
	} else {
		log.Print("Failed to create temporary directory prior to repacking deb file")
		return
	}

	ubuntu.tempDir = dirname
}

const VERSION_KEY = "version"
const ARCHITECTURE_KEY = "architecture"
const __DPKG_STATUS_FIELD = 0
const __DPKG_NAME_FIELD = 1
const __DPKG_VERSION_FIELD = 2
const __DPKG_ARCH_FIELD = 3
const __FAKEROOT_COMMAND = "fakeroot"
const __FAKEROOT_ARG = "-u"
const __REPACK_COMMAND = "dpkg-repack"
const __TEMP_ROOT = "/tmp"
const __TEMP_PREFIX = "pkgthing"
const __DPKG_IS_INSTALLED = "ii"
const __DPKG_FIELD_SIZE = 4
const __DPKG_COMMAND = "dpkg"
const __DPKG_LIST_ARG = "-l"
