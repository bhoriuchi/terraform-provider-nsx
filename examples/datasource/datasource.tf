provider "nsx" {
  # user = "nsx-user"
  # password = "nsx-user-password"
  # nsx_manager = "nsx-manager-fqdn"
  allow_unverified_ssl = "true"
}

data "nsx_security_tag" "tag1" {
  name_regex = "(?i)terraform"
}

output "tag1_id" {
  value = "${data.nsx_security_tag.tag1.id}"
}

output "tag1_name" {
  value = "${data.nsx_security_tag.tag1.name}"
}

output "tag1_description" {
  value = "${data.nsx_security_tag.tag1.description}"
}

output "tag1_vm_count" {
  value = "${data.nsx_security_tag.tag1.vm_count}"
}