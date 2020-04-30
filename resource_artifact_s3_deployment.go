package main

import (
	"bytes"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
)

func resourceArtifactS3Deployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArtifactS3DeploymentCreate,
		Read:   resourceArtifactS3DeploymentRead,
		Update: resourceArtifactS3DeploymentUpdate,
		Delete: resourceArtifactS3DeploymentDelete,

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

func resourceArtifactS3DeploymentCreate(d *schema.ResourceData, m interface{}) error {
	repository_path := d.Get("repository_path").(string)
	d.SetId(repository_path)
	return resourceArtifactS3DeploymentRead(d, m)
}

func resourceArtifactS3DeploymentRead(d *schema.ResourceData, m interface{}) error {
	return dataSourceArtifactRead(d, m)
}

func resourceArtifactS3DeploymentUpdate(d *schema.ResourceData, m interface{}) error {

	return resourceArtifactS3DeploymentRead(d, m)
}

func resourceArtifactS3DeploymentDelete(d *schema.ResourceData, m interface{}) error {
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
