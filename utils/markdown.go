package utils

import (
	"bytes"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func MarkdownToHTML(md string) string {
	var buf bytes.Buffer

	goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub flavored markdown
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"), // or "monokai", "github", etc.
			),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // allow raw HTML in markdown
		),
	).Convert([]byte(md), &buf)

	return buf.String()
}
