package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func dataSourceArtifact() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArtifactRead,

		Schema: map[string]*schema.Schema{
			"repository_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			//
			//"body": {
			//	//
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"checksums": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"download_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"original_checksums": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func checksumsToMap(c Checksums) map[string]string {
	return map[string]string{
		"md5":    c.Md5,
		"sha1":   c.Sha1,
		"sha256": c.Sha256,
	}
}

func grabBody(uri string) error {
	artifact_binary_resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(artifact_binary_resp.Body)
	if err != nil {
		return err
	}
	defer artifact_binary_resp.Body.Close()
	return nil
}

func dataSourceArtifactRead(d *schema.ResourceData, m interface{}) error {
	repository_path := d.Get("repository_path").(string)
	d.SetId(repository_path)
	api_url := fmt.Sprintf("https://artifacts.amfamlabs.com/api/storage/%s", repository_path)
	resp, err := http.Get(api_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// TODO: better logic if 404-not found, but we need to delete this with setId
		d.SetId("")
		// we don't need to continue
		return nil
	}
	var f FileInfo
	jsonerr := json.NewDecoder(resp.Body).Decode(&f)
	if jsonerr != nil {
		return jsonerr
	}
	d.Set("checksums", checksumsToMap(f.Checksums))
	d.Set("download_uri", f.DownloadUri)
	d.Set("original_checksums", checksumsToMap(f.OriginalChecksums))
	d.Set("path", f.Path)
	d.Set("repo", f.Repo)
	d.Set("size", f.Size)
	return nil
}
