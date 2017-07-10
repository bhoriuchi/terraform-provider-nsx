# VMware NSX Provider

The VMware NSX provider is used to interact with the resources supported by VMware NSX. The provider needs to be configured with the proper credentials and manager before it can be used.

* [Provider](#vmware-nsx-provider)
* [Data Sources](#data-sources)
  * [nsx_security_tag](#nsx_security_tag)
  * [nsx_vm](#nsx_vm)
* [Resources](#resources)
  * [nsx_security_tag](#nsx_security_tag-1)
  * [nsx_vm](#nsx_vm-1)


```
Note: This provider is experimental and not full-featured
```

**Example Usage**

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
resource "nsx_security_tag" "tag1" {
  name = "Application.web"
  description = "Opens port 80"
}

# Create or lookup a second Security Tag
resource "nsx_security_tag" "tag2" {
  name = "Application.secure"
  description = "Opens port 443"
}

# Assign Security Tags to a Virtual Machine
resource "nsx_vm" "web01" {
  vm_id = "${var.vm_id}"
  security_tags = [
    "${nsx_security_tag.tag1.id}",
    "${nsx_security_tag.tag2.id}"
  ]
}
```

**Argument Reference**

---

The following arguments are used to configure the VMware NSX Provider:

* `user` - (Required) This is the username for NSX API operations. Can also be specified with the NSX_USER environment variable.
* `password` - (Required) This is the password for NSX API operations. Can also be specified with the NSX_PASSWORD environment variable.
* `nsx_manager` - (Required) This is the NSX manager name for NSX API operations. Can also be specified with the NSX_MANAGER environment variable.
* `nsx_version` - (Optional) This is the version of the NSX manager. It is used for determining what API features are available and defaults to `6.3`. Can also be specified with the NSX_VERSION environment variable.
* `allow_unverified_ssl` - (Optional) Boolean that can be set to true to disable SSL certificate verification. This should be used with care as it could allow an attacker to intercept your auth token. If omitted, default value is false. Can also be specified with the NSX_ALLOW_UNVERIFIED_SSL environment variable.

# Data Sources

* [nsx_security_tag](#nsx_security_tag)
* [nsx_vm](#nsx_vm)

### nsx_security_tag

Looks up an nsx tag by name using a regex search

**Example Usage**

---

```
data "nsx_security_tag" "tag1" {
  name_regex = "(?i)web"
}
```

**Argument Reference**

The following arguments are supported:

* `name_regex` - The regex string to use when searching for a tag

**Attributes Reference**

id is set to the ID of the first matching security tag. In addition, the following
attributes are exported:

* `name` - The name of the security tag
* `is_universal` - Boolean. Tag is universal when true
* `type_name` - NSX type name for the SecurityTag
* `description` - Tags description
* `vm_count` - Number of vms the tag is attached to


### nsx_vm

Looks up vm by id

**Example Usage**

---

```
data "nsx_vm" "vm1" {
    vm_id = "vm-123"
}
```

**Argument Reference**

---

The following arguments are supported

* `vm_id` - Virtual Machine id

**Attributes Reference**

id is set to the ID of the vm. In addition, the following attributes are exported:

* `vm_id` - The id passed in
* `security_tag_ids` - A list of the security tag ids attached to the vm
* `security_tag_names` - A list of the security tag names attached to the vm

---
---

# Resources

* [nsx_security_tag](#nsx_security_tag-1)
* [nsx_vm](#nsx_vm-1)

### nsx_security_tag

Provides an NSX security tag resource. This can be used to create, modify, delete, and lookup security tags.

**Example Usage**

---

```
resource "nsx_security_tag" "tag1" {
  tag = "Application.web"
  description = "Opens port 80"
}
```

**Argument Reference**

---

The following arguments are supported:

* `name` - (Optional) The name of the NSX security tag to apply
* `description` - (Optional) The description of the NSX security tag
* `is_universal` - (Optional) Boolean that creates the NSX security as a universal security tag. NSX 6.3 and higher required. Defaults to false
* `persistent` - (Optional) Boolean that prevents the NSX security tag from being destroyed when true during a destroy operation. This is useful when using this resource for lookup in the `nsx_vm` resource. Defaults to false
* `safe_destroy` - (Optional) Boolean that prevents the NSX security tag from being destroyed when one or more virtual machines are attached to it. Default to true

### nsx_vm

Provides an NSX virtual machine resource. This can be used to attach and detach security tags

**Example Usage**

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

**Argument Reference**

---

The following arguments are supported:

* `vm_id` - (Required) The vSphere managed object reference id or BIOS uuid of the virtual machine
* `security_tags` (Optional) A list of NSX security tag ids or names. Can be used to attach and detach security tags to the virtual machine

# Build
* For use on the local machine run `make bin` from the root directory of this project
* For specific os/architecture set the environment variable `GOX_OS_ARCH` using [gox](https://github.com/mitchellh/gox) os/arch combinations like "darwin/amd64" use `make dist`.
This also requires that docker is installed on the machine building the distributions

# Runtime
Dependencies on [resty](https://github.com/go-resty/resty) result in dynamic bindings to net in glibc. (guessing)  
This will cause Terraform to fail to exec the provider in alpine containers like `hashicorp/terraform`.

Use `stealthybox/infra` for an alpine terraform with glibc:
```bash
docker run -v/$PWD://terra -w//terra stealthybox/infra terraform plan
```
... or just run terraform on your local machine like a normal person.
