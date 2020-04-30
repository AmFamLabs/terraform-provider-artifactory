# terraform-provider-artifactory

## Abstract
This provider's is given *some* input, an artifact from JFrog's
Artifactory will create a S3 object deployment that's state can be
traced/tracked for influencing a vanilla lambda deployment that wants an s3
object for deployment.

## Feature Wishlist

- *MVP?* a `resource` and `data`
 - [ ] a `resource` that is a `artifactory_artifact_s3_deployment` which will produce a
   resource that will have the attributes of `s3_bucket` and `s3_key` that can be
   used to feed a vanilla
   [`aws_lambda_function`](https://www.terraform.io/docs/providers/aws/r/lambda_function.html)
   resource defined with `s3_bucket` and `s3_key` from the source resource. 
   - Changes to the source `artifactory_artifact_s3_deployment` resource
     will effect the `s3_key` through a hash if Artifactory deployment isn't
     already doing this for us. (I think it is), but a `s3_bucket` 
     should be required and we should expect to use a vanilla
     `aws_s3_bucket` data/resource attribute.
   - `source_code_hash` on the `aws_lambda_function` can also use the one
     provided conveniently by artifactory. i.e.;
     
     ```terraform
     ```
   - the "deployment" resource should signify the object is placed into s3 and will
     have the appropriate CRUD operations expected of a `resource` in terraform
     - broken down, the resource name follows terraform's pattern of jumping
       between domains with `s3_deployment` as a 'destination'. `deployment` to
       signify this is put as one of potential other actions on `s3`
 - [X] a `data` that is a `artiactory_artifact` which will utilize Artifactory's
   REST API for grabbing information of a given `repository_path`


### Caveat for Terraform Enterprise

Via [terraform custom providers][tfe_custom_providers]

Custom providers (plugins) have to either have their binary added to the repo
in the same path of any terraform that would require it, or it must be packaged
and bundled using the [`terraform_bundle`][terraform_bundle] tool.


# Terraform Resources

_theres probably a way to generate this_

## `data.artifactory_artifact`
```terraform
data "artifactory_artifact" "test_artifact" {
  repository_path = "lambda/propinc/ingest/replicate-2.30.0.zip"
}

output "checksums" {
  value = data.artifactory_artifact.test_artifact.checksums
}

output "download_uri" {
  value = data.artifactory_artifact.test_artifact.download_uri
}

output "path" {
  value = data.artifactory_artifact.test_artifact.path
}

output "repo" {
  value = data.artifactory_artifact.test_artifact.repo
}

output "size" {
  value = data.artifactory_artifact.test_artifact.size
}
```



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


[terraform_bundle]: https://github.com/hashicorp/terraform/tree/master/tools/terraform-bundle#installing-a-bundle-in-on-premises-terraform-enterprise
[tfe_custom_providers]: https://www.terraform.io/docs/cloud/run/install-software.html#custom-and-community-providers