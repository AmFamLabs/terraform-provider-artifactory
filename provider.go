package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The artifactory URL",
				Required:    true,
			},
		},
		//ResourcesMap: map[string]*schema.Resource{
		//	"artifactory_artifact": resourceArtifact(),
		//},
		DataSourcesMap: map[string]*schema.Resource{
			"artifactory_artifact": dataSourceArtifact(),
		},
	}
}
