# terraform-provider-artifactory

This provider's goal is given *some* input, an artifact from JFrog's
Artifactory will create a S3 object deployment that's state can be
traced/tracked for influencing a vanilla lambda deployment that wants an s3
object for deployment.

## Feature Wishlist

- *MVP?* a `resource` and `data`
 - [ ] a `resource` that is a `artifactory_artifact_s3_deployment` which will produce a
   resource that will have an attributes of `s3_bucket` and `s3_key` that can be
   used to feed a vanilla
   [`aws_lambda_function`](https://www.terraform.io/docs/providers/aws/r/lambda_function.html) defining a resource with
   `s3_bucket` and `s3_key`. 
   - Changes to the source `artifactory_artifact_s3_deployment` resource
     will hash if Artifactory deployment isn't already doing this for us. (I
     think it is)
 - [X] a `data` that is a `artiactory_artifact` which will utilize Artifactory's
   REST API for grabbing information of a given `repository_path`



## Development

### Prerequisites

- `go` installed

### Make plugin

```sh
go build
```

### Run plugin

```sh
terraform init
terraform plan && terraform apply
```
