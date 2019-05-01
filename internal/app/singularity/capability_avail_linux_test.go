// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package singularity

import (
	"testing"

	"github.com/sylabs/singularity/internal/pkg/test"
)

func TestCapabilityAvail(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	var emptyCap CapAvailConfig
	invalidCap := CapAvailConfig{
		Caps: "justARandomString,anotherRandomString",
		Desc: true,
	}

	validCap := CapAvailConfig{
		Caps: "CAP_ALL", // CAP_ALL should always be available
		Desc: false,
	}

	tests := []struct {
		name      string
		c         CapAvailConfig
		shallPass bool
	}{
		{
			name:      "empty capability",
			c:         emptyCap,
			shallPass: true,
		},
		{
			name:      "invalid capability",
			c:         invalidCap,
			shallPass: false,
		},
		{
			name:      "valid capabilities",
			c:         validCap,
			shallPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CapabilityAvail(tt.c)
			if tt.shallPass == true && err != nil {
				t.Fatalf("valid case failed: %s\n", err)
			}
			if tt.shallPass == false && err == nil {
				t.Fatal("invalid case succeeded")
			}
		})
	}
}
