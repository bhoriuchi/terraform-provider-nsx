variable "user" {}
variable "password" {}
variable "nsx_manager" {}
variable "allow_unverified_ssl" {}
variable "vm_id" {}

provider "nsx" {
  user = "${var.user}"
  password = "${var.password}"
  nsx_manager = "${var.nsx_manager}"
  allow_unverified_ssl = "${var.allow_unverified_ssl}"
}

resource "nsx_tag" "tag1" {
  tag_name = "Terraform.Test1"
  description = "Test"
  create_if_missing = true
}

resource "nsx_tag" "tag2" {
  tag_name = "Terraform.Test2"
  description = "Test"
  create_if_missing = true
}

resource "nsx_vm" "puppet01" {
  vm_id = "${var.vm_id}"
  security_tags = [
    "${nsx_tag.tag1.id}",
    "Terraform.Test2"
  ]
}