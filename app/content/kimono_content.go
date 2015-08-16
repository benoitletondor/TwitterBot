package content

import (
	"encoding/json"
	"errors"
	"fmt"
)

type KimonoContent struct {
	Url string
}

func (kimono KimonoContent) callAPI() ([]Content, error) {
	content, err := getWebserviceContent(kimono.Url)
	if err != nil {
		fmt.Println("Error while calling API")
		return nil, err
	} else {
		// Fill the record with the data from the JSON
		var jsonObj map[string]interface{}
		err = json.Unmarshal(content, &jsonObj)
		if err != nil {
			fmt.Println("An error occurred while converting our JSON to an object")
			return nil, err
		}

		contents := make([]Content, 0)

		if val, ok := jsonObj["results"]; ok {

			if results, ok := val.(map[string]interface{}); ok {

				if val, ok := results["collection1"]; ok {

					if collection, ok := val.([]interface{}); ok {

						for _, item := range collection {

							if property, ok := item.(map[string]interface{}); ok {

								if property1, ok := property["property1"]; ok {

									if contentData, ok := property1.(map[string]interface{}); ok {

										if href, ok := contentData["href"].(string); ok {

											if property2, ok := property["property2"]; ok {

												if contentData2, ok := property2.(map[string]interface{}); ok {

													if text, ok := contentData2["text"].(string); ok {

														contents = append(contents, Content{Text: text, Url: href})

													} else {
														fmt.Println("Error mapping text as string")
														return nil, errors.New("json mapping error")
													}
												} else {
													fmt.Println("Error mapping property2 as json obj")
													return nil, errors.New("json mapping error")
												}
											} else {
												fmt.Println("Missing mappings property2")
												return nil, errors.New("json mapping error")
											}

										} else {
											fmt.Println("Error mapping href as string")
											return nil, errors.New("json mapping error")
										}

									} else {
										fmt.Println("Error mapping property1 as json obj")
										return nil, errors.New("json mapping error")
									}

								} else {
									fmt.Println("Missing mappings property1")
									return nil, errors.New("json mapping error")
								}

							} else {
								fmt.Println("Error mappings property1")
								return nil, errors.New("json mapping error")
							}

						}

					} else {
						fmt.Println("Error mappings collection1 as array")
						return nil, errors.New("json mapping error")
					}

				} else {
					fmt.Println("Error mappings collection1")
					return nil, errors.New("json mapping error")
				}

			} else {
				fmt.Println("Error mappings results as json obj")
				return nil, errors.New("json mapping error")
			}

			return contents, nil

		} else {
			fmt.Println("No field results in json")
			return nil, errors.New("json mapping error")
		}
	}
}
