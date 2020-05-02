package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"path/filepath"
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
				ForceNew: true,
			},
			"s3_key": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"s3_prefix": {
				Type:     schema.TypeString,
				Optional: true,
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
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"original_checksums": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"download_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArtifactS3DeploymentCreate(d *schema.ResourceData, m interface{}) error {
	resourceArtifactS3DeploymentRead(d, m)
	sess := m.(*session.Session)
	uploader := s3manager.NewUploader(sess)
	artifact_binary_resp, err := http.Get(d.Get("download_uri").(string))
	if err != nil {
		return err
	}
	defer artifact_binary_resp.Body.Close()
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(d.Get("s3_bucket").(string)),
		Key:    aws.String(d.Get("s3_key").(string)),
		Body:   artifact_binary_resp.Body,
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceArtifactS3DeploymentRead(d *schema.ResourceData, m interface{}) error {
	repository_path := d.Get("repository_path").(string)
	s3_bucket := d.Get("s3_bucket").(string)
	s3_prefix := d.Get("s3_prefix").(string)

	var f FileInfo
	err := getFileInfo(repository_path, &f)
	if err != nil {
		return err
	}

	filename := filepath.Base(f.Path)
	s3_key := fmt.Sprintf("%s/%s", s3_prefix, filename)
	d.SetId(fmt.Sprintf("%s__s3://%s/%s", repository_path, s3_bucket, s3_key))
	d.Set("s3_key", s3_key)
	d.Set("checksums", checksumsToMap(f.Checksums))
	d.Set("download_uri", f.DownloadUri)
	d.Set("original_checksums", checksumsToMap(f.OriginalChecksums))
	d.Set("path", f.Path)
	d.Set("repo", f.Repo)
	d.Set("size", f.Size)

	return nil
}

func resourceArtifactS3DeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceArtifactS3DeploymentRead(d, m)
}

func resourceArtifactS3DeploymentDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*session.Session)
	svc := s3.New(sess)
	s3_bucket := d.Get("s3_bucket").(string)
	s3_key := d.Get("s3_key").(string)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3_bucket),
		Key:    aws.String(s3_key),
	}
	_, err := svc.DeleteObject(input)
	if err != nil {
		return err
	}
	return nil
}
