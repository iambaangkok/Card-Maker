package generic

import "github.com/iambaangkok/Card-Maker/internal/renderer"

// ResolveViewport returns the PNG screenshot clip size in CSS pixels for a card type.
// Precedence per axis: card type → project default → built-in renderer default (poker 240×336).
func ResolveViewport(project ProjectConfig, schema CardTypeSchema) (width, height float64) {
	width = float64(renderer.CardViewportWidth)
	height = float64(renderer.CardViewportHeight)
	if project.DefaultViewportWidth != nil && *project.DefaultViewportWidth > 0 {
		width = *project.DefaultViewportWidth
	}
	if project.DefaultViewportHeight != nil && *project.DefaultViewportHeight > 0 {
		height = *project.DefaultViewportHeight
	}
	if schema.ViewportWidth != nil && *schema.ViewportWidth > 0 {
		width = *schema.ViewportWidth
	}
	if schema.ViewportHeight != nil && *schema.ViewportHeight > 0 {
		height = *schema.ViewportHeight
	}
	return width, height
}

// ResolveOutputScale returns the Chromedp screenshot scale (PNG pixel multiplier vs viewport).
// Precedence: card type → project default → renderer.DefaultOutputScale.
func ResolveOutputScale(project ProjectConfig, schema CardTypeSchema) float64 {
	s := renderer.DefaultOutputScale
	if project.DefaultOutputScale != nil && *project.DefaultOutputScale > 0 {
		s = *project.DefaultOutputScale
	}
	if schema.OutputScale != nil && *schema.OutputScale > 0 {
		s = *schema.OutputScale
	}
	return s
}
