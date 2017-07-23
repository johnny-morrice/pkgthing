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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/johnny-morrice/pkgthing"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search pkgthing for packages",

	Run: func(cmd *cobra.Command, args []string) {

	},
}

var searchTerm string
var searchKeyText string

func validateSearchArgs() {
	if searchTerm == "" || searchKeyText == "" {
		die(errors.New("Must specify term and field"))
	}
}

func makeSearchTerm(searchKey pkgthing.SearchKey) pkgthing.PackageSearchTerm {
	return pkgthing.PackageSearchTerm{
		System:     system,
		SearchTerm: searchTerm,
		SearchKey:  searchKey,
	}
}

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().StringVar(&searchTerm, "term", "", "Search term")
	searchCmd.PersistentFlags().StringVar(&searchKeyText, "field", "name", "Search field")
}
