---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarqube_project Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  
---

# sonarqube_project (Resource)



## Example Usage

```terraform
resource "sonarqube_project" "my-project" {
  key    = "my-key"
  name   = "my-project"
  public = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String)
- `name` (String)

### Optional

- `public` (Boolean)

### Read-Only

- `id` (String) The ID of this resource.
