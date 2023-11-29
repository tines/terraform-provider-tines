# Install

- `make install`

# Setup

- `go mod download`

# Run tests

- `make tests`

# Run linter

- `make lint`

## Docs

- `make docs`

# VSCode golang specific settings

Extensions: `hashicorp.terraform` and `golang.go`

Config:

```json
  "[go]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "golang.go"
  },
  "go.lintFlags": ["--enable-all", "--new"],
  "[terraform]": {
    "editor.defaultFormatter": "hashicorp.terraform",
    "editor.formatOnSave": true
  }
```

# .terraformrc

You will need the following in your `~/.terraformrc` for local development

```terraform
# ~/.terraformrc
provider_installation{
    dev_overrides{
        "tines/tines" = "REPLACE THIS WITH YOUR $GOPATH"
    }

    direct{}
}
```

# Using terraform-provider-tines locally

- Ensure above steps are performed
- Run `make install-local`
