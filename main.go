package main

import (
	"fmt"
	"log"

	"gotham/webhookapi"
)

func main() {
	fmt.Println("something")
	url := "https://webhook.site/177f5305-823d-4464-8dc2-5930b0586ae2"
	//  Create a new Client instance
	client, err := webhookapi.NewClient(url) // Replace with your actual API URL
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Define the destination configuration
	destination := webhookapi.Destination{
		Protocol:   webhookapi.HttpProto,
		URI:        "https://webhook.site/177f5305-823d-4464-8dc2-5930b0586ae2", // Replace with your actual endpoint
		HttpMethod: webhookapi.HttpPost,
		Encoding:   webhookapi.JSON,
	}

	// Define the extension ID
	extensionId := "test-extension-id" // Replace with your actual extension ID

	// Call the Send method
	response, err := client.Send(destination, extensionId)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	// Print the response
	fmt.Printf("Response: %s\n", response.Body)
}