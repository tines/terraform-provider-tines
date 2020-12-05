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

* You may fail at this step, check your distribution and update the Makefile "OS_ARCH" if not using OSX (i.e. `linux_amd64` for Linux)

## Authentication

Authentication parameters can be set as environment variables

```
export TF_VAR_tines_email=example@email.com
export TF_VAR_tines_email=token
export TF_VAR_tines_base_url=https://dappled-horse-1234.tines.io/
```

# Circular Logic

Tines utilizes circular logic and loops frequently, however, that is a barrier in using Terraform effectively because it cannot know which resources to create first. When running `terraform plan` or `terraform apply` you may run into an error like:

`Error: Cycle ....`

In order to get passed this error, set the source and receiver IDs of the Agent which finishes the loop and remove the dependency from the beginning of the loop. 

[More information](https://serverfault.com/questions/1005761/what-does-error-cycle-means-in-terraform#:~:text=When%20Terraform%20returns%20this%20error,that%20it's%20no%20longer%20contradictory.&text=The%20%2Ddraw%2Dcycles%20command%20causes,reported%20using%20the%20color%20red.)
