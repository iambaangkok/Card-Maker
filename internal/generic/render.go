package generic

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/iambaangkok/Card-Maker/internal/renderer"
)

// TemplateContext is the data passed into HTML templates for generic cards.
type TemplateContext struct {
	Card          GenericCard
	Schema        CardTypeSchema
	ReferenceData map[string]interface{} // e.g. "effects" -> []map[string]interface{}
}

// buildTemplateFuncMap creates a FuncMap with effectLookup that uses refData.
// effectLookup(effectStr) returns the effect description for tooltips; parses "Name:Level" to match by name.
func buildTemplateFuncMap(refData map[string]interface{}) template.FuncMap {
	effectMap := make(map[string]string)
	if effects, ok := refData["effects"]; ok {
		if list, ok := effects.([]interface{}); ok {
			for _, e := range list {
				if m, ok := e.(map[string]interface{}); ok {
					if name, ok := m["name"].(string); ok {
						if desc, ok := m["description"].(string); ok {
							effectMap[name] = desc
						}
					}
				}
			}
		}
	}
	return template.FuncMap{
		"effectLookup": func(v interface{}) string {
			effectStr := fmt.Sprint(v)
			name := effectStr
			if idx := strings.Index(effectStr, ":"); idx >= 0 {
				name = effectStr[:idx]
			}
			return effectMap[name]
		},
	}
}

// RenderProject renders all cards for the given project configuration using
// the provided registry and renderer. It expects template paths in the
// schemas to be relative to project.TemplateDir, and writes outputs into
// project.OutputDir.
//
// If writeHTML is true, the parsed HTML for each card is written alongside
// the rendered image.
func RenderProject(project ProjectConfig, reg TypeRegistry, r renderer.ChromeRendererImpl, writeHTML bool) error {
	if err := os.MkdirAll(project.OutputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	refData, err := LoadReferenceData(project)
	if err != nil {
		return fmt.Errorf("load reference data: %w", err)
	}

	fnMap := buildTemplateFuncMap(refData)

	for _, schema := range reg.List() {
		templatePath := filepath.Join(project.TemplateDir, schema.Template)
		t, err := template.New(filepath.Base(templatePath)).Funcs(fnMap).ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("parse template %q: %w", templatePath, err)
		}

		templateName := filepath.Base(templatePath)

		cards, err := LoadCards(schema, project.DataDir)
		if err != nil {
			return fmt.Errorf("load cards for type %q: %w", schema.ID, err)
		}

		for _, card := range cards {
			ctx := TemplateContext{
				Card:          card,
				Schema:        schema,
				ReferenceData: refData,
			}

			var buf bytes.Buffer
			if err := t.ExecuteTemplate(&buf, templateName, ctx); err != nil {
				return fmt.Errorf("execute template for card %q of type %q: %w", card.ID, schema.ID, err)
			}
			htmlStr := buf.String()

			baseName := fmt.Sprintf("%s_%s", schema.ID, card.ID)

			if writeHTML {
				htmlPath := filepath.Join(project.OutputDir, baseName+".html")
				if err := os.WriteFile(htmlPath, []byte(htmlStr), 0o644); err != nil {
					return fmt.Errorf("write html %q: %w", htmlPath, err)
				}
			}

			imgPath := filepath.Join(project.OutputDir, baseName+".png")
			if err := r.RenderHTMLToPNG(htmlStr, imgPath); err != nil {
				return fmt.Errorf("render png for card %q of type %q: %w", card.ID, schema.ID, err)
			}
		}
	}

	return nil
}

