// cloudinary/cloudinary.go
package cloudinary

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var CLD *cloudinary.Cloudinary
var ctx = context.Background()

func InitCloudinary() {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		fmt.Println("Warning: Cloudinary credentials not set in environment variables. Image uploads will fail.")
		return
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		fmt.Printf("Failed to initialize Cloudinary: %v\n", err)
		return
	}
	CLD = cld
	fmt.Println("Cloudinary initialized successfully.")
}

// UploadImage uploads a file to Cloudinary and returns its URL.
// filePath is the path to the temporary file or a multipart.FileHeader.
func UploadImage(file interface{}) (string, error) {
	if CLD == nil {
		return "", fmt.Errorf("Cloudinary not initialized")
	}

	uploadResult, err := CLD.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "gin-blog-profile-pictures", // Optional: organizes uploads into a specific folder
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}
