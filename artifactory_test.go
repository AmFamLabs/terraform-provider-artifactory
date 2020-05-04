package main

import (
	"reflect"
	"testing"
)

func TestChecksumsToMap(t *testing.T) {
	base := Checksums{
		Md5:    "18cfe70b69a4f06b5a0244689db39775",
		Sha1:   "fa0b1073315c85861fd8221c0af802b3b2da1860",
		Sha256: "5920550b1cbf63e8887e13d2a0d8c71c8df69770a55bb14511c0276592c046c9",
	}
	cases := []struct {
		input    Checksums
		expected map[string]string
	}{
		{
			input: base,
			expected: map[string]string{
				"md5":    base.Md5,
				"sha1":   base.Sha1,
				"sha256": base.Sha256,
			},
		},
		{
			input: Checksums{
				Md5:    "this",
				Sha1:   "would",
				Sha256: "work",
			},
			expected: map[string]string{
				"md5":    "this",
				"sha1":   "would",
				"sha256": "work",
			},
		},
	}
	for _, c := range cases {
		out := checksumsToMap(c.input)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching input with output. %#v vs. %#v", out, c.expected)
		}
	}
}

func TestGetFileInfo_ArtifactExists(t *testing.T) {
	var f FileInfo
	repository_path_that_exists := "lambda/propinc/ingest/replicate-2.30.0.zip"
	err := getFileInfo(repository_path_that_exists, &f)
	if err != nil {
		t.Fatalf("Error with getFileInfo %#v", err)
	}
}

func TestGetFileInfo_ArtifactDoesntExist(t *testing.T) {
	var f FileInfo
	err := getFileInfo("i/dont/exist", &f)
	if err == nil {
		t.Fatalf("Error with getFileInfo when trying to fail with a non-existant artifact. Err was %#v", err)
	}
}
