package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceArtifact() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArtifactRead,

		Schema: map[string]*schema.Schema{
			"repository_path": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceArtifactRead(d *schema.ResourceData, m interface{}) error {
	return resourceArtifactRead(d, m)
}
