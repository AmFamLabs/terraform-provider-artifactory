package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"session_name": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"external_id": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"policy": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "The artifactory URL",
				Required:    true,
			},
			"assume_role": assumeRoleSchema(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"artifactory_artifact_s3_deployment": resourceArtifactS3Deployment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"artifactory_artifact": dataSourceArtifact(),
		},
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			// Initial credentials loaded from SDK's default credential chain. Such as
			// the environment, shared credentials (~/.aws/credentials), or EC2 Instance
			// Role. These credentials will be used to to make the STS Assume Role API.
			sess := session.Must(session.NewSession())

			assumeRoleList := d.Get("assume_role").(*schema.Set).List()

			if len(assumeRoleList) == 1 {
				assumeRole := assumeRoleList[0].(map[string]interface{})
				// Create the credentials from AssumeRoleProvider to assume the role
				creds := stscreds.NewCredentials(sess, assumeRole["role_arn"].(string))
				new_sess, err := session.NewSessionWithOptions(session.Options{
					Config: aws.Config{Credentials: creds},
				})
				return new_ess, err

				// alternative...?
				//assumeRoleInput := &sts.AssumeRoleInput{
				//	//DurationSeconds:   aws.Int64(int64((p.Duration - jitter) / time.Second)),
				//	RoleArn:         aws.String(assumeRole["role_arn"].(string)),
				//	RoleSessionName: aws.String(assumeRole["session_name"].(string)),
				//	// ExternalId:
				//	// Tags:              p.Tags,
				//	// PolicyArns:        p.PolicyArns,
				//	// TransitiveTagKeys: p.TransitiveTagKeys,
				//}
				//sts_svc := sts.New(sess)

				//result, err := sts_svc.AssumeRole(assumeRoleInput)
				//if err != nil {
				//	return nil, err
				//}

				//if v := assumeRole["policy"].(string); v != "" {
				//	config.AssumeRolePolicy = v
				//}

				//log.Printf("[INFO] assume_role configuration set: (ARN: %q, SessionID: %q, ExternalID: %q, Policy: %q)", config.AssumeRoleARN, config.AssumeRoleSessionName, config.AssumeRoleExternalID, config.AssumeRolePolicy)
			} else {
				log.Printf("[INFO] No assume_role block read from configuration")
				return sess, nil
			}

		},
	}
}
