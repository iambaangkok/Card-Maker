package generic

// FieldSchema describes a single field on a card type.
type FieldSchema struct {
	Name     string      `yaml:"name" json:"name"`
	Type     string      `yaml:"type" json:"type"` // e.g. "string", "int", "bool", "string_list"
	Required bool        `yaml:"required" json:"required"`
	Default  interface{} `yaml:"default,omitempty" json:"default,omitempty"`
}

// CardTypeSchema defines a card type, its data source files, and template.
type CardTypeSchema struct {
	ID        string        `yaml:"id" json:"id"`
	DataFiles []string      `yaml:"data_files" json:"data_files"`
	Template  string        `yaml:"template" json:"template"`
	Fields    []FieldSchema `yaml:"fields" json:"fields"`
	Render    *bool         `yaml:"render,omitempty" json:"render,omitempty"` // if nil or true, render; if false, skip
	// ViewportWidth/Height: optional PNG clip in CSS pixels (must match root card box in the template).
	// If unset, project defaults apply; if those are unset, renderer built-in defaults apply.
	ViewportWidth  *float64 `yaml:"viewport_width,omitempty" json:"viewport_width,omitempty"`
	ViewportHeight *float64 `yaml:"viewport_height,omitempty" json:"viewport_height,omitempty"`
	// OutputScale: PNG resolution multiplier vs CSS viewport (Chromedp Scale). E.g. 2 → 480×672 px from 240×336 layout.
	OutputScale *float64 `yaml:"output_scale,omitempty" json:"output_scale,omitempty"`
}

// ProjectConfig groups card types and paths for a single project/game.
type ProjectConfig struct {
	Name          string           `yaml:"name" json:"name"`
	DataDir       string           `yaml:"data_dir" json:"data_dir"`
	TemplateDir   string           `yaml:"template_dir" json:"template_dir"`
	ImageDir      string           `yaml:"image_dir" json:"image_dir"`
	OutputDir     string           `yaml:"output_dir" json:"output_dir"`
	CardTypes     []CardTypeSchema `yaml:"card_types" json:"card_types"`
	ReferenceData map[string]string `yaml:"reference_data,omitempty" json:"reference_data,omitempty"` // key -> data file (e.g. "effects" -> "effects.yaml")
	// Default viewport for all card types unless overridden per schema (CSS px; match template outer size).
	DefaultViewportWidth  *float64 `yaml:"default_viewport_width,omitempty" json:"default_viewport_width,omitempty"`
	DefaultViewportHeight *float64 `yaml:"default_viewport_height,omitempty" json:"default_viewport_height,omitempty"`
	// DefaultOutputScale: project-wide PNG scale unless a card type sets output_scale.
	DefaultOutputScale *float64 `yaml:"default_output_scale,omitempty" json:"default_output_scale,omitempty"`
}

// GenericCard is a generic representation of a card instance.
type GenericCard struct {
	TypeID string                 `yaml:"type_id" json:"type_id"`
	ID     string                 `yaml:"id" json:"id"`
	Fields map[string]interface{} `yaml:"fields" json:"fields"`
}

// TypeRegistry allows lookup of card type schemas by ID.
type TypeRegistry interface {
	Get(id string) (CardTypeSchema, bool)
	List() []CardTypeSchema
}

// InMemoryTypeRegistry is a simple in-memory implementation of TypeRegistry.
type InMemoryTypeRegistry struct {
	types map[string]CardTypeSchema
}

// NewInMemoryTypeRegistry constructs a registry from the given schemas.
func NewInMemoryTypeRegistry(schemas []CardTypeSchema) *InMemoryTypeRegistry {
	m := make(map[string]CardTypeSchema, len(schemas))
	for _, s := range schemas {
		m[s.ID] = s
	}
	return &InMemoryTypeRegistry{types: m}
}

// Get returns the schema for the given id, if present.
func (r *InMemoryTypeRegistry) Get(id string) (CardTypeSchema, bool) {
	s, ok := r.types[id]
	return s, ok
}

// List returns all registered schemas.
func (r *InMemoryTypeRegistry) List() []CardTypeSchema {
	out := make([]CardTypeSchema, 0, len(r.types))
	for _, s := range r.types {
		out = append(out, s)
	}
	return out
}

