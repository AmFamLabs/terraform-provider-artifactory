provider "artifactory" {
  url = "https://artifacts.amfamlabs.com"
}

// we only need READs for MVP, methinks
data "artifactory_artifact" "test_artifact" {
  repository_path = "lambda/propinc/ingest/fake.zip"
}

// this is the terraform use case
//resource "artifactory_artifact" "test_artifact" {
//  repository_path = "lambda/propinc/ingest/fake.zip"
//  body = ""
//}
