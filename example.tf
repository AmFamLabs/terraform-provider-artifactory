provider "artifactory" {
  url = "https://artifacts.amfamlabs.com"
  //will work :)
  //assume_role {
  //  //role_arn     = "arn:aws:iam::${data.consul_keys.config.var["${var.workspace}_id"]}:role/terraform_enterprise"
  //  //role_arn     = "arn:aws:iam::some_id:role/terraform_enterprise"
  //  role_arn = "arn:aws:iam::REDACTED:role/core-operator"
  //  session_name = "artifactory"
  //}
}

// we only need READs for MVP, methinks
data "artifactory_artifact" "test_artifact" {
  repository_path = "lambda/propinc/ingest/replicate-2.30.0.zip"
}




// this is the terraform use case
resource "artifactory_artifact_s3_deployment" "test_artifact" {
  repository_path = "lambda/propinc/ingest/replicate-2.30.0.zip"
  s3_bucket = "yolk-propinc-live-tmp-bucket"
  s3_prefix = "lambda/deployments"
}

resource "artifactory_artifact_s3_deployment" "test_artifact_deux" {
  repository_path = "lambda/propinc/ingest/${local.artifact_name}"
  s3_bucket = "yolk-propinc-live-tmp-bucket"
  s3_prefix = local.s3_prefix
}


locals {
  repo      = data.artifactory_artifact.test_artifact.repo
  path      = data.artifactory_artifact.test_artifact.path
  checksums = data.artifactory_artifact.test_artifact.checksums
  s3_prefix = "lambda/deployments-deux/"
  artifact_name = "replicate-2.30.0.zip"
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
