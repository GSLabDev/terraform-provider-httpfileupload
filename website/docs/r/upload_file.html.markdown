---
layout: "httpfileupload"
page_title: "HttpFileUpload: httpfileupload_file"
sidebar_current: "docs-httpfileupload-resource-inventory-folder"
description: |-
  Provides a Http file upload resource. This can be used to upload file in server.
---

# httpfileupload\_file

Provides a Http file upload resource. This can be used to upload file in server.

## Example Usage

```hcl
# Upload a file
resource "httpfileupload_file" "my_file"{
  host_url = "localhost:8080/upload"
  file_path = "/home/user/Documents/hello.txt"
}
```

## Argument Reference

The following arguments are supported:

* `host_url` - (Required) The Url for http file upload
* `file_path` - (Required) The file path of the File to be uploaded