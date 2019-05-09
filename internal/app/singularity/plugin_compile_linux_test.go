// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package singularity

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sylabs/singularity/internal/pkg/test"
)

func TestGetSingularityDirFromBasedir(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	emptyDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %s", err)
	}

	curDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %s", err)
	}
	singularitySrcDir := filepath.Join(curDir, "..", "..", "..")

	tests := []struct {
		name      string
		path      string
		shallPass bool
	}{
		{
			name:      "empty path",
			path:      "",
			shallPass: false,
		},
		{
			name:      "invalid path",
			path:      "/not/a/valid/path",
			shallPass: false,
		},
		{
			name:      "empty directory",
			path:      emptyDir,
			shallPass: false,
		},
		{
			name:      "valid Singularity top source directory",
			path:      singularitySrcDir,
			shallPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getSingularitySrcDirFromBasedir(tt.path)
			if tt.shallPass == false && err == nil {
				t.Fatalf("test %s was expected to fail but succeeded", tt.name)
			}
			if tt.shallPass == true && err != nil {
				t.Fatalf("test %s was expected to succeed but failed", tt.name)
			}
		})
	}
}

func TestCompilePlugin(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	dummyPluginDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temporary file: %s", err)
	}
	defer os.Remove(dummyPluginDir)

	tests := []struct {
		name      string
		srcDir    string
		destDir   string
		tags      string
		shallPass bool
	}{
		{
			name:      "empty srcDir; empty destDir; empty tags",
			srcDir:    "",
			destDir:   "",
			tags:      "",
			shallPass: false,
		},
		{
			name:      "dummy srcDir; empty destDir; empty tags",
			srcDir:    dummyPluginDir,
			destDir:   "",
			shallPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CompilePlugin(tt.srcDir, tt.destDir, tt.tags)
			if tt.shallPass == false && err == nil {
				t.Fatalf("invalid case (%s) succeeded", tt.name)
			}
			if tt.shallPass == true && err != nil {
				t.Fatalf("valid case (%s) failed: %s", tt.name, err)
			}
		})
	}
}
