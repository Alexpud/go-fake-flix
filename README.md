# go-fake-flix

## Run the API

```bash
go run ./cmd/server/main.go
```

The server listens on `:8080`.

## Swagger / OpenAPI docs (swaggo + Scalar)

This project uses:

- `swaggo/swag` to **generate** the Swagger/OpenAPI spec files into `./docs/`
- Scalar (via `go-scalar-api-reference`) to **render** the spec at `/reference`

### 1) Install the `swag` generator

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Make sure your Go bin directory is on your PATH:

- Windows (PowerShell):

```powershell
$env:PATH += ';' + (go env GOPATH) + '\bin'
```

### 2) Generate the spec files

From the repo root:

```bash
swag init -g ./cmd/server/main.go -o ./docs
```

This generates:

- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

### 3) View the docs in the browser

Start the API:

```bash
go run ./cmd/server/main.go
```

Then open:

- Scalar UI: `http://localhost:8080/reference`
- OpenAPI 3 (converted): `http://localhost:8080/openapi.json`
- Raw Swagger 2.0: `http://localhost:8080/docs/swagger.json`

### Common issues

- **`cannot find type definition: <TypeName>`**: one of your Swagger annotations references a Go type that doesn’t exist
  (example: `@Success 200 {object} ResponseSuccess`). Create the type or update the annotation.
- **Blank `/reference` page / URL parse errors on Windows**: Scalar needs an **HTTP URL** for `SpecURL` (not a Windows `file://C:\...` path).
