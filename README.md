# terraform-provider-artifactory

This provider's goal is given *some* input, an artifact from JFrog's
Artifactory will create a S3 object deployment that's state can be
traced/tracked.

## Feature Wishlist

- *MVP?* given an artifact path use resource outputs to provide input to an
  s3_bucket_object
- resource content doesn't get written to state, but `aws_s3_bucket_object::content_base64` will
- `aws_s3_bucket_object:content*` isn't written to tfstate somehow
- a resource that handles s3 through golang equivalent of boto calls bypassing `terraform_provider_aws`



## Development

### Prerequisites

- `go` installed

### Make plugin

```sh
go build
```
