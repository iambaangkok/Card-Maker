package generic

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/iambaangkok/Card-Maker/internal/renderer"
)

// TemplateContext is the data passed into HTML templates for generic cards.
type TemplateContext struct {
	Card   GenericCard
	Schema CardTypeSchema
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

	for _, schema := range reg.List() {
		templatePath := filepath.Join(project.TemplateDir, schema.Template)
		t, err := template.ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("parse template %q: %w", templatePath, err)
		}

		cards, err := LoadCards(schema, project.DataDir)
		if err != nil {
			return fmt.Errorf("load cards for type %q: %w", schema.ID, err)
		}

		for _, card := range cards {
			ctx := TemplateContext{
				Card:   card,
				Schema: schema,
			}

			var buf bytes.Buffer
			if err := t.Execute(&buf, ctx); err != nil {
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

