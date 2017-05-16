# Sandbox Setup/Run

1. Create a `terraform.tfvars` file with your environemnt's specifics and save it to `/sandbox`

**Example**

```
user = "administrator@vsphere.local"
password = "password"
nsx_manager = "nsx-manager.directv.com"
allow_unverified_ssl = "true"
vm_id = "vm-49"
```

2. Run build to build the nsx provider
3. Modify `main.tf` to your liking
4. Run `terraform apply --var-file "terraform.tfvars"`