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
		Delete: resourceArtifactS3DeploymentDelete,

		Schema: map[string]*schema.Schema{
			"repository_path": {
				Type:        schema.TypeString,
				Description: "The source artifactory repository_path.",
				Required:    true,
				ForceNew:    true,
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
			"s3_region": {
				Type:        schema.TypeString,
				Description: "The target S3 region to override the default.",
				Optional:    true,
				ForceNew:    true,
			},
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
			"s3_key": {
				Type:        schema.TypeString,
				Description: "This is computed from the bucket_prefix and repository_path.",
				Computed:    true,
			},
			"s3_region_actual": {
				Type:        schema.TypeString,
				Description: "The actual S3 region, whether set explictly or inherited.",
				Computed:    true,
			},
		},
	}
}

func doesS3ObjectExist(_s3 *s3.S3, s3_bucket string, s3_key string) (bool, error) {
	input := &s3.HeadObjectInput{Bucket: aws.String(s3_bucket), Key: aws.String(s3_key)}

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

func getS3Region(d *schema.ResourceData, sess *session.Session) string {
	// get it if it's set on this resource
	s3_region := d.Get("s3_region").(string)
	if s3_region == "" {
		// use whatever is set on provider
		s3_region = *sess.Config.Region
	}
	return s3_region
}

func resourceArtifactS3DeploymentCreate(d *schema.ResourceData, m interface{}) error {
	// get configuration values
	repository_path := d.Get("repository_path").(string)
	s3_bucket := d.Get("s3_bucket").(string)
	s3_prefix := d.Get("s3_prefix").(string)

	// create session and ensure region is set
	sess := m.(*session.Session)
	s3_region := getS3Region(d, sess)
	svc := s3.New(sess, &aws.Config{Region: aws.String(s3_region)})

	// get info about the file from artifactory
	var f FileInfo
	err := getFileInfo(repository_path, &f)
	if err != nil {
		// if can't retrieve info about file from artifactory, creation fails
		d.SetId("")
		return err
	}

	// collect the params needed for file upload
	filename := filepath.Base(f.Path)
	s3_key := filepath.Join(s3_prefix, filename)

	// check if file already exists in s3
	it_does_exist, err := doesS3ObjectExist(svc, s3_bucket, s3_key)
	if err != nil {
		return err
	}

	if it_does_exist {
		// if it already exists creation fails
		d.SetId("")
		return fmt.Errorf(
			"'s3://%s' already exists, remove the object on s3 manually and re-apply.",
			filepath.Join(s3_bucket, s3_key),
		)
	}

	// set computed resource params
	download_uri := f.DownloadUri
	d.Set("checksums", checksumsToMap(f.Checksums))
	d.Set("download_uri", download_uri)
	d.Set("original_checksums", checksumsToMap(f.OriginalChecksums))
	d.Set("path", f.Path)
	d.Set("repo", f.Repo)
	d.Set("size", f.Size)

	d.Set("s3_key", s3_key)
	d.Set("s3_region_actual", s3_region)

	// download the file from artifactory and upload it to s3
	uploader := s3manager.NewUploaderWithClient(svc)
	artifact_binary_resp, err := http.Get(download_uri)
	if err != nil {
		return err
	}

	defer artifact_binary_resp.Body.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3_bucket),
		Key:    aws.String(s3_key),
		Body:   artifact_binary_resp.Body,
	})
	if err != nil {
		return err
	}

	return resourceArtifactS3DeploymentRead(d, m)
}

func resourceArtifactS3DeploymentRead(d *schema.ResourceData, m interface{}) error {
	repository_path := d.Get("repository_path").(string)
	s3_bucket := d.Get("s3_bucket").(string)
	s3_key := d.Get("s3_key").(string)
	s3_region := d.Get("s3_region_actual").(string)

	sess := m.(*session.Session)
	svc := s3.New(sess, &aws.Config{Region: aws.String(s3_region)})
	it_does_exist, err := doesS3ObjectExist(svc, s3_bucket, s3_key)
	if err != nil {
		return err
	}

	if it_does_exist {
		d.SetId(fmt.Sprintf("%s__s3://%s", repository_path, filepath.Join(s3_bucket, s3_key)))
	} else {
		d.SetId("")
	}

	return nil
}

func resourceArtifactS3DeploymentDelete(d *schema.ResourceData, m interface{}) error {
	s3_bucket := d.Get("s3_bucket").(string)
	s3_key := d.Get("s3_key").(string)
	s3_region := d.Get("s3_region_actual").(string)

	sess := m.(*session.Session)
	svc := s3.New(sess, &aws.Config{Region: aws.String(s3_region)})

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
