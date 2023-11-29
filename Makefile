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
	brew install terraform
	go install golang.org/x/tools/gopls@latest
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Generate docs
.PHONY: docs
docs:
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Run local install
.PHONY: install-local
install:
	go install
