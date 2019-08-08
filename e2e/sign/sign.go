// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package sign

import (
	"testing"

	"github.com/sylabs/singularity/e2e/internal/e2e"
)

type ctx struct {
	env e2e.TestEnv
	imgPath string
}

const imgURL = "library://alpine:latest"

func (c *ctx) singularitySignHelpOption(t *testing.T) {
	c.env.RunSingularity(
		t,
		e2e.WithPrivileges(false),
		e2e.WithCommand("sign"),
		e2e.WithArgs("--help"),
		e2e.ExpectExit(
			0,
			e2e.ExpectOutput(e2e.ContainMatch, "Attach a cryptographic signature to an image"),
		),
	)
}


func (c *ctx) singularitySignIDOption(t *testing.T) {
	cmdArgs := []string{"--id", "0", c.successImage}
	c.env.RunSingularity(
		t,
		e2e.WithPrivileges(false),
		e2e.WithCommand("sign"),
		e2e.WithArgs(cmdArgs...),
		e2e.ExpectExit(
			0,
			e2e.ExpectOutput(e2e.ContainMatch, "Container is signed by 1 key(s):"),
		),
	)
}

// RunE2ETests is the main func to trigger the test suite
func RunE2ETests(env e2e.TestEnv) func(*testing.T) {
	c := &ctx{
		env: env,
		imgPath: filepath.Join(env.TestDir, "testImage.sif")
	}

	e2e.PullImage(t, c.env, imageURL, tt.imagePath)

	return func(t *testing.T) {
		t.Run("singularitySignHelpOption", c.singularity)
		t.Run("singularitySignIDOption", c.singularitySignIDOption)
	}
}
