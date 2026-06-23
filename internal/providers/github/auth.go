package github

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"cturner8/binmate/internal/database"
)

// resolveAskpass executes the given askpass script and returns the token it prints.
// Returns an error if the script path is empty, the script exits with a non-zero
// status, or its output is empty.
func resolveAskpass(scriptPath string) (string, error) {
	if scriptPath == "" {
		return "", fmt.Errorf("askpass script path is empty")
	}

	out, err := exec.Command(scriptPath).Output()
	if err != nil {
		return "", fmt.Errorf("askpass script %q failed: %w", scriptPath, err)
	}

	token := strings.TrimSpace(string(out))
	if token == "" {
		return "", fmt.Errorf("askpass script %q returned empty output", scriptPath)
	}

	return token, nil
}

// ResolveToken resolves a GitHub token by trying the following sources in order:
//  1. BINMATE_GITHUB_ASKPASS script
//  2. GITHUB_ASKPASS script
//  3. BINMATE_GITHUB_TOKEN environment variable
//  4. GITHUB_TOKEN environment variable
//
// Returns the first non-empty token found, or an empty string if none is set.
func ResolveToken() (string, error) {
	if script := os.Getenv("BINMATE_GITHUB_ASKPASS"); script != "" {
		token, err := resolveAskpass(script)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	if script := os.Getenv("GITHUB_ASKPASS"); script != "" {
		token, err := resolveAskpass(script)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	if token := os.Getenv("BINMATE_GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", nil
}

// CreateHTTPClient creates an HTTP client authenticated with the given token.
// An empty token returns a plain unauthenticated client.
func CreateHTTPClient(token string) (*http.Client, error) {
	client := &http.Client{}

	if token == "" {
		return client, nil
	}

	transport := &authenticatedTransport{
		token:     token,
		transport: http.DefaultTransport,
	}

	client.Transport = transport
	return client, nil
}

// NewClientForBinary creates an HTTP client appropriate for the given binary.
// If the binary does not require authentication a plain client is returned.
// If authentication is required, ResolveToken is called; an error is returned
// if no token can be resolved.
func NewClientForBinary(binary *database.Binary) (*http.Client, error) {
	if !binary.Authenticated {
		return &http.Client{}, nil
	}

	token, err := ResolveToken()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve GitHub token: %w", err)
	}

	if token == "" {
		return nil, fmt.Errorf("authentication required for binary %q but no GitHub token found; set an askpass script or static token in your environment", binary.Name)
	}

	return CreateHTTPClient(token)
}

// authenticatedTransport is an http.RoundTripper that adds GitHub authentication.
type authenticatedTransport struct {
	token     string
	transport http.RoundTripper
}

func (t *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.transport.RoundTrip(req)
}
