package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"os/exec"
	"os"
)

func resourceUpload() *schema.Resource {
	return &schema.Resource{
		Create: resourceUploadCreate,
		Read:   resourceUploadRead,
		Update: resourceUploadUpdate,
		Delete: resourceUploadDelete,

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
	file_name := d.Get("file_name").(string)
	d.SetId(file_name)
	c := exec.Command("curl", "-i", "-X", "POST", "-F", "file=@" + file_name, "http://" + ip_address + ":8080/upload")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

	if err != nil {
		panic("The flask server may not be running")
	}
	return nil
}

func resourceUploadRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUploadUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUploadDelete(d *schema.ResourceData, m interface{}) error {
	c := exec.Command("curl", "-i", "-X", "DELETE", "-F", "file=@" + d.Get("file_name").(string), "http://" + d.Get("ip_address").(string) + ":8080/upload")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err !=nil {
		panic(err)
	}
	return nil
}