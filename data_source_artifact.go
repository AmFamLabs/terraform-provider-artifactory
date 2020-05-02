package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			// everything else is computed
			//
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
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceArtifactRead(d *schema.ResourceData, m interface{}) error {
	repository_path := d.Get("repository_path").(string)
	d.SetId(repository_path)
	var f FileInfo
	err := getFileInfo(repository_path, &f)
	if err != nil {
		return err
	}
	if f.Path == "" {
		d.SetId("")
		// we dont continue
		return nil
	}
	d.Set("checksums", checksumsToMap(f.Checksums))
	d.Set("download_uri", f.DownloadUri)
	d.Set("original_checksums", checksumsToMap(f.OriginalChecksums))
	d.Set("path", f.Path)
	d.Set("repo", f.Repo)
	d.Set("size", f.Size)
	return nil
}
