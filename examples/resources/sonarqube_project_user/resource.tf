resource "sonarqube_project_user" "my-project-user" {
  login       = "my-login"
  permissions = ["user"]
  project_key = "my-project-key"
}
