provider "upload" {
  ip_address = "10.136.50.28"
}

resource "upload_file" "my_file"{
  file_name = "/home/ismail/Documents/hello.txt"
}