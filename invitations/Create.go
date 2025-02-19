package invitations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type InvitationsCreateResult struct {
	Result []InvitationsCreateResponse `json:"result"`
	Error  *Error                      `json:"error"`
}

type InvitationsCreateResponse struct {
	ID        int64                 `json:"id"`
	Email     string                `json:"email"`
	Role      InvitationsCreateRole `json:"role"`
	CreatedAt string                `json:"created_at"`
}

type InvitationsCreateOptions struct {
	Email string                `json:"email"`
	Role  InvitationsCreateRole `json:"role"`
}

type InvitationsCreateRole string

const (
	Manager    InvitationsCreateRole = "manager"
	Member     InvitationsCreateRole = "member"
	Viewer     InvitationsCreateRole = "viewer"
	Restricted InvitationsCreateRole = "restricted"
)

func (c *InvitationsClient) Create(lo InvitationsCreateOptions) (*InvitationsCreateResult, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(lo)
	if err != nil {
		return nil, err
	}

	req, err := c.client.NewRequest("POST", "/v1/invitations", b)
	if err != nil {
		return nil, err
	}

	r, err := c.client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusCreated && r.StatusCode != http.StatusOK {
		var target Error
		decode := json.NewDecoder(r.Body)
		err = decode.Decode(&target)
		if err != nil {
			return nil, err
		}
		return &InvitationsCreateResult{Error: &target}, fmt.Errorf("%s", target.Detail)
	}

	if r.StatusCode == http.StatusCreated {
		var target InvitationsCreateResponse
		decode := json.NewDecoder(r.Body)
		err = decode.Decode(&target)
		if err != nil {
			return nil, err
		}
		return &InvitationsCreateResult{Result: []InvitationsCreateResponse{target}}, nil
	}

	return nil, fmt.Errorf("user already in the workspace")
}
