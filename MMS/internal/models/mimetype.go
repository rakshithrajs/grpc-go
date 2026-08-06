package models

import (
	"mime"

	"github.com/rakshithrajs/cloud/MMS/internal/config"
)

// MimeType is a supported file MIME type.
type MimeType string

const (
	// MimeTypeImagePNG is the MIME type for PNG images.
	MimeTypeImagePNG MimeType = "image/png"
	// MimeTypeImageJPEG is the MIME type for JPEG images.
	MimeTypeImageJPEG MimeType = "image/jpeg"
	// MimeTypeImageGIF is the MIME type for GIF images.
	MimeTypeImageGIF MimeType = "image/gif"
	// MimeTypeImageWebP is the MIME type for WebP images.
	MimeTypeImageWebP MimeType = "image/webp"
	// MimeTypeImageSVG is the MIME type for SVG images.
	MimeTypeImageSVG MimeType = "image/svg+xml"
	// MimeTypeApplicationPDF is the MIME type for PDF documents.
	MimeTypeApplicationPDF MimeType = "application/pdf"
	// MimeTypeTextPlain is the MIME type for plain text files.
	MimeTypeTextPlain MimeType = "text/plain"
	// MimeTypeTextMarkdown is the MIME type for Markdown text files.
	MimeTypeTextMarkdown MimeType = "text/markdown"
	// MimeTypeApplicationJSON is the MIME type for JSON files.
	MimeTypeApplicationJSON MimeType = "application/json"
)

var mimeTypeMap = map[string]MimeType{
	string(MimeTypeImagePNG):      MimeTypeImagePNG,
	string(MimeTypeImageJPEG):     MimeTypeImageJPEG,
	string(MimeTypeImageGIF):      MimeTypeImageGIF,
	string(MimeTypeImageWebP):     MimeTypeImageWebP,
	string(MimeTypeImageSVG):      MimeTypeImageSVG,
	string(MimeTypeApplicationPDF): MimeTypeApplicationPDF,
	string(MimeTypeTextPlain):     MimeTypeTextPlain,
	string(MimeTypeTextMarkdown):  MimeTypeTextMarkdown,
	string(MimeTypeApplicationJSON): MimeTypeApplicationJSON,
}

// ParseMimeType takes a string representation of a MIME type and returns the corresponding MimeType constant.
func ParseMimeType(s string) MimeType {
	if mediaType, _, err := mime.ParseMediaType(s); err == nil {
		s = mediaType
	}
	if mt, ok := mimeTypeMap[s]; ok {
		return mt
	}
	return config.NullString
}
