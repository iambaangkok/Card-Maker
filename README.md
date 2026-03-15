# Card Maker

- read CSV (legacy mode) or YAML/JSON (generic mode)
- read HTML template
- put data from data files into HTML template
- export cards as image
  - .PNG

- You need the fonts in you computer
- Not sure if image importing works ):
  - might have to serve the static images with http.FileServer

## Usage

### Legacy mode (current weapon/part/item flow)

Runs the existing CSV-based pipeline:

```bash
go run ./cmd/card-maker
```

Configuration comes from `app.yaml` and the `csvs/` directory.

### Generic mode (configurable projects)

The generic engine is enabled by passing a `--project` flag pointing to a
project configuration file (YAML or JSON) that describes card types, data
files, and templates:

```bash
go run ./cmd/card-maker --project ./configs/example-project.yaml
```

See `configs/example-project.yaml` for the basic structure of a project
configuration. Card records are defined in YAML/JSON files referenced from
`data_dir`, and templates live under `template_dir`.
