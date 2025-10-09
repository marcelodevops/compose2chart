## README.md (installation & usage)
helm plugin to create helm chart from docker compose
# helm-compose2chart-plugin — Go (multi-file repo)

**What I changed**
- Converted the single-file prototype into a small multi-file repository with a Cobra-based CLI, `plugin.yaml` for Helm plugin installation, a `README.md`, basic unit tests, and a simple converter implementation separated into `internal/convert`.
- The project layout below is ready to `go build` and `helm plugin install` (installation steps provided in README).

---

## Repository layout

```
helm-compose2chart-plugin/
├── plugin.yaml
├── README.md
├── go.mod
├── cmd/
│   └── root.go
├── main.go
├── internal/
│   └── convert/
│       ├── convert.go
│       └── convert_test.go
└── templates/
    └── _helpers.tpl
```

---



# helm-compose2chart-plugin

Generate a Helm chart from a `docker-compose.yml` file.

## Install locally as a Helm plugin

1. Build the binary:

```bash
go build -o helm-compose2chart
```

2. Copy binary and plugin.yaml into a folder under your helm plugins dir, e.g. on Unix:

```bash
mkdir -p ~/.local/share/helm/plugins/helm-compose2chart
cp helm-compose2chart plugin.yaml ~/.local/share/helm/plugins/helm-compose2chart/
```
  On Mac OS

  ```bash
mkdir -p ~/Library/helm/plugins/helm-compose2chart
cp helm-compose2chart plugin.yaml ~/Library/helm/plugins/helm-compose2chart/
  ```
3. Now run:

```bash
helm-compose2chart -f docker-compose.yml -o ./mychart -n mychart
```

(Or you can run `./helm-compose2chart` directly without installing as a plugin.)

## Notes & limitations
- This is a best-effort starter — it handles images, simple ports, env, and basic templates.
- Complex volumes, build contexts, and networks are emitted into `values.yaml` for manual handling.
- Contributions welcome!



