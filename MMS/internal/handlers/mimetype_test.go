package handlers

import (
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
)

func TestToProtoMimeType(t *testing.T) {
	tests := []struct {
		name     string
		input    models.MimeType
		expected MMSpb.MimeType
	}{
		{
			name:     "image/png",
			input:    models.MimeTypeImagePNG,
			expected: MMSpb.MimeType_MIME_TYPE_IMAGE_PNG,
		},
		{
			name:     "image/jpeg",
			input:    models.MimeTypeImageJPEG,
			expected: MMSpb.MimeType_MIME_TYPE_IMAGE_JPEG,
		},
		{
			name:     "image/gif",
			input:    models.MimeTypeImageGIF,
			expected: MMSpb.MimeType_MIME_TYPE_IMAGE_GIF,
		},
		{
			name:     "image/webp",
			input:    models.MimeTypeImageWebP,
			expected: MMSpb.MimeType_MIME_TYPE_IMAGE_WEBP,
		},
		{
			name:     "image/svg+xml",
			input:    models.MimeTypeImageSVG,
			expected: MMSpb.MimeType_MIME_TYPE_IMAGE_SVG,
		},
		{
			name:     "application/pdf",
			input:    models.MimeTypeApplicationPDF,
			expected: MMSpb.MimeType_MIME_TYPE_APPLICATION_PDF,
		},
		{
			name:     "text/plain",
			input:    models.MimeTypeTextPlain,
			expected: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
		},
		{
			name:     "text/markdown",
			input:    models.MimeTypeTextMarkdown,
			expected: MMSpb.MimeType_MIME_TYPE_TEXT_MARKDOWN,
		},
		{
			name:     "application/json",
			input:    models.MimeTypeApplicationJSON,
			expected: MMSpb.MimeType_MIME_TYPE_APPLICATION_JSON,
		},
		{
			name:     "unknown mime type falls back to unspecified",
			input:    models.MimeType(config.NullString),
			expected: MMSpb.MimeType_MIME_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toProtoMimeType(tt.input)
			if got != tt.expected {
				t.Errorf("toProtoMimeType(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
