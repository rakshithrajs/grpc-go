package handlers

import (
	"context"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// logPrefix returns a formatted string for logging purposes, including the function name.
func logPrefix(fn string) string { return "[" + fn + "]: " }

var protoMimeTypeMap = map[models.MimeType]MMSpb.MimeType{
	models.MimeTypeImagePNG:      MMSpb.MimeType_MIME_TYPE_IMAGE_PNG,
	models.MimeTypeImageJPEG:     MMSpb.MimeType_MIME_TYPE_IMAGE_JPEG,
	models.MimeTypeImageGIF:      MMSpb.MimeType_MIME_TYPE_IMAGE_GIF,
	models.MimeTypeImageWebP:     MMSpb.MimeType_MIME_TYPE_IMAGE_WEBP,
	models.MimeTypeImageSVG:      MMSpb.MimeType_MIME_TYPE_IMAGE_SVG,
	models.MimeTypeApplicationPDF: MMSpb.MimeType_MIME_TYPE_APPLICATION_PDF,
	models.MimeTypeTextPlain:     MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
	models.MimeTypeTextMarkdown:  MMSpb.MimeType_MIME_TYPE_TEXT_MARKDOWN,
	models.MimeTypeApplicationJSON: MMSpb.MimeType_MIME_TYPE_APPLICATION_JSON,
}

// toProtoMimeType converts a models.MimeType to its corresponding MMSpb.MimeType.
func toProtoMimeType(mt models.MimeType) MMSpb.MimeType {
	if pt, ok := protoMimeTypeMap[mt]; ok {
		return pt
	}
	return MMSpb.MimeType_MIME_TYPE_UNSPECIFIED
}

// UserIDFromContext extracts the user ID from the gRPC context metadata.
func UserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return config.NullString, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	}

	userIDs := md.Get(config.UserIDMetadataKey)
	if len(userIDs) == 0 || userIDs[0] == config.NullString {
		return config.NullString, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	}

	return userIDs[0], nil
}
