package url

import (
	"strings"
	"testing"
)

func TestParseGitHubReleaseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want *ParsedGitHubRelease
	}{
		{
			name: "valid tar.gz URL",
			url:  "https://github.com/github/copilot-cli/releases/download/v0.0.400/copilot-linux-x64.tar.gz",
			want: &ParsedGitHubRelease{
				Owner:     "github",
				Repo:      "copilot-cli",
				Version:   "v0.0.400",
				AssetName: "copilot-linux-x64.tar.gz",
				Format:    ".tar.gz",
			},
		},
		{
			name: "valid zip URL",
			url:  "https://github.com/oven-sh/bun/releases/download/bun-v1.0.0/bun-linux-x64.zip",
			want: &ParsedGitHubRelease{
				Owner:     "oven-sh",
				Repo:      "bun",
				Version:   "bun-v1.0.0",
				AssetName: "bun-linux-x64.zip",
				Format:    ".zip",
			},
		},
		{
			name: "valid tgz URL",
			url:  "https://github.com/cli/cli/releases/download/v2.0.0/gh_2.0.0_linux_amd64.tgz",
			want: &ParsedGitHubRelease{
				Owner:     "cli",
				Repo:      "cli",
				Version:   "v2.0.0",
				AssetName: "gh_2.0.0_linux_amd64.tgz",
				Format:    ".tar.gz",
			},
		},
		{
			name: "URL with version without v prefix",
			url:  "https://github.com/sharkdp/bat/releases/download/0.24.0/bat-v0.24.0-x86_64-unknown-linux-gnu.tar.gz",
			want: &ParsedGitHubRelease{
				Owner:     "sharkdp",
				Repo:      "bat",
				Version:   "0.24.0",
				AssetName: "bat-v0.24.0-x86_64-unknown-linux-gnu.tar.gz",
				Format:    ".tar.gz",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGitHubReleaseURL(tt.url)

			if err != nil {
				t.Errorf("ParseGitHubReleaseURL() unexpected error: %v", err)
				return
			}

			if got.Owner != tt.want.Owner {
				t.Errorf("Owner = %v, want %v", got.Owner, tt.want.Owner)
			}
			if got.Repo != tt.want.Repo {
				t.Errorf("Repo = %v, want %v", got.Repo, tt.want.Repo)
			}
			if got.Version != tt.want.Version {
				t.Errorf("Version = %v, want %v", got.Version, tt.want.Version)
			}
			if got.AssetName != tt.want.AssetName {
				t.Errorf("AssetName = %v, want %v", got.AssetName, tt.want.AssetName)
			}
			if got.Format != tt.want.Format {
				t.Errorf("Format = %v, want %v", got.Format, tt.want.Format)
			}
		})
	}
}

func TestParseGitHubReleaseURL_NonGitHubURL(t *testing.T) {
	url := "https://example.com/owner/repo/releases/download/v1.0.0/asset.tar.gz"
	got, err := ParseGitHubReleaseURL(url)

	if err == nil {
		t.Errorf("ParseGitHubReleaseURL() expected error but got none")
	}
	if !strings.Contains(err.Error(), "not a GitHub URL") {
		t.Errorf("ParseGitHubReleaseURL() expected 'not a GitHub URL' error, got %v", err)
	}
	if got != nil {
		t.Errorf("ParseGitHubReleaseURL() expected nil, got %v", got)
	}
}

func TestParseGitHubReleaseURL_MissingReleases(t *testing.T) {
	url := "https://github.com/owner/repo/download/v1.0.0/asset.tar.gz"
	got, err := ParseGitHubReleaseURL(url)

	if err == nil {
		t.Errorf("ParseGitHubReleaseURL() expected error but got none")
	}
	if !strings.Contains(err.Error(), "invalid GitHub release URL") {
		t.Errorf("ParseGitHubReleaseURL() expected 'invalid GitHub release URL' error, got %v", err)
	}
	if got != nil {
		t.Errorf("ParseGitHubReleaseURL() expected nil, got %v", got)
	}
}

func TestParseGitHubReleaseURL_TooFewSegments(t *testing.T) {
	url := "https://github.com/owner/repo"
	got, err := ParseGitHubReleaseURL(url)

	if err == nil {
		t.Errorf("ParseGitHubReleaseURL() expected error but got none")
	}
	if !strings.Contains(err.Error(), "invalid GitHub release URL format") {
		t.Errorf("ParseGitHubReleaseURL() expected 'invalid GitHub release URL format' error, got %v", err)
	}
	if got != nil {
		t.Errorf("ParseGitHubReleaseURL() expected nil, got %v", got)
	}
}

func TestParseGitHubReleaseURL_UnsupportedFormat(t *testing.T) {
	url := "https://github.com/owner/repo/releases/download/v1.0.0/asset.exe"
	got, err := ParseGitHubReleaseURL(url)

	if err == nil {
		t.Errorf("ParseGitHubReleaseURL() expected error but got none")
	}
	if !strings.Contains(err.Error(), "unsupported file format") {
		t.Errorf("ParseGitHubReleaseURL() expected 'unsupported file format' error, got %v", err)
	}
	if got != nil {
		t.Errorf("ParseGitHubReleaseURL() expected nil, got %v", got)
	}
}

func TestParseGitHubReleaseURL_MalformedURL(t *testing.T) {
	url := "http://%zz"
	got, err := ParseGitHubReleaseURL(url)

	if err == nil {
		t.Errorf("ParseGitHubReleaseURL() expected error but got none")
	}
	if !strings.Contains(err.Error(), "invalid URL") {
		t.Errorf("ParseGitHubReleaseURL() expected 'invalid URL' error, got %v", err)
	}
	if got != nil {
		t.Errorf("ParseGitHubReleaseURL() expected nil, got %v", got)
	}
}

func TestParseGitHubReleaseURL_InvalidReleaseURL(t *testing.T) {
	url := "https://github.com/github/copilot-cli/issues/upload/v0.0.400/copilot-linux-x64.tar.gz"
	got, err := ParseGitHubReleaseURL(url)

	if err == nil {
		t.Errorf("ParseGitHubReleaseURL() expected error but got none")
	}
	if !strings.Contains(err.Error(), "invalid GitHub release URL: expected /releases/download/ in path") {
		t.Errorf("ParseGitHubReleaseURL() expected 'invalid GitHub release URL' error, got %v", err)
	}
	if got != nil {
		t.Errorf("ParseGitHubReleaseURL() expected nil, got %v", got)
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{
			name:     "tar.gz file",
			fileName: "binary.tar.gz",
			want:     ".tar.gz",
		},
		{
			name:     "zip file",
			fileName: "binary.zip",
			want:     ".zip",
		},
		{
			name:     "tgz file",
			fileName: "binary.tgz",
			want:     ".tar.gz",
		},
		{
			name:     "unsupported format",
			fileName: "binary.exe",
			want:     "",
		},
		{
			name:     "no extension",
			fileName: "binary",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectFormat(tt.fileName)
			if got != tt.want {
				t.Errorf("detectFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateBinaryID(t *testing.T) {
	tests := []struct {
		name      string
		assetName string
		want      string
	}{
		{
			name:      "copilot asset",
			assetName: "copilot-linux-x64.tar.gz",
			want:      "copilot",
		},
		{
			name:      "gh asset",
			assetName: "gh_2.0.0_linux_amd64.tar.gz",
			want:      "gh",
		},
		{
			name:      "bat asset",
			assetName: "bat-v0.24.0-x86_64-unknown-linux-gnu.tar.gz",
			want:      "bat",
		},
		{
			name:      "bun asset",
			assetName: "bun-linux-x64.zip",
			want:      "bun",
		},
		{
			name:      "simple name",
			assetName: "binary.tar.gz",
			want:      "binary",
		},
		{
			name:      "underscore separator",
			assetName: "my_binary_linux_amd64.zip",
			want:      "my",
		},
		{
			name:      "raw asset",
			assetName: "bun-linux-x64",
			want:      "bun",
		},
		{
			name:      "empty asset name",
			assetName: "",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateBinaryID(tt.assetName)
			if got != tt.want {
				t.Errorf("GenerateBinaryID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateBinaryName(t *testing.T) {
	tests := []struct {
		name      string
		assetName string
		want      string
	}{
		{
			name:      "copilot asset",
			assetName: "copilot-linux-x64.tar.gz",
			want:      "copilot",
		},
		{
			name:      "gh asset",
			assetName: "gh_2.0.0_linux_amd64.tar.gz",
			want:      "gh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateBinaryName(tt.assetName)
			if got != tt.want {
				t.Errorf("GenerateBinaryName() = %v, want %v", got, tt.want)
			}
		})
	}
}
