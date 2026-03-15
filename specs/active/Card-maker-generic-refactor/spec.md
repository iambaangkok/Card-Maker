# Specification: Generic, Data-Driven Card Maker

**Task ID:** cook-food-cardgame-card-maker-refactor  
**Created:** 2026-03-15  
**Status:** Ready for Planning  
**Version:** 1.0

## 1. Problem Statement

- **The Problem:** The current `Card-Maker` tool is tightly coupled to a specific prototype boardgame (weapon frames, parts, and items). Card entities, mappers, CSV schemas, templates, and image paths are all baked around that game’s concepts, making it hard to reuse the tool for other games like CookCook.
- **Current Situation:** The Go application reads CSVs for `WeaponFrame`, `WeaponPart`, and `Item`, maps them to strongly-typed domain structs with game-specific stats (damage, range, accuracy, etc.), injects them into HTML templates, and renders PNGs/PDFs via a Chrome-based renderer. Adding a new game or card type requires new structs, mappers, templates, and often logic changes in `main.go`.
- **Desired Outcome:** A generic, data-driven card generator that lets a designer define card types, fields, templates, and assets via configuration (CSV/JSON/YAML and HTML templates) without changing core Go code. The tool should remain game-agnostic while supporting multiple projects; CookCook’s ingredient-type cards will be the first real consumer.

## 2. User Personas

### Primary User: Game Designer / System Designer
- **Who:** A designer creating analog or hybrid board/card games who is comfortable editing structured data (CSV/JSON/YAML) and simple HTML templates, but does not want to modify Go code.
- **Goals:**
  - Define card types (e.g., weapon, item, ingredient, menu) and their fields.
  - Configure how cards look (templates, fonts, icons) and where data fields appear.
  - Generate print-ready card images (PNG) and/or PDFs from those definitions.
- **Pain points:**
  - Currently must work within hard-coded weapon/part/item concepts.
  - Adding a new game or card concept means new Go structs, mapping code, and templates.

### Secondary User: Developer / Toolsmith
- **Who:** A developer comfortable with Go who maintains the tool.
- **Goals:**
  - Keep the core engine small, testable, and extensible.
  - Onboard new games with minimal changes (ideally config + templates only).
- **Pain points:**
  - Domain-specific entities and mappers are scattered through the codebase.
  - Renderer and templates are aware of weapon-specific layout and assumptions.

## 3. Functional Requirements

### FR-1: Generic Card Type Definition
**Description:** The system must support defining multiple card “types” and their fields via configuration, without adding new Go structs per type.

**User Story:**
> As a game designer, I want to define card types (e.g., `weapon_frame`, `weapon_part`, `item`, `ingredient`, `menu`) and their fields in configuration so that I can reuse the same card engine for different games without touching Go code.

**Acceptance Criteria:**
- [ ] There is a configuration format (YAML/JSON) that allows declaring:
  - Card type IDs (e.g., `weapon_frame`, `ingredient`).
  - For each type, a set of fields (name, type, optional/required, default), including support for nested object fields and lists of objects.
  - For each type, the data source file name(s) or data source.
- [ ] Existing weapon-related cards can be expressed using this configuration model.
- [ ] New card types can be added by editing configuration and CSVs only.

**Priority:** Must Have

### FR-2: Data-Driven Card Records from YAML/JSON
**Description:** Card instances are loaded from YAML or JSON files according to the configured schema for each card type.

**User Story:**
> As a game designer, I want to define card instances in YAML or JSON files, following a schema for each card type, so that I can manage card content in structured text and regenerate cards easily.

**Acceptance Criteria:**
- [ ] For each configured card type, the system can load one or more YAML or JSON files from a configurable directory.
- [ ] YAML/JSON properties are mapped to fields defined in the card type schema (with sensible validation and error reporting for missing/extra properties).
- [ ] Non-scalar fields (e.g., lists, tags, effects) are represented using native YAML/JSON structures (arrays, objects) rather than delimiter-separated strings.
- [ ] Existing `Effects`, `WeaponFrame`, `WeaponPart`, and `Item` data can be migrated from CSV into YAML/JSON and loaded through the generic mechanism.

**Priority:** Must Have

### FR-3: Template-Based Rendering Per Card Type
**Description:** Each card type is associated with one or more HTML templates that define its layout, using the existing Go HTML template engine and static assets.

**User Story:**
> As a game designer, I want to assign an HTML template to each card type so that I can fully control the visual layout of each family of cards without changing Go code.

**Acceptance Criteria:**
- [ ] Configuration maps each card type to an HTML template file path.
- [ ] Templates can access card fields in a generic way (e.g., `.Fields["name"]`) and, optionally, type-specific helper methods.
- [ ] Static assets (icons, background images, fonts) continue to be served via HTTP and are usable in templates.
- [ ] Existing HTML templates (`WeaponFrame.html`, `WeaponPart.html`, `Item.html`) can be adapted to work with the generic model with minimal changes.

**Priority:** Must Have

### FR-4: Renderer Pipeline (HTML → PNG/PDF)
**Description:** The tool must continue to render configured card templates into PNG images (and optionally PDFs) using the existing Chrome-based renderer.

**User Story:**
> As a game designer, I want to generate PNG (and optionally PDF) files for all card instances so that I can print them or use them in digital prototypes.

**Acceptance Criteria:**
- [ ] For each card instance and its associated template, the renderer can produce a PNG file in a configurable output directory.
- [ ] The existing renderer (`ChromeRendererImpl`) remains in use with only minimal changes (e.g., configuration wiring).
- [ ] Optional: the same pipeline can output PDFs via the existing `RenderHTMLToPDF` without changing template semantics.
- [ ] Output file naming is configurable and can include card type and card identifier (e.g., `ingredient_chili-paste.png`).

**Priority:** Must Have

### FR-5: Project / Game Configuration
**Description:** The engine must support multiple “projects” or “games” (e.g., the current weapon game and CookCook) using different configurations and assets.

**User Story:**
> As a game designer, I want to configure different projects (like the weapon prototype and CookCook) so that I can switch between them without modifying the core engine.

**Acceptance Criteria:**
- [ ] There is a project-level configuration file (or flag) indicating:
  - CSV input directory.
  - Card type definitions to load.
  - Template paths.
  - Renderer output directory and flags (e.g., which card types to render).
- [ ] Switching projects is possible via config (e.g., different `app.yaml` or project config files) without recompiling or editing Go code.
- [ ] CookCook’s ingredient cards can be introduced as a new project that reuses the same engine.

**Priority:** Should Have

### FR-6: Validation and Error Reporting
**Description:** The tool should clearly report configuration and data issues instead of failing silently or panicking in the middle of rendering.

**User Story:**
> As a game designer, I want clear validation errors when my config or CSVs are wrong so that I can fix issues quickly without reading Go stack traces.

**Acceptance Criteria:**
- [ ] On startup, configuration is validated (missing templates, unknown card types, invalid field types).
- [ ] On CSV load, the tool reports missing/extra columns, invalid data types, and constraint violations per row.
- [ ] Errors include enough context (file, row, card type, field name) to be actionable.

**Priority:** Should Have

## 4. Non-Functional Requirements

- **Technology Stack:** Must remain in Go, using the existing Go `html/template` engine and the current HTML → PNG/PDF workflow via `chromedp`. No new major runtime dependencies beyond what is already in use unless absolutely necessary.
- **Configurability:** All game- and card-specific details must move into configuration and templates; core Go code should not encode specific game rules or stats.
- **Extensibility:** Adding a new card type or project should require only:
  - Updating configuration and CSVs.
  - Adding HTML templates and assets.
  - Minimal or no changes to engine code.
- **Performance:** Reasonable performance when rendering hundreds of cards in a batch; no need for real-time performance.
- **Maintainability:** Core engine code should be organized to separate:
  - Configuration loading.
  - Data loading and validation.
  - Rendering orchestration.

## 5. Out of Scope

- ❌ Implementing or simulating game rules, turn structure, or scoring logic inside `Card-Maker`. Rules stay external (e.g., in CookCook’s markdown rulebooks).
- ❌ Providing a full-fledged GUI editor for templates or card data; this spec focuses on a CLI/tool-driven flow using files.
- ❌ Online distribution, multiplayer support, or live card editing in a browser.
- ❌ Any CookCook-specific visual design decisions beyond acknowledging that CookCook will be the first consumer.

## 6. Edge Cases & Error Handling

| Scenario | Expected Behavior |
|----------|-------------------|
| CSV row missing a required field | Row is rejected with a clear error message indicating card type, row index, and missing field. |
| CSV contains unknown columns not present in schema | Tool logs a warning or error depending on configuration (fail-fast vs tolerant). |
| Template refers to a field that is not defined | Rendering fails for that card with a clear error stating template name and missing field. |
| Asset (image/font) referenced in template is missing | Renderer logs a warning and continues if possible; behavior is documented. |
| Multiple projects share the same card type name | Namespaces or separate configuration files avoid collisions; this is described in project config conventions. |

| Error | User Message | System Action |
|-------|--------------|---------------|
| Invalid configuration file | "Invalid config: [details]" | Abort startup, exit non-zero. |
| CSV parsing error | "Error parsing [file]: [details]" | Skip bad row or abort (configurable), log details. |
| Renderer cannot reach Chrome / fails | "Rendering failed for [card id]: [error]" | Abort or retry strategy defined; no silent failures. |

## 7. Success Metrics

| Metric | Target | How to Measure |
|--------|--------|----------------|
| Game-agnostic engine | 2+ distinct projects (weapon prototype + CookCook) using same core | Both projects can be configured and rendered via the same binary. |
| New card type onboarding effort | No Go code changes required | Ability to add a new card type using only config/templates/CSVs. |
| Migration completeness | All existing `WeaponFrame`, `WeaponPart`, `Item` cards supported | Side-by-side comparison of outputs before/after refactor. |
| Stability | No runtime panics on valid configs | Run regression renders with existing CSVs and templates. |

## 8. Open Questions

- [ ] How expressive should the schema for fields be (e.g., simple scalar/list types vs richer types like nested objects or computed fields)?
- [ ] Should template helpers remain hard-coded in Go (e.g., `HasRange`, `GetDamageStr`) or move to more generic helpers or data pre-processing?
- [ ] How should multiple projects be organized on disk (separate directories with their own `app.yaml`, or one global file listing projects)?
- [ ] Do we need explicit support for localization (multiple languages) at this stage?

## 9. Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-03-15 | Initial specification for generic, data-driven Card-Maker refactor. |

## Next Steps

1. Review this specification and confirm the generic scope and priorities.
2. Design a technical plan (`plan.md`) for how to evolve current entities/mappers into a schema-driven engine while preserving the existing weapon-based project.
3. Define the initial CookCook project configuration and simple ingredient-type card schema (card template details to be specified later).
4. Implement the refactor incrementally, validating that existing outputs remain stable while enabling new projects.

*Specification created with SDD 5.0*

