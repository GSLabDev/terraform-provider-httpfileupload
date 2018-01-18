package httpfileupload

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("HTTPFILEUPLOAD_HOST_URL"); v == "" {
		t.Fatal("HTTPFILEUPLOAD_HOST_URL must be set for acceptance tests")
	}

	if v := os.Getenv("HTTPFILEUPLOAD_FILE_PATH"); v == "" {
		t.Fatal("HTTPFILEUPLOAD_FILE_PATH must be set for acceptance tests")
	}

	if _, err := os.Stat(os.Getenv("HTTPFILEUPLOAD_FILE_PATH")); os.IsNotExist(err) {
		// path/to/whatever does not exist
		t.Fatal("Path is Invalid")
	}
}
func TestAccHttpFileUploadFile_Basic(t *testing.T) {
	resourceName := "httpfileupload_file.my_file"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHttpFileUploadDestroy(resourceName),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckHTTPFileUploadBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHttpFileUploadExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckHttpFileUploadDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("[ERROR] Resource Not found")
		}

		hostURL := rs.Primary.Attributes["host_url"]
		filePath := rs.Primary.Attributes["file_path"]
		_, file := filepath.Split(filePath)
		resp, err := http.Get(hostURL + "/" + file)
		if err != nil {
			log.Println("[ERROR] Error ", err)
		}
		if resp.StatusCode == 200 {
			return fmt.Errorf("[ERROR] File with name " + file + "was found")
		}
		return nil
	}
}

func testAccCheckHttpFileUploadExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No File ID is set")
		}
		hostURL := rs.Primary.Attributes["host_url"]
		filePath := rs.Primary.Attributes["file_path"]
		_, file := filepath.Split(filePath)
		resp, err := http.Get(hostURL + "/" + file)
		if err != nil {
			log.Println("[ERROR] Error ", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("[ERROR] File with name " + file + "was not found")
		}
		return nil
	}
}

func testAccCheckHTTPFileUploadBasic() string {
	return fmt.Sprintf(`
provider "httpfileupload" {

}

resource "httpfileupload_file" "my_file" {
  host_url = "%s"
  file_path = "%s"
}
`, os.Getenv("HTTPFILEUPLOAD_HOST_URL"), os.Getenv("HTTPFILEUPLOAD_FILE_PATH"))
}
