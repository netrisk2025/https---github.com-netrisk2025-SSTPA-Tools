package manifest

import "fmt"

type NormalizedReferenceItem struct {
	FrameworkName    string   `json:"frameworkName"`
	FrameworkVersion string   `json:"frameworkVersion"`
	ExternalID       string   `json:"externalId"`
	ExternalType     string   `json:"externalType"`
	Name             string   `json:"name"`
	ShortDescription string   `json:"shortDescription"`
	LongDescription  string   `json:"longDescription"`
	SourceURI        string   `json:"sourceUri"`
	Tags             []string `json:"tags"`
}

func (item NormalizedReferenceItem) Validate() error {
	if item.FrameworkName == "" {
		return fmt.Errorf("framework name is required")
	}

	if item.FrameworkVersion == "" {
		return fmt.Errorf("framework version is required")
	}

	if item.ExternalID == "" {
		return fmt.Errorf("external id is required")
	}

	if item.ExternalType == "" {
		return fmt.Errorf("external type is required")
	}

	if item.Name == "" {
		return fmt.Errorf("name is required")
	}

	if item.SourceURI == "" {
		return fmt.Errorf("source uri is required")
	}

	return nil
}
