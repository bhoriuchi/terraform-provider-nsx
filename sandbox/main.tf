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

resource "nsx_tag" "tag1" {
  tag_id = "securitytag-9"
}

resource "nsx_tag" "tag2" {
  tag_id = "securitytag-3"
}