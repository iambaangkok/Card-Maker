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
    - `ReferenceData`: `map[string]interface{}` for lookup data (e.g. effects).
  - Optionally support type-specific view-model adaptors:
    - e.g., for weapon cards, convert generic fields into a struct implementing existing helper methods (`HasDamage`, etc.).
  - Define a clear contract for templates to avoid tight coupling to Go internals.
- **Interfaces:**
  - `type TemplateContext struct { Card GenericCard; Schema CardTypeSchema; ReferenceData map[string]interface{} }`
  - Reference data loading is **generic**: any key in `project.ReferenceData` maps to a YAML/JSON file; no hard-coded effect-specific logic in the loader.

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
    Name          string           `yaml:"name" json:"name"`
    DataDir       string           `yaml:"data_dir" json:"data_dir"`
    TemplateDir   string           `yaml:"template_dir" json:"template_dir"`
    ImageDir      string           `yaml:"image_dir" json:"image_dir"`
    OutputDir     string           `yaml:"output_dir" json:"output_dir"`
    CardTypes     []CardTypeSchema `yaml:"card_types" json:"card_types"`
    ReferenceData map[string]string `yaml:"reference_data,omitempty" json:"reference_data,omitempty"` // key -> file (generic)
}

type GenericCard struct {
    TypeID string                 `yaml:"type_id" json:"type_id"`
    ID     string                 `yaml:"id" json:"id"`
    Fields map[string]interface{} `yaml:"fields" json:"fields"`
}
```

### 4.2 Reference Data (Non-Card Types)

Reference data is loaded once per project and passed to all templates. It is **generic**: any key can map to any YAML/JSON file. Common use: effect definitions, tag glossaries, etc.

```go
// ProjectConfig
ReferenceData map[string]string  // key -> file path (e.g. "effects" -> "effects.yaml")

// TemplateContext
ReferenceData map[string]interface{}  // key -> loaded data (e.g. "effects" -> []map[string]interface{})
```

**Config example:**
```yaml
reference_data:
  effects: "effects.yaml"
  # Future: tags: "tags.yaml", abilities: "abilities.yaml"
```

**effects.yaml structure** (each entry: `name`, `type`, `has_level`, `description`):
```yaml
- name: Nimble
  type: Passive
  has_level: true
  description: "Movement +1 per level."
```

**Template usage:** Templates use the generic `refLookup` FuncMap helper and specify which fields to use:

```gotemplate
{{with refLookup "effects" . "name" "description"}}{{.}}{{end}}
```

`refLookup(refKey, lookupValue, keyField, returnField)` — looks up in `ReferenceData[refKey]` (a list of maps), matches by `keyField` (parsing `"X:Y"` to get `X`), returns the `returnField` value. The template chooses key and return fields (e.g. `"name"`/`"description"`, `"id"`/`"label"`).

### 4.3 Example YAML Schemas

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

Templates receive:

```go
type TemplateContext struct {
    Card          GenericCard
    Schema        CardTypeSchema
    ReferenceData map[string]interface{}  // e.g. "effects" -> []map[string]interface{}
}
```

**Field access:**
```gotemplate
{{ .Card.Fields.name }}
{{ index .Card.Fields "damage" }}
```

**Reference data lookup:** Use the generic `refLookup(refKey, lookupValue, keyField, returnField)` helper:

```gotemplate
{{range .Card.Fields.effects}}
  <div class="Effect">...</div>
  {{with refLookup "effects" . "name" "description"}}
  <div class="Description">{{.}}</div>
  {{end}}
{{end}}
```

The template specifies `keyField` (field to match, e.g. `"name"`) and `returnField` (field to return, e.g. `"description"` or `"type"`). Parses `"Name:Level"` to extract key for matching.

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

- [x] Convert `Effects.csv` to `effects.yaml` as **reference data** (not a card type).
- [x] Add `reference_data` to `ProjectConfig` and `LoadReferenceData` (generic: any key → any YAML/JSON file).
- [x] Pass `ReferenceData` into `TemplateContext` for all card types.
- [x] Convert `WeaponFrame.csv`, `WeaponPart.csv`, and `Item.csv` to YAML under `data/example/`.
- [x] Define `CardTypeSchema` entries for `weapon_frame`, `weapon_part`, `item`.
- [x] Adjust templates to use `GenericCard` fields.
- [ ] **Port effects usage to HTML templates:** Use `.ReferenceData.effects` in `WeaponPart.html` and `Item.html` to optionally display effect type badges or tooltips (e.g. `title="{{.description}}"`). Add a `template.FuncMap` helper (e.g. `effectLookup`) to resolve `"Reflect:1"` → effect name + level when matching against reference data.
- [ ] **Port all items and weapon parts data:** Complete migration so `weapon_parts.yaml` has all 25 entries from `WeaponPart.csv` and `items.yaml` has all 13 entries from `Item.csv` (currently 6 and 7 respectively).
- [ ] Compare generated PNGs/PDFs against legacy outputs to check for regressions.

**Phase 4a: Effects Usage in Templates**

- Ensure `LoadReferenceData` remains generic (no effect-specific code; any `reference_data` key loads any YAML/JSON).
- Use generic `refLookup(refKey, lookupValue)` FuncMap helper for any reference data key.
- Update `WeaponPart.html` and `Item.html` to use `.ReferenceData.effects` for effect display (e.g. tooltip with description, or type badge).
- Templates must handle effects with level suffix (`Name:Level`) when matching reference data.

**Phase 4b: Complete Data Migration**

| Source | Target | Current | Target |
|--------|--------|--------|--------|
| `WeaponPart.csv` | `weapon_parts.yaml` | 6 entries | 25 entries |
| `Item.csv` | `items.yaml` | 7 entries | 13 entries |

**WeaponPart.csv columns:** Name, Manufacturer, Type, Damage, FireRate, Accuracy, MinRange, MaxRange, AmmoPerMag, Price, Compatibles, Tags, Effects.

**Item.csv columns:** Name, UsageLimit, Price, Effects.

Convert compatibles/tags from `Pistol/SMG/AR` to `[Pistol, SMG, AR]`. Convert effects from `Nimble/Handling:2` to `[Nimble, Handling:2]`. Add `usage_boxes: [1, 2, ...]` for items with `usage_limit > 0`.

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

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.1 | 2026-03-16 | [REFINED] Generalized reference data injection: replaced `effectLookup` with `refLookup(refKey, lookupValue)` for any reference_data key |
| 1.2 | 2026-03-16 | [REFINED] `refLookup` now accepts `keyField` and `returnField` — templates specify which field to match and which to return |

---

## Next Steps

1. Implement Phase 4a (effects usage in templates) and Phase 4b (complete data migration).
2. Run `/implement Card-maker-generic-refactor` to execute remaining tasks.
3. Compare generated PNGs against legacy outputs for regression validation.
4. Proceed with CookCook project (Phase 5) once weapon project migration is complete.

