provider "artifactory" {
  url = "https://artifacts.amfamlabs.com"
}

// we only need READs for MVP, methinks
data "artifactory_artifact" "test_artifact" {
  repository_path = "lambda/propinc/ingest/replicate-2.30.0.zip"
}




// this is the terraform use case
resource "artifactory_artifact_s3_deployment" "test_artifact" {
  repository_path = "lambda/propinc/ingest/fake.zip"
  //body = ""
}


locals {
  repo      = data.artifactory_artifact.test_artifact.repo
  path      = data.artifactory_artifact.test_artifact.path
  checksums = data.artifactory_artifact.test_artifact.checksums
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
