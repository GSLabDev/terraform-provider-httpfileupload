---
layout: "httpfileupload"
page_title: "Provider: Http File Upload"
sidebar_current: "docs-httpfileupload-index"
description: |-
  The Http File Upload provider is used to upload file to the host website
---

# HttpFileUpload Provider

The HttpFileUpload provider is used upload file to the host website.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the vRealize Automation Provider
provider "httpfileupload" {

}

resource "httpfileupload_file" "my_file"{
  host_url = "localhost:8080/upload"
  file_path = "/home/samtholiya/Documents/hello.txt"
}
```

## Acceptance Tests

The Active Directory provider's acceptance tests require the above provider
configuration fields to be set using the documented environment variables.

Once all these variables are in place, the tests can be run like this:

```
make testacc TEST=./builtin/providers/httpfileupload
```
