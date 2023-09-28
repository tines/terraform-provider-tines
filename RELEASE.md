# Release

Releases are managed via `goreleaser` and triggered in Github Actions via `release.yml`. The build is kicked off on a git tag. Steps:

- `git tag v0.1.0`
- `git push origin --tags`
