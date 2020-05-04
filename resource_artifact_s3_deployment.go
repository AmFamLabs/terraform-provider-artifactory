package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
				Type:        schema.TypeString,
				Description: "The source artifactory repository_path.",
				Required:    true,
			},
			"s3_bucket": {
				Type:        schema.TypeString,
				Description: "The target s3 bucket name.",
				Required:    true,
				ForceNew:    true,
			},
			"s3_prefix": {
				Type:        schema.TypeString,
				Description: "The target s3 prefix for the target s3 bucket.",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			"s3_key": {
				Type:        schema.TypeString,
				Description: "This is computed from the bucket_prefix and repository_artifact.",
				Computed:    true,
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

func doesS3ObjectExist(input *s3.HeadObjectInput, _s3 *s3.S3) (bool, error) {
	_, err := _s3.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.RequestFailure); ok {
			if aerr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, err
	} else {
		return true, nil
	}
}

func resourceArtifactS3DeploymentCreate(d *schema.ResourceData, m interface{}) error {
	// set the resource properies
	// this is a little weird that it's called here, but its a first pass and
	// the 'deployment' aspect is what we're playing with
	// so you can think of READ as gathering metadta from artifactory and
	// CREATE as writing to S3
	//
	// alternatively this could be implemented requiring a
	// data.artifactory_artifact as most of the input
	//
	// what we're playing with though is the ux of this provider when using in
	// terraform for our use case :)
	err := resourceArtifactS3DeploymentRead(d, m)
	s3_bucket := d.Get("s3_bucket").(string)
	s3_key := d.Get("s3_key").(string)
	if err != nil {
		return err
	}
	// s3 read.. CREATE logic
	sess := m.(*session.Session)
	svc := s3.New(sess)
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s3_bucket),
		Key:    aws.String(s3_key),
	}
	it_does_exist, err := doesS3ObjectExist(input, svc)
	if it_does_exist && err == nil {
		// explicitly say we didn't make this resource
		d.SetId("")
		return fmt.Errorf("'s3://%s' already exists, remove the object on s3 manually and re-apply.", filepath.Join(s3_bucket, s3_key))
	}
	if err != nil {
		return err
	}
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

func setS3Key(d *schema.ResourceData, path string, prefix string) error {
	filename := filepath.Base(path)
	s3_key := filepath.Join(prefix, filename)
	d.Set("s3_key", s3_key)
	return nil
}

func resourceArtifactS3DeploymentRead(d *schema.ResourceData, m interface{}) error {
	// TODO add s3 read
	repository_path := d.Get("repository_path").(string)
	s3_bucket := d.Get("s3_bucket").(string)
	s3_prefix := d.Get("s3_prefix").(string)

	var f FileInfo
	err := getFileInfo(repository_path, &f)
	if err != nil {
		return err
	}

	filename := filepath.Base(f.Path)
	s3_key := filepath.Join(s3_prefix, filename)
	d.SetId(fmt.Sprintf("%s__s3://%s", repository_path, filepath.Join(s3_bucket, s3_key)))
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
	// TODO recreate when keys change over "ForceRecreate" option in Schema
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
