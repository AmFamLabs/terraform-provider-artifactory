# terraform-provider-artifactory

## Abstract
This provider's is given *some* input, an artifact from JFrog's
Artifactory will create a S3 object deployment that's state can be
traced/tracked for influencing a vanilla lambda deployment that wants an s3
object for deployment.


## Troubleshooting

See this project's [Wiki][wiki].

## Feature Wishlist
<details>
<summary>completed features</summary>

- *MVP?* a `resource` and `data`
 - [X] a `resource` that is a `artifactory_artifact_s3_deployment` which will produce a
   resource that will have the attributes of `s3_bucket` and `s3_prefix` that can be
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
     resource "aws_lambda_function" "fun" {
        source_code_hash = base64encode(data.artifactory_artifact.test_artifact.checksums.sha256)
        s3_bucket = ...
     }
     ```
   - the "deployment" resource should signify the object is placed into s3 and will
     have the appropriate CRUD operations expected of a `resource` in terraform
     - broken down, the resource name follows terraform's pattern of jumping
       between domains with `s3_deployment` as a 'destination'. `deployment` to
       signify this is put as one of potential other actions on `s3`
 - [X] a `data` that is a `artiactory_artifact` which will utilize Artifactory's
   REST API for grabbing information of a given `repository_path`

</details>

### Caveat for Terraform Enterprise

Via [terraform custom providers][tfe_custom_providers]

Custom providers (plugins) have to either have their binary added to the repo
in the same path of any terraform that would require it, or it must be packaged
and bundled using the [`terraform_bundle`][terraform_bundle] tool.


# Terraform Resources

_theres probably a way to generate this_


## `artifactory_artifact` (resource)

One "deployment" should be made per function, since it is fully managed by the
s3_key, the underlying zip on s3 that more than one function may use will be
deleted. We assume artifactory will have uniqueness on each "repository_path",
the artifact itself, or the basename will be the basename on s3 with
`s3_prefix` like `s3://$S3_PREFIX/$ARTIFACT_BASENAME`.

To put another way - if your function will require a specific version, it will
need to be exposed by Artifactory's path provided as `repository_path` to the
resource.

```terraform
resource "artifactory_artifact_s3_deployment" "test_artifact" {
  repository_path = "lambda/propinc/ingest/replicate-2.30.0.zip"
  s3_bucket = "yolk-propinc-live-tmp-bucket"
  s3_prefix = "lambda/deployments"
}


... aws_lambda_function goes here ...
```

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
[wiki]: https://git.amfamlabs.com/terraform/terraform-provider-artifactory/-/wikis/home