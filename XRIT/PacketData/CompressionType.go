package PacketData

import "fmt"

// CompressionType
const (
	NO_COMPRESSION = 0
	LRIT_RICE      = 1
	JPEG           = 2
	GIF            = 5
	ZIP            = 10
)

var CompressionType = map[int]string{
	NO_COMPRESSION: "No Compression",
	LRIT_RICE:      "Goloumb Rice (LRIT)",
	JPEG:           "JPEG",
	GIF:            "GIF",
	ZIP:            "ZIP",
}

var CompressionTypeExtension = map[int]string{
	NO_COMPRESSION: ".bin",
	LRIT_RICE:      ".bin",
	JPEG:           ".jpg",
	GIF:            ".gif",
	ZIP:            ".zip",
}

func GetCompressionTypeString(compressionType int) string {
	v, ok := CompressionType[compressionType]
	if ok {
		return v
	}

	return fmt.Sprintf("Unknown (%d)", compressionType)
}

func GetCompressionTypeExtension(compressionType int) string {
	v, ok := CompressionTypeExtension[compressionType]
	if ok {
		return v
	}

	return ".data"
}
