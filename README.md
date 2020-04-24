# terraform-provider-sonarqube

A Terraform provider for SonarQube


```hcl
provider "sonarqube" {
  url   = "https://your-sonarqube-url"
  token = "your-personal-token"
}

resource "sonarqube_project" "my_project" {
  key    = "terraform-testje"
  name   = "Terraform Test Project"
  public = false
}
```
