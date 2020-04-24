# terraform-provider-sonarqube

A Terraform provider for SonarQube


```hcl
provider "sonarqube" {
  url   = "https://your-sonarqube-url"
  token = "your-personal-token"
}

data "sonarqube_user" "my_user" {
  email = "my-email-address"
}

resource "sonarqube_project_user" "myproject__my_user" {
  project_key = sonarqube_project.my_project.key
  login       = data.sonarqube_user.my_user.login
  permission  = "admin"
}

resource "sonarqube_project" "my_project" {
  key    = "terraform-testje"
  name   = "Terraform Test Project"
  public = false
}
```
