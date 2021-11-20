# terraform-provider-sonarqube

A Terraform provider for SonarQube


```hcl
provider "sonarqube" {
  url   = "https://your-sonarqube-url"
  token = "your-personal-token"
}

# Managing global settings
resource "sonarqube_settings_value" "email_from" {
  key   = "email.from"
  value = "sonarqube@example.org"
}

# Retrieving an existing user
data "sonarqube_user" "my_user" {
  email = "my-email-address"
}

# Create new proejct
resource "sonarqube_project" "my_project" {
  key    = "terraform-testje"
  name   = "Terraform Test Project"
  public = false
}

# Add user to the project
resource "sonarqube_project_user" "myproject__my_user" {
  project_key = sonarqube_project.my_project.key
  login       = data.sonarqube_user.my_user.login
  permission  = ["admin","user"]
}

```
