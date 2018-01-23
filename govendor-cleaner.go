package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const command = "govendor"
const searchTag = "/vendor/"

// From: https://github.com/kardianos/vendor-spec
type govendor struct {
	// Comment is free text for human use. Example "Revision abc123 introduced
	// changes that are not backwards compatible, so leave this as def876."
	Comment string `json:"comment,omitempty"`

	// Package represents a collection of vendor packages that have been copied
	// locally. Each entry represents a single Go package.
	Package []struct {
		// Import path. Example "rsc.io/pdf".
		// go get <Path> should fetch the remote package.
		Path string `json:"path"`

		// Origin is an import path where it was copied from. This import path
		// may contain "vendor" segments.
		//
		// If empty or missing origin is assumed to be the same as the Path field.
		Origin string `json:"origin"`

		// The revision of the package. This field must be persisted by all
		// tools, but not all tools will interpret this field.
		// The value of Revision should be a single value that can be used
		// to fetch the same or similar revision.
		// Examples: "abc104...438ade0", "v1.3.5"
		Revision string `json:"revision"`

		// RevisionTime is the time the revision was created. The time should be
		// parsed and written in the "time.RFC3339" format.
		RevisionTime string `json:"revisionTime"`

		// Comment is free text for human use.
		Comment string `json:"comment,omitempty"`

		// Added, not included in vendor-spec
		Version      string `json:"version,omitempty"`
		VersionExact string `json:"versionExact,omitempty"`
	} `json:"package"`
}

func removeVendor(path string) error {
	args := []string{"remove", path}
	return exec.Command(command, args...).Run()
}

func fetchVendor(path, version string) error {
	str := fmt.Sprintf("%s@%s", path, version)
	args := []string{"fetch", str}
	return exec.Command(command, args...).Run()
}

// getPkgRevisionFromVendor returns a revision or a version rule
func getPkgRevisionFromVendor(file, dep string) (string, error) {
	depVendor, err := readVendorFile(file)
	if err != nil {
		return "", err
	}

	// Find a vendoring item that is transitive
	for _, pkg := range depVendor.Package {
		if pkg.Path == dep {
			if pkg.VersionExact != "" {
				return pkg.VersionExact, nil
			}
			return pkg.Revision, nil
		}
	}
	return "", fmt.Errorf("Failed to find revision for %s in %s", dep, file)
}

func readVendorFile(path string) (*govendor, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data govendor
	json.Unmarshal(byteValue, &data)
	return &data, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument. Please provide the vendor.json file")
		return
	}

	// Get our gopath to be able to find packages
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// Get vendor.json as a parsed structure
	data, err := readVendorFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find a vendoring item that is transitive, got an origin with "vendor" in it
	for _, pkg := range data.Package {
		if pkg.Origin != "" {

			found := strings.Index(pkg.Origin, searchTag) + len(searchTag)
			if found > 0 {

				// Construct the vendor file name for base package
				// This file should contain the correct rule for the dependency
				base := pkg.Origin[:found]
				dep := pkg.Origin[found:]
				vendorFile := fmt.Sprintf("%s/src/%svendor.json", gopath, base)

				fmt.Println("Get revision for:", dep, "in:", vendorFile)
				rev, err := getPkgRevisionFromVendor(vendorFile, dep)
				if err != nil {
					fmt.Println("! Getting revision failed:", err)
				} else {
					// Run govendor in a shell
					fmt.Println("  Remove dependency:", dep)
					if err := removeVendor(dep); err != nil {
						fmt.Println("! Failed to govendor remove", dep, err)
					}

					fmt.Println("  Fetching dependency", dep, rev)
					if err := fetchVendor(dep, rev); err != nil {
						fmt.Println("! Failed to govendor fetch", dep, err)
					}
				}

			}

		}
	}
	fmt.Println("Done")
	return
}
