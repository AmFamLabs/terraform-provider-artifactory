package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Checksums struct {
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
}

type FileInfo struct {
	Checksums         Checksums `json:"checksums"`
	DownloadUri       string    `json:"downloadUri"`
	OriginalChecksums Checksums `json:"originalChecksums"`
	Path              string    `json:"path"`
	Repo              string    `json:"repo"`
	Size              string    `json:"size"`
	Uri               string    `json:"uri"`
}

func getFileInfo(repository_path string, f *FileInfo) error {

	resp, err := http.Get(fmt.Sprintf("https://artifacts.amfamlabs.com/api/storage/%s", repository_path))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("We didn't find that artifact :(")
	}
	jsonerr := json.NewDecoder(resp.Body).Decode(&f)
	if jsonerr != nil {
		return jsonerr
	}
	return nil
}

func checksumsToMap(c Checksums) map[string]string {
	return map[string]string{
		"md5":    c.Md5,
		"sha1":   c.Sha1,
		"sha256": c.Sha256,
	}
}
