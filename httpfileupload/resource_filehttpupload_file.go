package httpfileupload

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUpload() *schema.Resource {
	return &schema.Resource{
		Create: resourceUploadCreate,
		Read:   resourceUploadRead,
		Delete: resourceUploadDelete,
		Schema: map[string]*schema.Schema{
			"host_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("FILEHTTPUPLOAD_HOST_URL", nil),
				Description: descriptions["host_url"],
			},

			"file_path": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				DefaultFunc:  schema.EnvDefaultFunc("FILEHTTPUPLOAD_FILE_PATH", nil),
				Description:  descriptions["file_path"],
				ValidateFunc: validatePath,
			},
		},
	}
}
func validatePath(v interface{}, k string) (warnings []string, errors []error) {
	operation := v.(string)
	log.Printf("[DEBUG] Validating Path %s", operation)
	if _, err := os.Stat(os.Getenv("HTTPFILEUPLOAD_FILE_PATH")); os.IsNotExist(err) {
		return returnError("Path does not exist", fmt.Errorf("[ERROR] Invalid Path"))
	}
	return nil, nil
}
func returnError(message string, err error) (warnings []string, errors []error) {
	var errorVar []error
	var warningVar []string
	return append(warningVar, message), append(errorVar, err)
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"host_url": "The ip address of the machine you want to upload file to.",

		"file_path": "Entire filepath including the filename.",
	}
}

func resourceUploadCreate(d *schema.ResourceData, m interface{}) error {
	hostUrl := d.Get("host_url").(string)
	filePath := d.Get("file_path").(string)
	_, file := filepath.Split(filePath)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", file)
	if err != nil {
		log.Println("error writing to buffer")
		return fmt.Errorf("[ERROR] Error in writing to buffer %s", err)
	}
	fh, err := os.Open(filePath)
	if err != nil {
		log.Println("error opening file")
		return fmt.Errorf("[ERROR] Error while opening file %s ", err)
	}

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(hostUrl, contentType, bodyBuf)
	if err != nil {
		return fmt.Errorf("[ERROR] Request failed")
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(resp_body))
	d.SetId(file)
	return nil
}

func resourceUploadRead(d *schema.ResourceData, m interface{}) error {
	hostUrl := d.Get("host_url").(string)
	filePath := d.Get("file_path").(string)
	_, file := filepath.Split(filePath)
	resp, err := http.Get(hostUrl + "/" + file)
	if err != nil {
		log.Println("[ERROR] Error ", err)
	}
	if resp.StatusCode == 200 {
		return nil
	}

	d.SetId("")
	return nil
}

func resourceUploadDelete(d *schema.ResourceData, m interface{}) error {
	err := resourceUploadRead(d, m)
	if d.Id() == "" {
		return fmt.Errorf("[ERROR] Resource does not exist")
	}
	hostUrl := d.Get("host_url").(string)
	filePath := d.Get("file_path").(string)
	_, file := filepath.Split(filePath)
	req, err := http.NewRequest("DELETE", hostUrl+"/"+file, nil)
	if err != nil {
		log.Println("[ERROR] Error ", err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] Error while requesting Token ", err)
	}
	log.Println(resp.Status)
	d.SetId("")
	return nil
}
