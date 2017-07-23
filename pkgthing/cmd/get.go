// Copyright © 2017 Johnny Morrice <john@functorama.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/johnny-morrice/pkgthing"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a package from pkgthing",
	Run: func(cmd *cobra.Command, args []string) {
		validateGetArgs()

		info := makePackageInfo()

		pkgthing := makePkgthing()
		pack, err := pkgthing.Get(info)

		if err != nil {
			die(err)
		}

		fmt.Println(pack.PackageInfo)

		writePackageFile(pack)
	},
}

func validateGetArgs() {
	ok := name != ""
	ok = ok && system != ""
	ok = ok && packageFilePath != ""

	if !ok {
		die(errors.New("Must supply name, system, and file"))
	}
}

func writePackageFile(pack pkgthing.Package) {
	file, err := os.Create(packageFilePath)

	if err != nil {
		die(err)
	}

	defer file.Close()

	_, err = file.Write(pack.Data)

	if err != nil {
		die(err)
	}
}

func makePackageInfo() pkgthing.PackageInfo {
	return pkgthing.PackageInfo{
		Name:   name,
		System: system,
	}
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVar(&name, "name", "", "Package name")
	getCmd.PersistentFlags().StringVar(&packageFilePath, "file", "", "Package file")
}
