provider "artifactory" {
  url = "https://artifacts.amfamlabs.com"
}

// we only need READs for MVP, methinks
data "artifactory_artifact" "test_artifact" {
  repository_path = "lambda/propinc/ingest/fake.zip"
}

resource "null_resource" "use_data" {
  triggers = {
    repository_path = data.artifactory_artifact.test_artifact.repository_path
    repo = data.artifactory_artifact.test_artifact.repo
    path = data.artifactory_artifact.test_artifact.path
  }
}

// this is the terraform use case
//resource "artifactory_artifact" "test_artifact" {
//  repository_path = "lambda/propinc/ingest/fake.zip"
//  body = ""
//}
