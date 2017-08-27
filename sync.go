package pkgthing

import (
	"log"

	"github.com/pkg/errors"
)

func AddAllPackages(lister PackageLister, getter PackageGetter, adder PackageAdder) error {
	const errMsg = "AddAllPackages failed"

	allInstalled, err := lister.GetInstalledPackages()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	// TODO concurrency
	for _, info := range allInstalled {
		pkg, err := getter.Get(info)

		if err != nil {
			log.Printf("Failed to get package for '%v': %s", info, err.Error())
			continue
		}

		_, err = adder.Add(pkg)

		if err != nil {
			log.Printf("Failed to add package '%v': %s", info, err.Error())
			continue
		}

		log.Printf("Uploaded package '%v'", info)
	}

	return nil
}
