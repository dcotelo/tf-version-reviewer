# tf-version-reviewer

Check Terraform version across recursive directories

  

# Install

  

`$ go get`

`$ go build`

  

# Usage

  `$ ./tfversion`
  

Usage arguments:

`-d`  Base project directory  that contains all terraform configurations

`-tf` Search .terraform-version and versions.tf file in project -default:false

`-v`  Display each configuration absolute path -default:false
