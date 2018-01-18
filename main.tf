provider "http_file_upload" {
  ip_address = "localhost"
}

resource "http_file_upload_file" "my_file"{
  file_name = "/home/samtholiya/Documents/hello.txt"
}
