package models

import "mime"

type MimeType string

const (
	MimeTypeImagePNG        MimeType = "image/png"
	MimeTypeImageJPEG       MimeType = "image/jpeg"
	MimeTypeImageGIF        MimeType = "image/gif"
	MimeTypeImageWebP       MimeType = "image/webp"
	MimeTypeImageSVG        MimeType = "image/svg+xml"
	MimeTypeApplicationPDF  MimeType = "application/pdf"
	MimeTypeTextPlain       MimeType = "text/plain"
	MimeTypeTextMarkdown    MimeType = "text/markdown"
	MimeTypeApplicationJSON MimeType = "application/json"
)

func ParseMimeType(s string) MimeType {
	if mediaType, _, err := mime.ParseMediaType(s); err == nil {
		s = mediaType
	}
	switch s {
	case string(MimeTypeImagePNG):
		return MimeTypeImagePNG
	case string(MimeTypeImageJPEG):
		return MimeTypeImageJPEG
	case string(MimeTypeImageGIF):
		return MimeTypeImageGIF
	case string(MimeTypeImageWebP):
		return MimeTypeImageWebP
	case string(MimeTypeImageSVG):
		return MimeTypeImageSVG
	case string(MimeTypeApplicationPDF):
		return MimeTypeApplicationPDF
	case string(MimeTypeTextPlain):
		return MimeTypeTextPlain
	case string(MimeTypeTextMarkdown):
		return MimeTypeTextMarkdown
	case string(MimeTypeApplicationJSON):
		return MimeTypeApplicationJSON
	default:
		return ""
	}
}

