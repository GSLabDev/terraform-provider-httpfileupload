package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"os/exec"
	"os"
        "net/http"
        "io/ioutil"
        "bytes"
        "fmt"
        "io"
        "mime/multipart"
)

func resourceUpload() *schema.Resource {
	return &schema.Resource{
		Create: resourceUploadCreate,
		Read:   resourceUploadRead,
		Update: resourceUploadUpdate,

		Schema: map[string]*schema.Schema{
			"ip_address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IP_ADDR", nil),
				Description: descriptions["ip_address"],
			},

			"file_name": &schema.Schema{
				Type:           schema.TypeString,
				Required:       true,
				DefaultFunc:    schema.EnvDefaultFunc("FILENAME", ""),
				Description:    descriptions["file_name"],
			},
		},
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"ip_address": "The ip address of the machine you want to upload file to.",

		"file_name": "Entire filepath including the filename.",
	}
}


func resourceUploadCreate(d *schema.ResourceData, m interface{}) error {
	ip_address := d.Get("ip_address").(string)
	d.SetId(ip_address)
	filename := d.Get("file_name").(string)
	d.SetId(filename)
        targetUrl:="http:/"+"/"+ip_address+":8080/upload"


        bodyBuf := &bytes.Buffer{}
        bodyWriter := multipart.NewWriter(bodyBuf)
        fileWriter, err := bodyWriter.CreateFormFile("file", filename)
        if err != nil {
            fmt.Println("error writing to buffer")
            return fmt.Errorf("[ERROR] Error in writing to buffer %s",err)
         }
       fh, err := os.Open(filename)
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
       resp, err := http.Post(targetUrl, contentType, bodyBuf)
       if err != nil {
          return fmt.Errorf("[ERROR] Request failed")
         }
       defer resp.Body.Close()
       resp_body, err := ioutil.ReadAll(resp.Body)
       if err != nil {
          return err
        }
       fmt.Println(resp.Status)
       fmt.Println(string(resp_body))
       return nil
}

func resourceUploadRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUploadUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}
