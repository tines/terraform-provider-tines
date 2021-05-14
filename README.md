# Terraform Provider for Tines.io

## Requirements

Terraform `0.13.x` or greater

Go `1.12.x` or greater

## Install

```
mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
git clone https://github.com/tuckner/terraform-provider-tines.git
cd $GOPATH/src/github.com/terraform-providers/terraform-provider-tines
export GO111MODULE="on"
make install
```

* You may fail at this step, check your distribution and update the Makefile "OS_ARCH" (i.e. `linux_amd64` for Linux)

## Authentication

Authentication parameters can be set as environment variables

```
export TINES_EMAIL=example@email.com
export TINES_TOKEN=token
export TINES_URL=https://dappled-horse-1234.tines.io/
```

Parameters can also be set in the provider configuration

```
provider "tines" {
    email    = var.tines_email
    base_url = var.tines_base_url
    token    = var.tines_token
}
```


## Examples

More examples can be found here: https://github.com/tuckner/tines-example-stories

## Export Conversion

Tines exports can be transformed into Terraform files utilizing the `export2terraform` script in the scripts directory. Alternatively, there is a service available that will convert story exports to Terraform files and email the resulting files:

https://quiet-vista-5142.tines.io/forms/91784f6d80499f2810ab9a31d0c15b72

# Circular Logic

Tines utilizes circular logic and loops frequently, however, that is a barrier in using Terraform effectively because it cannot know which resources to create first. When running `terraform plan` or `terraform apply` you may run into an error like:

`Error: Cycle ....`

In order to get passed this error, set the source and receiver IDs of the first agent which begins the loop and remove the agent receiver ID from the agent which triggers the loop in order remove the circular reference. 

[More information](https://serverfault.com/questions/1005761/what-does-error-cycle-means-in-terraform#:~:text=When%20Terraform%20returns%20this%20error,that%20it's%20no%20longer%20contradictory.&text=The%20%2Ddraw%2Dcycles%20command%20causes,reported%20using%20the%20color%20red.)

# Troubleshooting

If Terraform receives errors from the Tines API, below are the most common causes.

## Error Codes

`500` - The body of a request is most likely malformed in a way that Tines could not understand.

`422` - The agent/resource already exists or there is a bug in the in the logic to create/destroy the resource. Try to start with a fresh Terraform state and check that unique global resource names are being created.
