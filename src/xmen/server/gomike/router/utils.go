package router

import (
	"fmt"
	"net/url"
)

// parse ID from URL query parameters
func ParseIDFromURL(urlPath *url.URL) (string, error) {
	listIDs := urlPath.Query()["id"]

	if len(listIDs) == 0 {
		err := fmt.Errorf("ID is required in the URL")
		return "", err
	}
	if len(listIDs) > 1 {
		err := fmt.Errorf("only one ID is allowed in the URL")
		return "", err
	}

	return listIDs[0], nil
}
