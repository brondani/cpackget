/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package utils_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	errs "github.com/open-cmsis-pack/cpackget/cmd/errors"
	"github.com/open-cmsis-pack/cpackget/cmd/utils"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name     string
	path     string
	expected utils.PackInfo
	err      error
	short    bool
}

func absPath(filePath string) string {
	abs, _ := filepath.Abs(filePath)
	return abs
}

func TestExtractPackInfo(t *testing.T) {
	assert := assert.New(t)
	cwd, _ := os.Getwd()

	var tests = []testCase{
		// Cases where short=true
		{
			name:  "test short path bad pack name",
			path:  "this-is-not-a-valid-pack-name",
			short: true,
			err:   errs.ErrBadPackName,
		},
		{
			name:  "test short path bad pack name with invalid version",
			path:  "TheVendor.ThePack.not-a-valid-version",
			short: true,
			err:   errs.ErrBadPackNameInvalidVersion,
		},
		{
			name:  "test short path bad pack name with invalid vendor name",
			path:  "not-a-valid-vendor?.ThePack",
			short: true,
			err:   errs.ErrBadPackNameInvalidName,
		},
		{
			name:  "test short path bad pack name with invalid pack name",
			path:  "TheVendor.not-a-valid-pack-name?",
			short: true,
			err:   errs.ErrBadPackNameInvalidName,
		},
		{
			name:  "test short path successfully extract pack info",
			path:  "TheVendor.ThePack.0.0.1",
			short: true,
			expected: utils.PackInfo{
				Vendor:  "TheVendor",
				Pack:    "ThePack",
				Version: "0.0.1",
			},
		},
		{
			name:  "test short path successfully extract pack info without version",
			path:  "TheVendor.ThePack",
			short: true,
			expected: utils.PackInfo{
				Vendor: "TheVendor",
				Pack:   "ThePack",
			},
		},

		// Tests with full paths (short=false)
		{
			name: "test path with bad extension",
			path: "only-check-for-extension.txt",
			err:  errs.ErrBadPackNameInvalidExtension,
		},
		{
			name: "test path with bad name pack extension",
			path: "not-a-valid-pack-name.pack",
			err:  errs.ErrBadPackName,
		},
		{
			name: "test path with bad name zip extension",
			path: "not-a-valid-pack-name.zip",
			err:  errs.ErrBadPackName,
		},
		{
			name: "test path with bad name pdsc extension",
			path: "not-a-valid-pack-name.pdsc",
			err:  errs.ErrBadPackName,
		},
		{
			name: "test pdsc path with bad vendor name",
			path: "not-a-valid-vendor-name?.ThePack.pdsc",
			err:  errs.ErrBadPackNameInvalidVendor,
		},
		{
			name: "test pack path with bad vendor name",
			path: "not-a-valid-vendor-name?.ThePack.0.0.1.pack",
			err:  errs.ErrBadPackNameInvalidVendor,
		},
		{
			name: "test zip path with bad vendor name",
			path: "not-a-valid-vendor-name?.ThePack.0.0.1.pdsc",
			err:  errs.ErrBadPackNameInvalidVendor,
		},
		{
			name: "test pdsc path with bad pack name",
			path: "TheVendor.not-a-valid-pack-name?.pdsc",
			err:  errs.ErrBadPackNameInvalidName,
		},
		{
			name: "test pack path with bad pack name",
			path: "TheVendor.not-a-valid-pack-name?.0.0.1.pack",
			err:  errs.ErrBadPackNameInvalidName,
		},
		{
			name: "test zip path with bad pack name",
			path: "TheVendor.not-a-valid-pack-name?.0.0.1.zip",
			err:  errs.ErrBadPackNameInvalidName,
		},
		{
			name: "test pack path with bad version",
			path: "TheVendor.ThePack.not-a-valid-version?.pack",
			err:  errs.ErrBadPackNameInvalidVersion,
		},
		{
			name: "test zip path with bad version",
			path: "TheVendor.ThePack.not-a-valid-version?.zip",
			err:  errs.ErrBadPackNameInvalidVersion,
		},
		{
			name: "test path with with http URL",
			path: "http://vendor.com/TheVendor.ThePack.0.0.1.pack",
			expected: utils.PackInfo{
				Vendor:    "TheVendor",
				Pack:      "ThePack",
				Version:   "0.0.1",
				Extension: ".pack",
				Location:  "http://vendor.com/",
			},
		},
		{
			name: "test path with with relative path",
			path: "relative/path/to/TheVendor.ThePack.0.0.1.pack",
			expected: utils.PackInfo{
				Vendor:    "TheVendor",
				Pack:      "ThePack",
				Version:   "0.0.1",
				Extension: ".pack",
				Location:  "file://" + path.Join(cwd, "relative/path/to") + "/",
			},
		},
		{
			name: "test path with with relative path and dot-dot",
			path: "../path/to/TheVendor.ThePack.0.0.1.pack",
			expected: utils.PackInfo{
				Vendor:    "TheVendor",
				Pack:      "ThePack",
				Version:   "0.0.1",
				Extension: ".pack",
				Location:  "file://" + absPath(path.Join(cwd, "../path/to")) + "/",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info, err := utils.ExtractPackInfo(test.path, test.short)
			if test.err != nil {
				assert.True(errs.Is(err, test.err))
			} else {
				assert.Equal(info, test.expected)
			}
		})
	}
}