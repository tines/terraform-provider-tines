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
