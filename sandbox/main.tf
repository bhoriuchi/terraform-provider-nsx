variable "user" {}
variable "password" {}
variable "nsx_manager" {}
variable "allow_unverified_ssl" {}

provider "nsx" {
  user = "${var.user}"
  password = "${var.password}"
  nsx_manager = "${var.nsx_manager}"
  allow_unverified_ssl = "${var.allow_unverified_ssl}"
}

resource "nsx_tag" "my-tag" {
  id = "stuff"
}