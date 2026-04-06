package generic

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
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

// buildTemplateFuncMap creates a FuncMap with refLookup for generic reference data injection.
// refLookup(refKey, lookupValue, keyField, returnField) looks up in refData[refKey], matches by keyField
// (parsing "X:Y" to get "X"), returns the returnField value. Templates specify which fields to use.
func buildTemplateFuncMap(refData map[string]interface{}) template.FuncMap {
	return template.FuncMap{
		// menuGradeStars returns n copies of ★ for menu grade (2–5 typical); clamps 1–5; "?" if unset.
		"menuGradeStars": func(v interface{}) string {
			n := 0
			switch x := v.(type) {
			case int:
				n = x
			case int64:
				n = int(x)
			case float64:
				n = int(x)
			}
			if n < 1 {
				return "?"
			}
			if n > 5 {
				n = 5
			}
			return strings.Repeat("★", n)
		},
		// menuBoldOR wraps spaced " OR " in <strong> for requirement/effect lines (trusted YAML).
		"menuBoldOR": func(v interface{}) template.HTML {
			s := fmt.Sprint(v)
			s = strings.ReplaceAll(s, " OR ", " <strong>OR</strong> ")
			return template.HTML(s)
		},
		"refLookup": func(refKey string, lookupValue interface{}, keyField, returnField string) interface{} {
			data, ok := refData[refKey]
			if !ok {
				return ""
			}
			list, ok := data.([]interface{})
			if !ok {
				return ""
			}
			lookupStr := fmt.Sprint(lookupValue)
			keyVal := lookupStr
			if idx := strings.Index(lookupStr, ":"); idx >= 0 {
				keyVal = lookupStr[:idx]
			}
			for _, e := range list {
				m, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				if fmt.Sprint(m[keyField]) == keyVal {
					if v, ok := m[returnField]; ok {
						return v
					}
					return ""
				}
			}
			return ""
		},
	}
}

// RenderProject renders all cards for the given project configuration using
// the provided registry and renderer. It expects template paths in the
// schemas to be relative to project.TemplateDir. Each card type writes into
// project.OutputDir/<schema_id>/ (e.g. output/cookcook/ingredient_3/).
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

	schemas := reg.List()

	for typeIdx, schema := range schemas {
		if schema.Render != nil && !*schema.Render {
			log.Printf("skip %q (render: false)", schema.ID)
			continue
		}

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

		totalCards := len(cards)
		typeOutDir := filepath.Join(project.OutputDir, schema.ID)
		if err := os.MkdirAll(typeOutDir, 0o755); err != nil {
			return fmt.Errorf("create output dir for type %q: %w", schema.ID, err)
		}

		vpW, vpH := ResolveViewport(project, schema)
		outScale := ResolveOutputScale(project, schema)
		log.Printf("[%d/%d] %q: %d cards → %q (viewport %.2f×%.2f, output_scale %.2f)", typeIdx+1, len(schemas), schema.ID, totalCards, typeOutDir, vpW, vpH, outScale)

		for cardIdx, card := range cards {
			log.Printf("  [%d/%d] %s", cardIdx+1, totalCards, card.ID)
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
				htmlPath := filepath.Join(typeOutDir, baseName+".html")
				if err := os.WriteFile(htmlPath, []byte(htmlStr), 0o644); err != nil {
					return fmt.Errorf("write html %q: %w", htmlPath, err)
				}
			}

			imgPath := filepath.Join(typeOutDir, baseName+".png")
			if err := r.RenderHTMLToPNG(htmlStr, imgPath, vpW, vpH, outScale); err != nil {
				return fmt.Errorf("render png for card %q of type %q: %w", card.ID, schema.ID, err)
			}
		}
	}

	return nil
}

