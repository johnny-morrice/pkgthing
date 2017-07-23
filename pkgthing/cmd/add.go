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
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/johnny-morrice/pkgthing"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a package to pkgthing",
	Run: func(cmd *cobra.Command, args []string) {
		validateAddArgs()

		file := openPackageFile()
		defer file.Close()

		pack := makeNewPackage(file)

		pkgthing := makePkgthing()
		pkgthing.Add(pack)
	},
}

var packageFilePath string

func validateAddArgs() {
	ok := name != ""
	ok = ok && system != ""
	ok = ok && packageFilePath != ""

	if !ok {
		die(errors.New("Must supply name, system, and file"))
	}
}

func makeNewPackage(r io.Reader) pkgthing.Package {
	data, err := ioutil.ReadAll(r)

	if err != nil {
		die(err)
	}

	pack := pkgthing.Package{}
	pack.Name = name
	pack.System = system
	pack.Data = data
	return pack
}

func openPackageFile() io.ReadCloser {
	file, err := os.Open(packageFilePath)

	if err != nil {
		die(err)
	}

	return file
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVar(&name, "name", "", "Package name")
	addCmd.PersistentFlags().StringVar(&packageFilePath, "file", "", "Package file")
}
