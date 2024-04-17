Below is a sample code for this library.

```
package main

import (
	"github.com/flywithu/azure-golang"
	"fmt"
	"log"
	"os"
)

func main() {
	VISION_ENDPOINT := "https://flywithufreevision.cognitiveservices.azure.com"
	VISION_KEY := os.Getenv("VISION_KEY")
	filePath := "/mnt/mypool/python_picture/filter_pic_gotest/0000/00/00000000-7deb47bf83b0fdf78259e9a8b47ac986.JPG"

	// URLs for image analysis
	// landmarkImageURL := "https://raw.githubusercontent.com/Azure-Samples/cognitive-services-sample-data-files/master/ComputerVision/Images/landmark.jpg"
	kitchenImageURL := "https://learn.microsoft.com/en-us/azure/ai-services/computer-vision/images/windows-kitchen.jpg"
	eiffelTowerImageURL := "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d4/Eiffel_Tower_20051010.jpg/1023px-Eiffel_Tower_20051010.jpg"
	redShirtLogoImageURL := "https://publish-p47754-e237306.adobeaemcloud.com/adobe/dynamicmedia/deliver/dm-aid--08fdf594-c963-43c8-b686-d4ba06727971/noticia_madridistas_hp.app.png?preferwebp=true&width=1440"

	client := vision.ComputerVisionClient(VISION_ENDPOINT, VISION_KEY)

	// Tagging an image
	log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")

	log.Println("Image Tagging")
	tags, err := client.GetImageTags(filePath)
	if err != nil {
		log.Fatalf("Failed to get image tags: %v", err)
	}
	for i, tag := range tags.Tags {
		if i >= 3 { break }
		fmt.Printf("Tag: %s (Confidence: %.2f)\n", tag.Name, tag.Confidence)
	}

	// Describing an image
	log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")

	log.Println("Image Description")
	description, err := client.GetImageDesc(filePath)
	if err != nil {
		log.Fatalf("Failed to get image description: %v", err)
	}
	for i, tag := range description.Description.Tags {
		if i >= 3 { break }
		fmt.Printf("Description Tag: %s\n", tag)
	}
	for i, caption := range description.Description.Captions {
		if i >= 3 { break }
		fmt.Printf("Caption: %s (Confidence: %.2f)\n", caption.Text, caption.Confidence)
	}

	// Object Detection
	log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	log.Println("Object Detection")
	objects, err := client.GetImageObject(kitchenImageURL)
	if err != nil {
		log.Fatalf("Failed to detect objects: %v", err)
	}
	for i, obj := range objects.Objects {
		if i >= 3 { break }
		fmt.Printf("Object: %s (Confidence: %.2f)\n", obj.ObjectInfo.ObjectName, obj.ObjectInfo.Confidence)
		if obj.Parent != nil {
			fmt.Printf("Parent Object: %s (Confidence: %.2f)\n", obj.Parent.ObjectName, obj.Parent.Confidence)
		}
	}

	// Analyzing image for landmarks
	log.Println("++++++++++++++++++++++++++++++++++++++++++++++++")

	log.Println("Landmark Analysis")
	landmarks, err := client.GetImageAnalyze(eiffelTowerImageURL)
	if err != nil {
		fmt.Printf("Failed to analyze landmarks: %v", err)
	}
	for i, cat := range landmarks.Categories {
		if i >= 3 { break }
		fmt.Printf("Category: %v\n", cat)
		if cat.Detail != nil && len(cat.Detail.Landmarks) > 0 {
			fmt.Printf("Landmark: %v\n", cat.Detail.Landmarks[0].Name)
		}
	}

	// Analyzing brands
	log.Println("Analyze - Brands")
	brands, err := client.GetImageAnalyze(redShirtLogoImageURL)
	if err != nil {
		fmt.Printf("Failed to analyze brands: %v", err)
	}
	fmt.Printf("Brand : %v \n",brands.Brands)
	for i, tag := range brands.Tags {
		if i >= 3 { break }
		fmt.Printf("Brand Tag: %v \n", tag)
	}
}

```
