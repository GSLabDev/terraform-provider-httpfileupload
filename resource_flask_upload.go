package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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
		Update: resourceUploadUpdate,
		Delete: resourceUploadDelete,
		Schema: map[string]*schema.Schema{
			"host_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IP_ADDR", nil),
				Description: descriptions["host_url"],
			},

			"file_path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FILENAME", ""),
				Description: descriptions["file_path"],
			},
		},
	}
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
		fmt.Println("error writing to buffer")
		return fmt.Errorf("[ERROR] Error in writing to buffer %s", err)
	}
	fh, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening file")
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
	fmt.Println(string(resp_body))
	d.SetId(file)
	return nil
}

func resourceUploadRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUploadUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}
func resourceUploadDelete(d *schema.ResourceData, m interface{}) error {
	hostUrl := d.Get("host_url").(string)
	filePath := d.Get("file_path").(string)
	_, file := filepath.Split(filePath)

	in := []byte(`{ "file":"" }`)
	var raw map[string]string
	json.Unmarshal(in, &raw)
	raw["file"] = file
	out, _ := json.Marshal(raw)
	var jsonStr = []byte(out)
	req, err := http.NewRequest("DELETE", hostUrl, bytes.NewBuffer(jsonStr))
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
	fmt.Println(resp.Status)
	d.SetId("")
	return nil

}
