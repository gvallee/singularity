// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package signing

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/sylabs/sif/pkg/sif"
	"github.com/sylabs/singularity/pkg/sypgp"
)

const defaultLaunch = "#!/usr/bin/env run-singularity\n"
const sifVersion = ""

func TestComputeHashStr(t *testing.T) {
	// Create temporary resources for testing
	testDir, err := ioutil.TempDir("", "hashtest-")
	if err != nil {
		t.Fatal("cannot create temporary directory")
	}
	defer os.RemoveAll(testDir)
	imagePath := testDir + "test.sif"
	var sifDescr sif.CreateInfo
	sifDescr.Pathname = imagePath
	fimg, err := sif.CreateContainer(sifDescr)
	if err != nil {
		t.Fatalf("impossible to create test SIF: %s", err)
	}

	// Invalid cases, we simply ensure that the implementation does not implode
	computeHashStr(nil, nil)

	// Valid cases

	// Test with a valid SIF but no description
	hash1 := computeHashStr(fimg, nil)
	if hash1 == "" {
		t.Fatal("invalid hash")
	}

	// Test with a valid SIF and a description structure
	descr := []*sif.Descriptor{}
	var d sif.Descriptor
	descr = append(descr, &d)

	hash2 := computeHashStr(fimg, descr)
	if hash2 == "" {
		t.Fatal("invalid hash")
	}
}

func TestSifAddSignature(t *testing.T) {
	// temporary resources
	var fingerprint [20]byte
	var signature []byte
	testDir, err := ioutil.TempDir("", "addsigntest-")
	if err != nil {
		t.Fatal("cannot create temporary directory")
	}
	defer os.RemoveAll(testDir)
	imagePath := testDir + "test.sif"
	sifDescr := sif.CreateInfo{
		Pathname:   imagePath,
		Launchstr:  defaultLaunch,
		Sifversion: sifVersion,
		ID:         uuid.NewV4(),
	}
	_, err = sif.CreateContainer(sifDescr)
	if err != nil {
		t.Fatalf("impossible to create test SIF: %s", err)
	}
	img, loadErr := sif.LoadContainer(imagePath, true)
	if loadErr != nil {
		t.Fatalf("impossible to load SIF: %s", loadErr)
	}
	defer img.UnloadContainer()

	// Invalid cases
	err = sifAddSignature(nil, 0, 0, fingerprint, signature)
	if err == nil {
		t.Fatalf("adding a signature to a SIF with invalid arguments succeeded")
	}

	// Valid cases is tested while testing Sign()
	/*
		descr, descErr := descrToSign(&img, 0, false)
		if descErr != nil {
			t.Fatalf("signing requires a primary partition: %s", descErr)
		}
		sifhash := computeHashstr(&img, descr)


		err = sifAddSignature(&img, 0, 0, fingerprint, signature)
		if err != nil {
			t.Fatalf("cannot add signature to SIF: %s", err)
		}
	*/
}

func TestSign(t *testing.T) {
	tests := []struct {
		name            string
		imgPath         string
		keyServiceURI   string
		id              uint32
		isGroup         bool
		keyIdx          int
		token           string
		successExpected bool
	}{
		{"undefined image path", "", "", 0, false, 0, "", false},
		{"invalid key index", "", "", 0, false, -1, "", false},
	}

	fmt.Println("keyring:", sypgp.SecretPath())
	if sypgp.SecretPath() == "" {
		t.Skip("No keyring, skipping test")
	}
	elist, err := sypgp.LoadPrivKeyring()
	if err != nil || elist == nil {
		t.Skip("Private key ring is not available")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Sign(test.imgPath, test.keyServiceURI, test.id, test.isGroup, test.keyIdx, test.token)
			if test.successExpected == false && err == nil {
				t.Fatalf("signing with invalid parameters succeeded")
			}
			if test.successExpected == true && err != nil {
				t.Fatalf("signing with valid parameters failed: %s", err)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	tests := []struct {
		name            string
		cpath           string
		keyServiceURI   string
		id              uint32
		isGroup         bool
		token           string
		noPrompt        bool
		expectedSuccess bool
	}{
		{"invalid image", "a/invalid/path", "", 0, false, "", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Verify(test.cpath, test.keyServiceURI, test.id, test.isGroup, test.token, test.noPrompt)
			if test.expectedSuccess == true && err != nil {
				t.Fatalf("verifying with valid parameters failed: %s", err)
			}
			if test.expectedSuccess == false && err == nil {
				t.Fatalf("verifying with invalid parameters succeeded")
			}
		})
	}
}

func TestGetSignEntities(t *testing.T) {
	tests := []struct {
		name            string
		cpath           string
		expectedSuccess bool
	}{
		{"invalid path", "not/a/valid/path", false},
		{"empty path", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetSignEntities(test.cpath)
			if test.expectedSuccess == true && err != nil {
				t.Fatalf("test with valid parameters failed: %s", err)
			}
			if test.expectedSuccess == false && err == nil {
				t.Fatalf("test with invalid parameters succeeded")
			}
		})
	}
}

func TestGetSignEntitiesFp(t *testing.T) {
	// Test with an invalid file pointer
	_, err := GetSignEntitiesFp(nil)
	if err == nil {
		t.Fatal("test with an undefined file pointer succeeded")
	}
}
