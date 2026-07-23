package utils

import MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"

func MimeTypeToString(mt MMSpb.MimeType) string {
	switch mt {
	case MMSpb.MimeType_MIME_TYPE_IMAGE_PNG:
		return "image/png"
	case MMSpb.MimeType_MIME_TYPE_IMAGE_JPEG:
		return "image/jpeg"
	case MMSpb.MimeType_MIME_TYPE_IMAGE_GIF:
		return "image/gif"
	case MMSpb.MimeType_MIME_TYPE_IMAGE_WEBP:
		return "image/webp"
	case MMSpb.MimeType_MIME_TYPE_IMAGE_SVG:
		return "image/svg+xml"
	case MMSpb.MimeType_MIME_TYPE_APPLICATION_PDF:
		return "application/pdf"
	case MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN:
		return "text/plain"
	case MMSpb.MimeType_MIME_TYPE_TEXT_MARKDOWN:
		return "text/markdown"
	case MMSpb.MimeType_MIME_TYPE_APPLICATION_JSON:
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
