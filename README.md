# VMware NSX Provider

The VMware NSX provider is used to interact with the resources supported by VMware NSX. The provider needs to be configured with the proper credentials and manager before it can be used.

```
Note: This provider is experimental and not full-featured
```

### Example Usage

---

```
# Configure the VMware NSX Provider
provider "nsx" {
  user = "${var.user}"
  password = "${var.password}"
  nsx_manager = "${var.nsx_manager}"
  allow_unverified_ssl = "${var.allow_unverified_ssl}"
}

# Create or lookup a Security Tag
resource "nsx_tag" "tag1" {
  tag_name = "Application.web"
  description = "Opens port 80"
  create_if_missing = true
}

# Create or lookup a second Security Tag
resource "nsx_tag" "tag2" {
  tag_name = "Application.secure"
  description = "Opens port 443"
  create_if_missing = true
}

# Assign Security Tags to a Virtual Machine
resource "nsx_vm" "web01" {
  vm_id = "${var.vm_id}"
  security_tags = [
    "${nsx_tag.tag1.id}",
    "${nsx_tag.tag2.id}"
  ]
}
```

### Argument Reference

---

The following arguments are used to configure the VMware NSX Provider:

* `user` - (Required) This is the username for NSX API operations. Can also be specified with the NSX_USER environment variable.
* `password` - (Required) This is the password for NSX API operations. Can also be specified with the NSX_PASSWORD environment variable.
* `nsx_manager` - (Required) This is the NSX manager name for NSX API operations. Can also be specified with the NSX_MANAGER environment variable.
* `nsx_version` - (Optional) This is the version of the NSX manager. It is used for determining what API features are available and defaults to `6.3`. Can also be specified with the NSX_VERSION environment variable.
* `allow_unverified_ssl` - (Optional) Boolean that can be set to true to disable SSL certificate verification. This should be used with care as it could allow an attacker to intercept your auth token. If omitted, default value is false. Can also be specified with the NSX_ALLOW_UNVERIFIED_SSL environment variable.

# Resources

* [nsx_tag]()
* [nsx_vm]()

### nsx_tag

Provides an NSX security tag resource. This can be used to create, modify, delete, and lookup security tags.

### Example Usage

---

```
resource "nsx_tag" "tag1" {
  tag_name = "Application.web"
  description = "Opens port 80"
  create_if_missing = true
}
```

### Argument Reference

---

The following arguments are supported:

* `tag_id` - (Optional) The id of the NSX security tag (i.e. securitytag-123). Required if `tag_name` is omitted
* `tag_name` - (Optional) The name of the NSX security tag. Required if `tag_id` is omitted. Does a case sensitive lookup on the name and sets resource id and `tag_id` to the security tag id
* `description` - (Optional) The description of the NSX security tag
* `is_universal` - (Optional) Boolean that creates the NSX security as a universal security tag. NSX 6.3 and higher required. Defaults to false
* `create_if_missing` - (Optional) Boolean that creates the the NSX security tag if it is not found. Defaults to false
* `persistent` - (Optional) Boolean that prevents the NSX security tag from being destroyed when true during a destroy operation. This is useful when using this resource for lookup in the `nsx_vm` resource

### nsx_vm

Provides an NSX virtual machine resource. This can be used to attach and detach security tags

### Example Usage

---

```
resource "nsx_vm" "web01" {
  vm_id = "${var.vm_id}"
  security_tags = [
    "${nsx_tag.tag1.id}",
    "${nsx_tag.tag2.id}"
  ]
}
```

### Argument Reference

---

The following arguments are supported:

* `vm_id` - (Required) The vSphere managed object reference id or BIOS uuid of the virtual machine
* `security_tags` (Optional) A list of NSX security tag ids or names. Can be used to attach and detach security tags to the virtual machine
