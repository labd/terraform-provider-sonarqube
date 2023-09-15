resource "sonarqube_project" "my-project" {
  key    = "my-key"
  name   = "my-project"
  public = true
}
