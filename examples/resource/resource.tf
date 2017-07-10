provider "nsx" {
  # user = "nsx-user"
  # password = "nsx-user-password"
  # nsx_manager = "nsx-manager-fqdn"
  allow_unverified_ssl = "true"
}

data "nsx_security_tag" "tag1" {
  name_regex = "(?i)terraform"
}

data "nsx_security_tag" "tag2" {
  name_regex = "(?i)^anti_virus.virusfound.threat=high$"
}

resource "nsx_security_tag" "tag1" {
  name = "Terraform.Provider.Test=2"
  description = "Test 1 of terraform provider"
}

resource "nsx_vm" "vm1" {
  vm_id = "vm-33"
  security_tags = [
    "${nsx_security_tag.tag1.id}",
    "${data.nsx_security_tag.tag1.id}"
  ]
}