default: tests

# Run tests including acceptance
.PHONY: tests
tests:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Run linter fix
.PHONY: lintfix
lintfix:
	golangci-lint run --fix

# Run local install
.PHONY: install
install:
	brew install golangci-lint
	brew tap hashicorp/tap
	brew install hashicorp/tap/terraform
	go install golang.org/x/tools/gopls@latest
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# Generate docs
.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# Run local install
.PHONY: install-local
install-local:
	go install
