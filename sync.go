package pkgthing

import (
	"log"
	"sync"

	"github.com/pkg/errors"
)

type Syncer struct {
	Lister            PackageLister
	Getter            PackageGetter
	Adder             PackageAdder
	GetterConcurrency int
	AdderConcurrency  int
}

func (syncer Syncer) AddAllPackages() error {
	const errMsg = "AddAllPackages failed"

	if syncer.GetterConcurrency == 0 {
		syncer.GetterConcurrency = __DEFAULT_GETTER_CONCURRENCY
	}

	if syncer.AdderConcurrency == 0 {
		syncer.AdderConcurrency = __DEFAULT_ADDER_CONCURRENCY
	}

	allInstalled, err := syncer.Lister.GetInstalledPackages()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	wg := &sync.WaitGroup{}
	getSem := makeSem(syncer.GetterConcurrency)
	addSem := makeSem(syncer.AdderConcurrency)
	for _, i := range allInstalled {
		info := i
		wg.Add(1)
		go func() {
			lockSem(getSem)
			defer unlockSem(getSem)
			defer wg.Done()

			pkg, err := syncer.Getter.Get(info)

			if err != nil {
				log.Printf("Failed to get package for '%v': %s", info, err.Error())
				return
			}

			wg.Add(1)
			go func() {
				lockSem(addSem)
				defer unlockSem(addSem)
				defer wg.Done()
				_, err = syncer.Adder.Add(pkg)

				if err != nil {
					log.Printf("Failed to add package '%v': %s", info, err.Error())
					return
				}

				log.Printf("Synced package '%v'", info)
			}()
		}()
	}

	wg.Wait()

	return nil
}

func makeSem(c int) chan struct{} {
	return make(chan struct{}, c)
}

func lockSem(sem chan struct{}) {
	sem <- struct{}{}
}

func unlockSem(sem chan struct{}) {
	<-sem
}

const __DEFAULT_GETTER_CONCURRENCY = 10
const __DEFAULT_ADDER_CONCURRENCY = 10
