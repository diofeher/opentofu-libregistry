package github

import (
	"fmt"
	"net/http"
)

func (c Client) IsUserInOrganization(username string, org string) (bool, error) {
	// user/org is not case sensitive here
	check_url := fmt.Sprintf("https://api.github.com/orgs/%s/public_members/%s", org, username)

	resp, err := c.httpClient.Get(check_url)
	if err != nil {
		return false, err
	}
	resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return false, nil
	case http.StatusNoContent:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected status code %v when checking if %q is a member of %q", resp.StatusCode, username, org)
	}
}
