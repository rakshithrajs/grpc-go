package handlers

import (
	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
)

func toProtoMimeType(mt models.MimeType) MMSpb.MimeType {
	switch mt {
	case models.MimeTypeImagePNG:
		return MMSpb.MimeType_MIME_TYPE_IMAGE_PNG
	case models.MimeTypeImageJPEG:
		return MMSpb.MimeType_MIME_TYPE_IMAGE_JPEG
	case models.MimeTypeImageGIF:
		return MMSpb.MimeType_MIME_TYPE_IMAGE_GIF
	case models.MimeTypeImageWebP:
		return MMSpb.MimeType_MIME_TYPE_IMAGE_WEBP
	case models.MimeTypeImageSVG:
		return MMSpb.MimeType_MIME_TYPE_IMAGE_SVG
	case models.MimeTypeApplicationPDF:
		return MMSpb.MimeType_MIME_TYPE_APPLICATION_PDF
	case models.MimeTypeTextPlain:
		return MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN
	case models.MimeTypeTextMarkdown:
		return MMSpb.MimeType_MIME_TYPE_TEXT_MARKDOWN
	case models.MimeTypeApplicationJSON:
		return MMSpb.MimeType_MIME_TYPE_APPLICATION_JSON
	default:
		return MMSpb.MimeType_MIME_TYPE_UNSPECIFIED
	}
}
