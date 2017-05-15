variable "user" {}
variable "password" {}
variable "nsx_manager" {}
variable "allow_unverified_ssl" {}

provider "nsx" {
  user = "${var.user}"
  password = "${var.password}"
  nsx_manager = "${var.nsx_manager}"
  nsx_version = "6.3"
  allow_unverified_ssl = "${var.allow_unverified_ssl}"
}

resource "nsx_tag" "tag1" {
  tag_name = "Terraform.Test1"
  description = "Test"
  create_if_missing = true
  persistent = false
}