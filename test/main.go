package main

import (
	"fmt"

	witai "github.com/wit-ai/wit-go/v2"
)

func main() {
	witToken := "UVJRI7BCS4FOHCMPRMPUZBXUM5GJBHC6"

	if witToken == "" {
		fmt.Println("Wit.ai API token is missing. Set the environment variable WIT_AI_API_TOKEN.")
		return
	}

	client := witai.NewClient(witToken)
	fmt.Println(client)
	msg, err := client.Parse(&witai.MessageRequest{
		Query: "i dont like you",
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("/////////////////////// TRAITS ")
	for traitName, traitValues := range msg.Traits {
		fmt.Println("Trait:", traitName)
		for _, traitValue := range traitValues {
			fmt.Println("  Value:", traitValue.Value)
			fmt.Println("  Confidence:", traitValue.Confidence)
		}
	}
	fmt.Println("/////////////////////// TRAITS ")
	fmt.Println("/////////////////////// ENTITIES ")
	for entityname, entityValues := range msg.Entities {
		fmt.Println("Entity:", entityname)
		for _, entityValye := range entityValues {
			fmt.Println("  Value:", entityValye.Value)
			fmt.Println("  Confidence:", entityValye.Confidence)
		}
	}
	fmt.Println("/////////////////////// ENTITIES ")
}
