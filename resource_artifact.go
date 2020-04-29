package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceArtifact() *schema.Resource {
	return &schema.Resource{
		Create: resourceArtifactCreate,
		Read:   resourceArtifactRead,
		Update: resourceArtifactUpdate,
		Delete: resourceArtifactDelete,

		Schema: map[string]*schema.Schema{
			"repository_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksums": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArtifactCreate(d *schema.ResourceData, m interface{}) error {
	repository_path := d.Get("repository_path").(string)
	d.SetId(repository_path)
	return resourceArtifactRead(d, m)
}

func resourceArtifactRead(d *schema.ResourceData, m interface{}) error {
	return dataSourceArtifactRead(d, m)
}

func resourceArtifactUpdate(d *schema.ResourceData, m interface{}) error {

	return resourceArtifactRead(d, m)
}

func resourceArtifactDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

// TODO: useful as starting point for writing to s3
func bufBody(uri string) error {
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
