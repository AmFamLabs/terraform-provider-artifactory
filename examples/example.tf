provider "artifactory" {
  url = "${var.artifactory_url}"
  s3_region = "us-east-1"
  assume_role {
    role_arn     = "${var.iam_role_arn}"
    session_name = "artifactory"
}

// this is the terraform use case
resource "artifactory_artifact_s3_deployment" "test_artifact" {
  repository_path = "lambda/test-1.0.0.zip"
  s3_bucket = "live-tmp-bucket"
  s3_prefix = "lambda/deployments"
  s3_region = "us-east-1"
}


// we only need READs for MVP, methinks
data "artifactory_artifact" "test_artifact" {
  repository_path = "lambda/propinc/ingest/replicate-2.30.0.zip"
}

locals {
  repo      = data.artifactory_artifact.test_artifact.repo
  path      = data.artifactory_artifact.test_artifact.path
  checksums = data.artifactory_artifact.test_artifact.checksums
  s3_prefix = "lambda/test-deployments"
  artifact_name = "test-1.0.0.zip"
}

// this is the input
output "repository_path" {
  value = data.artifactory_artifact.test_artifact.repository_path
}

// these will be computed
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

output "sha256base64" {
  value = base64encode(artifactory_artifact_s3_deployment.test_artifact.checksums.sha256)
}