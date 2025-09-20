## README.md (installation & usage)
helm plugin to create helm chart from docker compose
# helm-compose-plugin — Go (multi-file repo)

**What I changed**
- Converted the single-file prototype into a small multi-file repository with a Cobra-based CLI, `plugin.yaml` for Helm plugin installation, a `README.md`, basic unit tests, and a simple converter implementation separated into `internal/convert`.
- The project layout below is ready to `go build` and `helm plugin install` (installation steps provided in README).

---

## Repository layout

```
helm-compose-plugin/
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


```markdown
# helm-compose-plugin

Generate a Helm chart from a `docker-compose.yml` file.

## Install locally as a Helm plugin

1. Build the binary:

```bash
go build -o helm-compose
```

2. Copy binary and plugin.yaml into a folder under your helm plugins dir, e.g. on Unix:

```bash
mkdir -p ~/.local/share/helm/plugins/helm-compose
cp helm-compose plugin.yaml ~/.local/share/helm/plugins/helm-compose/
```

3. Now run:

```bash
helm-compose -f docker-compose.yml -o ./mychart -n mychart
```

(Or you can run `./helm-compose` directly without installing as a plugin.)

## Notes & limitations
- This is a best-effort starter — it handles images, simple ports, env, and basic templates.
- Complex volumes, build contexts, and networks are emitted into `values.yaml` for manual handling.
- Contributions welcome!
```

---

### Next steps I can do for you (pick any, or I will pick one):
- Add more robust parsing of `ports` (host:container, ranges)
- Generate StatefulSets for services with named volumes
- Add CRD support / k8s Ingress generation
- Create a GitHub Actions workflow to build and release the plugin

Tell me which of the next steps you'd like and I will implement it directly in the canvas.
