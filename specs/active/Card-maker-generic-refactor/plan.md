# Technical Plan: Generic, Data-Driven Card Maker

**Task ID:** cook-foodcardgame-card-maker-refactor  
**Status:** Ready for Implementation  
**Based on:** `spec.md`

---

## 1. System Architecture

**Overview**

A single Go binary acts as a **generic card rendering engine**. It is configured per project (game) via YAML/JSON configuration files and card data files. The engine:

1. Loads **project configuration** (renderer settings, card type definitions, data locations, template paths).
2. Builds an in-memory **card type registry** from schemas.
3. Loads **card records** for each type from YAML/JSON into a generic representation.
4. For each card instance:
   - Binds it to the appropriate **Go HTML template**.
   - Produces an HTML string.
   - Hands it to the existing **Chrome renderer** to emit PNG (and optionally PDF) into an output directory.

**Architecture Decisions**

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Language | Go (existing) | Reuse current codebase and ecosystem. |
| Data format for card records | YAML or JSON files | Human-editable, supports nested structures, fits spec. |
| Config format | YAML (primary), with optional JSON | Align with existing `app.yaml`, keep flexibility. |
| Templating engine | Go `html/template` | Already used, safe and powerful for HTML generation. |
| Rendering engine | `chromedp` driving headless Chrome | Already implemented and tuned; supports PNG/PDF. |
| Card model in Go | Generic map-based model + optional typed helpers | Allow new card types without new structs; preserve convenience for some types. |
| Project structure | Project config + per-project data/templates directories | Clean separation between games. |

---

## 2. Technology Stack

**Stack Table**

| Layer | Technology | Version (target) | Rationale |
|-------|------------|------------------|-----------|
| Language | Go | Current repo Go version | No change from existing tool. |
| Config loading | `cleanenv` + standard library | Already in use for `app.yaml`. |
| Data parsing | `encoding/json`, `gopkg.in/yaml.v3` (or similar) | Parse YAML/JSON for card records and schemas. |
| Templates | `html/template` | Existing HTML templates for cards. |
| Renderer | `chromedp`, `cdproto/page` | Existing HTML → PNG/PDF pipeline. |
| HTTP static server | `net/http` | Serve images/fonts for templates. |

**Dependencies (conceptual JSON)**

```json
{
  "dependencies": {
    "github.com/ilyakaznacheev/cleanenv": "existing",
    "github.com/chromedp/chromedp": "existing",
    "github.com/chromedp/cdproto/page": "existing",
    "gopkg.in/yaml.v3": "new (for YAML parsing)"
  }
}
```

---

## 3. Component Design

### 3.1 Project Configuration Loader

- **Purpose:** Load and validate project-level configuration (card types, paths, renderer settings).
- **Responsibilities:**
  - Extend existing `config.Config` to support projects and generic card types.
  - Load YAML/JSON config files from a known location (e.g., `configs/<project>.yaml`).
  - Validate presence of data/template directories and flags.
- **Interfaces:**
  - `LoadConfig() (Config, error)` – existing, extended.
  - New: `ProjectConfig` struct capturing:
    - `Name`, `CardTypes`, `DataDir`, `TemplateDir`, `OutputDir`, etc.

### 3.2 Card Type Registry & Schema

- **Purpose:** Represent each card type’s schema and make it available to other components.
- **Responsibilities:**
  - Keep a registry keyed by card type ID (`"weapon_frame"`, `"ingredient"`, etc.).
  - For each type, store:
    - Field definitions (name, type, required, default).
    - Data source specification (YAML or JSON file(s)).
    - Template file path(s).
  - Offer schema-based validation for card instances.
- **Interfaces:**
  - `type FieldSchema struct { Name string; Type string; Required bool; Default any }`
  - `type CardTypeSchema struct { ID string; Fields []FieldSchema; DataFiles []string; Template string }`
  - `type TypeRegistry interface { Get(id string) (CardTypeSchema, bool); List() []CardTypeSchema }`

### 3.3 Card Data Loader

- **Purpose:** Load card records from YAML/JSON files into a generic structure based on schemas.
- **Responsibilities:**
  - For each `CardTypeSchema`, open related YAML/JSON files from `DataDir`.
  - Decode into either:
    - `[]map[string]any` for maximum flexibility, or
    - A small generic struct `GenericCard` with `TypeID`, `ID`, `Fields map[string]any`.
  - Validate each record against the schema (missing required fields, unknown fields, type mismatches).
  - Optionally, adapt or convert existing weapon data into this structure.
- **Interfaces:**
  - `type GenericCard struct { TypeID string; ID string; Fields map[string]any }`
  - `LoadCards(schema CardTypeSchema, dataDir string) ([]GenericCard, error)`

### 3.4 Template Binding & View Model Builder

- **Purpose:** Prepare data for HTML templates in a way that keeps templates simple and expressive.
- **Responsibilities:**
  - Provide a consistent template context for all card types:
    - `Card`: `GenericCard`
    - Optional helpers: derived values (e.g., arrays for bullet icons, formatted text).
  - Optionally support type-specific view-model adaptors:
    - e.g., for weapon cards, convert generic fields into a struct implementing existing helper methods (`HasDamage`, etc.).
  - Define a clear contract for templates to avoid tight coupling to Go internals.
- **Interfaces:**
  - `type TemplateContext struct { Card GenericCard; Type CardTypeSchema; Helpers map[string]any }`
  - `BuildContext(card GenericCard, schema CardTypeSchema) TemplateContext`

### 3.5 Rendering Orchestrator

- **Purpose:** Drive the end-to-end flow: load config → load schemas → load cards → render HTML → invoke renderer.
- **Responsibilities:**
  - Replace most of the game-specific logic in `main.go` with a generic loop over card types.
  - For each card type:
    - Parse its HTML template using `html/template`.
    - For each card:
      - Build `TemplateContext`.
      - Execute template to produce HTML.
      - Optionally write raw HTML to disk (config flag).
      - Call `ChromeRendererImpl.RenderHTMLToPNG` / `RenderHTMLToPDF`.
  - Ensure errors are logged with sufficient context.
- **Interfaces:**
  - `RenderProject(cfg ProjectConfig, registry TypeRegistry, renderer ChromeRenderer) error`

### 3.6 CLI / Entry Point

- **Purpose:** Provide a simple way to run the engine for a given project.
- **Responsibilities:**
  - Parse command-line arguments / env vars:
    - `--project <name>` or `PROJECT_CONFIG=<path>`.
  - Use existing `app.yaml` for defaults (renderer URL, global output dir) and project config for card-specific settings.
- **Interfaces:**
  - `func main()` orchestrates: load config → select project → call `RenderProject`.

---

## 4. Data Model

### 4.1 Core Entities (Go)

```go
type FieldSchema struct {
    Name     string      `yaml:"name" json:"name"`
    Type     string      `yaml:"type" json:"type"` // "string", "int", "float", "bool", "string_list", "object", "object_list", etc.
    Required bool        `yaml:"required" json:"required"`
    Default  interface{} `yaml:"default,omitempty" json:"default,omitempty"`
}

type CardTypeSchema struct {
    ID        string        `yaml:"id" json:"id"`
    DataFiles []string      `yaml:"data_files" json:"data_files"`
    Template  string        `yaml:"template" json:"template"`
    Fields    []FieldSchema `yaml:"fields" json:"fields"`
}

type ProjectConfig struct {
    Name       string           `yaml:"name" json:"name"`
    DataDir    string           `yaml:"data_dir" json:"data_dir"`
    TemplateDir string          `yaml:"template_dir" json:"template_dir"`
    OutputDir  string           `yaml:"output_dir" json:"output_dir"`
    CardTypes  []CardTypeSchema `yaml:"card_types" json:"card_types"`
}

type GenericCard struct {
    TypeID string                 `yaml:"type_id" json:"type_id"`
    ID     string                 `yaml:"id" json:"id"`
    Fields map[string]interface{} `yaml:"fields" json:"fields"`
}
```

### 4.2 Example YAML Schemas

**Card type schema example (weapon frame):**

```yaml
id: weapon_frame
template: WeaponFrame.html
data_files:
  - weapon_frames.yaml
fields:
  - name: name
    type: string
    required: true
  - name: type
    type: string
    required: true
  - name: damage
    type: int
    required: false
  - name: tags
    type: string_list
    required: false
```

**Card record example (CookCook ingredient):**

```yaml
- type_id: ingredient
  id: coconut-milk
  fields:
    name: Coconut Milk
    category: Coconut Milk
    description: Rich coconut cream base for soups and curries.
    icon: coconut-milk
```

---

## 5. API Contracts

This tool is primarily a CLI application, not a network service. The “API” is:

### 5.1 Command-Line Interface

| Option | Description |
|--------|-------------|
| `--project <name>` | Select a project to render, corresponding to a project config entry. |
| `--config <path>` | Override path to main config file (defaults to `app.yaml`). |
| `--output-dir <path>` | Override output directory. |
| `--html` | Enable writing intermediate HTML files. |

Example:

```bash
card-maker --project cookcook --config configs/app.yaml --output-dir ./output/cookcook
```

### 5.2 Template Context

Templates will receive a context with at least:

```go
type TemplateContext struct {
    Card   GenericCard
    Schema CardTypeSchema
}
```

Within HTML templates, fields are accessed as:

```gotemplate
{{ .Card.Fields.name }}
{{ index .Card.Fields "damage" }}
```

Type-specific helpers (if needed) can be injected via `template.FuncMap`.

---

## 6. Security Considerations

- **Data sources are local files only;** no untrusted remote input is processed.
- **Sandboxing Chrome:** Ensure the headless Chrome invocation does not navigate to external URLs; all HTML is in-memory.
- **Template safety:** Continue to use `html/template` (not `text/template`) to benefit from HTML escaping.
- **File paths:** Validate and normalize configurable directories to avoid path traversal issues.

Checklist:

- [ ] No network access in templates or render pipeline beyond local static assets.
- [ ] No execution of untrusted code; templates are authored by the project owner.
- [ ] Limit/validate output paths to stay within configured `OutputDir`.

---

## 7. Performance Strategy

- **Batch rendering:** Reuse a Chrome context where possible for multiple cards to avoid high startup overhead.
- **Parallelism:** Optionally introduce parallel rendering per card type (configurable concurrency) while being mindful of Chrome’s resource usage.
- **Template reuse:** Parse each card type’s template once, reuse for all its cards.
- **I/O:** Stream data from YAML/JSON in manageable chunks; expected dataset sizes (hundreds of cards) are small enough for in-memory processing.

Targets:

- Rendering hundreds of cards should complete within a reasonable time on a standard dev machine (minutes, not hours).
- Memory footprint remains modest; no need for special optimization yet.

---

## 8. Implementation Phases

### Phase 1: Introduce Generic Types & Config (No Behavior Change)

- Add `ProjectConfig`, `CardTypeSchema`, `FieldSchema`, and `GenericCard` definitions.
- Extend `config.Config` / `app.yaml` to include a default project referencing existing CSV-based entities (still using current flow).
- Add a feature-flagged path in `main.go` that **reads but does not yet use** the new generic configuration.

### Phase 2: Implement YAML/JSON Loader & Registry

- Implement `TypeRegistry` and generic `LoadCards` for YAML/JSON.
- Create sample YAML/JSON for one existing type (e.g., `Effect`) to validate the loader.
- Add validation and error reporting according to the spec.

### Phase 3: Generic Rendering Orchestrator (Alongside Existing Flow)

- Implement `RenderProject` that:
  - Iterates over configured card types.
  - Parses templates.
  - Renders cards via `ChromeRendererImpl`.
- Wire a CLI flag or config option to choose between:
  - Legacy flow (CSV + typed entities).
  - New generic flow (YAML/JSON + `GenericCard`).

### Phase 4: Migrate Existing Weapon Project to Generic Engine

- Convert `Effects.csv`, `WeaponFrame.csv`, `WeaponPart.csv`, and `Item.csv` to YAML/JSON under a project data directory.
- Define corresponding `CardTypeSchema` entries and templates.
- Adjust templates to use `GenericCard` fields (or lightweight adaptors + `FuncMap`).
- Compare generated PNGs/PDFs against legacy outputs to check for regressions.

### Phase 5: Add CookCook Project

- Define a `cookcook` project configuration with:
  - `ingredient` and any other needed card types.
  - YAML/JSON data files representing ingredient cards.
  - Initial simple ingredient template (to be designed separately).
- Verify CookCook cards render successfully through the same engine.

### Phase 6: Cleanup & Hardening

- Remove or deprecate legacy, game-specific entities and mappers as appropriate.
- Improve logging and error messages.
- Add basic tests for:
  - Config loading.
  - YAML/JSON parsing and validation.
  - Template execution for at least one card type.

---

## 9. Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|-----------|
| Migration complexity from CSV to YAML/JSON | Medium | Medium | Automate conversion with scripts; start with a single entity type. |
| Template breakage when moving to generic context | Medium | Medium | Introduce adaptors or helper functions; migrate templates incrementally. |
| Over-generalization making templates harder to write | Medium | Low-Med | Keep `GenericCard` simple; allow type-specific helpers where needed. |
| Increased complexity in configuration | Low-Med | Medium | Provide clear examples and defaults; validate configs on startup with friendly errors. |
| Performance regressions due to new rendering loop | Low | Low | Reuse existing renderer logic; benchmark against current behavior. |

---

## 10. Open Questions

- How much convenience should be preserved for the original weapon project (e.g., keep some typed structs and helper methods) vs fully migrating to `GenericCard`?
- Should multiple templates per card type (variants/skins) be supported in the first iteration?
- Do we need per-project overrides for renderer viewport size and DPI, or is one global setting sufficient?
- How will ingredient-type cards for CookCook handle multilingual text, if at all?

---

## Next Steps

1. Confirm this plan aligns with the desired level of generality and the YAML/JSON direction.
2. Implement Phase 1 and Phase 2 behind a feature flag to avoid disturbing existing behavior.
3. Once stable, proceed with migrating the weapon project and adding the CookCook project configuration and templates.

