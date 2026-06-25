package github

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"cturner8/binmate/internal/database"
)

// writeScript writes a shell script that echoes token to a temp file and makes it executable.
func writeScript(t *testing.T, token string) string {
	t.Helper()
	dir := t.TempDir()
	script := filepath.Join(dir, "askpass.sh")
	content := "#!/bin/sh\necho " + token + "\n"
	if err := os.WriteFile(script, []byte(content), 0o700); err != nil {
		t.Fatalf("writeScript: %v", err)
	}
	return script
}

// clearAuthEnv unsets all authentication-related environment variables.
func clearAuthEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"BINMATE_GITHUB_TOKEN", "GITHUB_TOKEN", "BINMATE_GITHUB_ASKPASS", "GITHUB_ASKPASS"} {
		t.Setenv(key, "")
	}
}

func TestResolveAskpass(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell scripts not supported on Windows")
	}

	t.Run("returns token from script stdout", func(t *testing.T) {
		script := writeScript(t, "my-secret-token")
		got, err := resolveAskpass(script)
		if err != nil {
			t.Fatalf("resolveAskpass() unexpected error: %v", err)
		}
		if got != "my-secret-token" {
			t.Errorf("resolveAskpass() = %q, want %q", got, "my-secret-token")
		}
	})

	t.Run("error on empty script path", func(t *testing.T) {
		_, err := resolveAskpass("")
		if err == nil {
			t.Fatal("resolveAskpass(\"\") expected error, got nil")
		}
	})

	t.Run("error on non-existent script", func(t *testing.T) {
		_, err := resolveAskpass("/nonexistent/script.sh")
		if err == nil {
			t.Fatal("resolveAskpass() expected error for missing script, got nil")
		}
	})

	t.Run("error when script returns empty output", func(t *testing.T) {
		dir := t.TempDir()
		script := filepath.Join(dir, "empty.sh")
		if err := os.WriteFile(script, []byte("#!/bin/sh\necho\n"), 0o700); err != nil {
			t.Fatal(err)
		}
		_, err := resolveAskpass(script)
		if err == nil {
			t.Fatal("resolveAskpass() expected error for empty output, got nil")
		}
	})
}

func TestResolveToken(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell scripts not supported on Windows")
	}

	t.Run("BINMATE_GITHUB_TOKEN takes precedence over GITHUB_TOKEN", func(t *testing.T) {
		clearAuthEnv(t)
		t.Setenv("BINMATE_GITHUB_TOKEN", "binmate-token")
		t.Setenv("GITHUB_TOKEN", "github-token")

		got, err := ResolveToken()
		if err != nil {
			t.Fatalf("ResolveToken() unexpected error: %v", err)
		}
		if got != "binmate-token" {
			t.Errorf("ResolveToken() = %q, want %q", got, "binmate-token")
		}
	})

	t.Run("falls back to GITHUB_TOKEN when BINMATE_GITHUB_TOKEN unset", func(t *testing.T) {
		clearAuthEnv(t)
		t.Setenv("GITHUB_TOKEN", "github-token")

		got, err := ResolveToken()
		if err != nil {
			t.Fatalf("ResolveToken() unexpected error: %v", err)
		}
		if got != "github-token" {
			t.Errorf("ResolveToken() = %q, want %q", got, "github-token")
		}
	})

	t.Run("BINMATE_GITHUB_ASKPASS takes precedence over GITHUB_ASKPASS", func(t *testing.T) {
		clearAuthEnv(t)
		binmateScript := writeScript(t, "binmate-askpass-token")
		githubScript := writeScript(t, "github-askpass-token")
		t.Setenv("BINMATE_GITHUB_ASKPASS", binmateScript)
		t.Setenv("GITHUB_ASKPASS", githubScript)

		got, err := ResolveToken()
		if err != nil {
			t.Fatalf("ResolveToken() unexpected error: %v", err)
		}
		if got != "binmate-askpass-token" {
			t.Errorf("ResolveToken() = %q, want %q", got, "binmate-askpass-token")
		}
	})

	t.Run("BINMATE_GITHUB_ASKPASS takes precedence over env tokens", func(t *testing.T) {
		clearAuthEnv(t)
		script := writeScript(t, "askpass-token")
		t.Setenv("BINMATE_GITHUB_ASKPASS", script)
		t.Setenv("BINMATE_GITHUB_TOKEN", "env-token")

		got, err := ResolveToken()
		if err != nil {
			t.Fatalf("ResolveToken() unexpected error: %v", err)
		}
		if got != "askpass-token" {
			t.Errorf("ResolveToken() = %q, want %q", got, "askpass-token")
		}
	})

	t.Run("GITHUB_ASKPASS used when BINMATE_GITHUB_ASKPASS unset", func(t *testing.T) {
		clearAuthEnv(t)
		script := writeScript(t, "github-askpass-token")
		t.Setenv("GITHUB_ASKPASS", script)

		got, err := ResolveToken()
		if err != nil {
			t.Fatalf("ResolveToken() unexpected error: %v", err)
		}
		if got != "github-askpass-token" {
			t.Errorf("ResolveToken() = %q, want %q", got, "github-askpass-token")
		}
	})

	t.Run("returns empty string when no token sources configured", func(t *testing.T) {
		clearAuthEnv(t)

		got, err := ResolveToken()
		if err != nil {
			t.Fatalf("ResolveToken() unexpected error: %v", err)
		}
		if got != "" {
			t.Errorf("ResolveToken() = %q, want empty string", got)
		}
	})

	t.Run("returns error when askpass script fails", func(t *testing.T) {
		clearAuthEnv(t)
		t.Setenv("BINMATE_GITHUB_ASKPASS", "/nonexistent/script.sh")

		_, err := ResolveToken()
		if err == nil {
			t.Fatal("ResolveToken() expected error for failed askpass, got nil")
		}
	})
}

func TestCreateHTTPClient(t *testing.T) {
	t.Run("empty token returns plain client", func(t *testing.T) {
		client, err := CreateHTTPClient("")
		if err != nil {
			t.Fatalf("CreateHTTPClient(\"\") unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("CreateHTTPClient(\"\") returned nil client")
		}
		if client.Transport != nil {
			t.Errorf("CreateHTTPClient(\"\") Transport = %v, want nil (plain client)", client.Transport)
		}
	})

	t.Run("non-empty token sets Authorization header", func(t *testing.T) {
		var receivedAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		client, err := CreateHTTPClient("testtoken")
		if err != nil {
			t.Fatalf("CreateHTTPClient() unexpected error: %v", err)
		}

		if _, err := client.Get(srv.URL); err != nil {
			t.Fatalf("client.Get() unexpected error: %v", err)
		}

		want := "Bearer testtoken"
		if receivedAuth != want {
			t.Errorf("Authorization header = %q, want %q", receivedAuth, want)
		}
	})
}

func TestNewClientForBinary(t *testing.T) {
	t.Run("unauthenticated binary returns plain client", func(t *testing.T) {
		clearAuthEnv(t)
		binary := &database.Binary{Authenticated: false, Name: "test"}
		client, err := NewClientForBinary(binary, false)
		if err != nil {
			t.Fatalf("NewClientForBinary() unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("NewClientForBinary() returned nil client")
		}
		if client.Transport != nil {
			t.Errorf("NewClientForBinary() Transport = %v, want nil (plain client)", client.Transport)
		}
	})

	t.Run("authenticated binary with token returns authenticated client", func(t *testing.T) {
		clearAuthEnv(t)
		t.Setenv("BINMATE_GITHUB_TOKEN", "my-token")

		var receivedAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		binary := &database.Binary{Authenticated: true, Name: "test"}
		client, err := NewClientForBinary(binary, false)
		if err != nil {
			t.Fatalf("NewClientForBinary() unexpected error: %v", err)
		}

		if _, err := client.Get(srv.URL); err != nil {
			t.Fatalf("client.Get() unexpected error: %v", err)
		}

		if receivedAuth != "Bearer my-token" {
			t.Errorf("Authorization header = %q, want %q", receivedAuth, "Bearer my-token")
		}
	})

	t.Run("returns authenticated client when binary is not configured for authentication but global provider authentication is enabled", func(t *testing.T) {
		clearAuthEnv(t)
		t.Setenv("BINMATE_GITHUB_TOKEN", "my-token")

		var receivedAuth string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		binary := &database.Binary{Authenticated: false, Name: "test"}
		client, err := NewClientForBinary(binary, true)
		if err != nil {
			t.Fatalf("NewClientForBinary() unexpected error: %v", err)
		}

		if _, err := client.Get(srv.URL); err != nil {
			t.Fatalf("client.Get() unexpected error: %v", err)
		}

		if receivedAuth != "Bearer my-token" {
			t.Errorf("Authorization header = %q, want %q", receivedAuth, "Bearer my-token")
		}
	})

	t.Run("authenticated binary with no token returns error", func(t *testing.T) {
		clearAuthEnv(t)
		binary := &database.Binary{Authenticated: true, Name: "test"}
		_, err := NewClientForBinary(binary, false)
		if err == nil {
			t.Fatal("NewClientForBinary() expected error when no token set, got nil")
		}
	})
}
