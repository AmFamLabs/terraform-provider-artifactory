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
			"checksums": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArtifactCreate(d *schema.ResourceData, m interface{}) error {
	return resourceArtifactRead(d, m)
}

func resourceArtifactRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceArtifactUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceArtifactRead(d, m)
}

func resourceArtifactDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
