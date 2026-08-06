package handlers

import MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"

var mimeTypeStringMap = map[MMSpb.MimeType]string{
	MMSpb.MimeType_MIME_TYPE_IMAGE_PNG:      "image/png",
	MMSpb.MimeType_MIME_TYPE_IMAGE_JPEG:     "image/jpeg",
	MMSpb.MimeType_MIME_TYPE_IMAGE_GIF:      "image/gif",
	MMSpb.MimeType_MIME_TYPE_IMAGE_WEBP:     "image/webp",
	MMSpb.MimeType_MIME_TYPE_IMAGE_SVG:      "image/svg+xml",
	MMSpb.MimeType_MIME_TYPE_APPLICATION_PDF: "application/pdf",
	MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN:     "text/plain",
	MMSpb.MimeType_MIME_TYPE_TEXT_MARKDOWN:  "text/markdown",
	MMSpb.MimeType_MIME_TYPE_APPLICATION_JSON: "application/json",
}

// MimeTypeToString converts a given MMSpb.MimeType to its corresponding string representation.
func MimeTypeToString(mt MMSpb.MimeType) string {
	if s, ok := mimeTypeStringMap[mt]; ok {
		return s
	}
	return "application/octet-stream"
}
