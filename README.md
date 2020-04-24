# terraform-provider-sonarqube

A Terraform provider for SonarQube


```hcl
resource "sonarqube_project" "my_project" {
  key    = "terraform-testje"
  name   = "Terraform Test Project"
  public = false
}
```
