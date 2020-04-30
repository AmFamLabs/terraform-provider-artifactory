package main

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"fmt"
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
			"s3_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"s3_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// these can be available
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
	svc := m.(*s3.S3)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		return err
	}
	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

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
