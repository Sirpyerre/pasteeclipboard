package monitor

import (
	"strings"
	"testing"

	"github.com/Sirpyerre/pasteeclipboard/internal/database"
)

func TestDetectContentType_URL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://google.com", "link"},
		{"http://example.org/path", "link"},
		{"https://solobaja.mx/en/avita-en/", "link"},
		{"www.google.com", "link"},
		{"www.example.org/path?query=1", "link"},
	}

	for _, tt := range tests {
		result := DetectContentType(tt.input)
		if result != tt.expected {
			t.Errorf("DetectContentType(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDetectContentType_Email(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user@example.com", "email"},
		{"test.email@domain.org", "email"},
		{"jpichardini@solobaja.mx", "email"},
		{"name+tag@gmail.com", "email"},
	}

	for _, tt := range tests {
		result := DetectContentType(tt.input)
		if result != tt.expected {
			t.Errorf("DetectContentType(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDetectContentType_Phone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"612 132 8036", "phone"},
		{"1234567890", "phone"},
		{"+52-612-132-8036", "phone"},
		{"(612) 132-8036", "phone"},
	}

	for _, tt := range tests {
		result := DetectContentType(tt.input)
		if result != tt.expected {
			t.Errorf("DetectContentType(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDetectContentType_Text(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello world", "text"},
		{"Some random text here", "text"},
		{"123", "text"}, // Too short for phone
		{"not-an-email", "text"},
		{"ftp://server.com", "text"}, // Not http/https
	}

	for _, tt := range tests {
		result := DetectContentType(tt.input)
		if result != tt.expected {
			t.Errorf("DetectContentType(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestContainsDigits(t *testing.T) {
	tests := []struct {
		input    string
		n        int
		expected bool
	}{
		{"1234567", 7, true},
		{"123456", 7, false},
		{"abc123def456", 6, true},
		{"no digits here", 1, false},
		{"612 132 8036", 10, true},
	}

	for _, tt := range tests {
		result := containsDigits(tt.input, tt.n)
		if result != tt.expected {
			t.Errorf("containsDigits(%q, %d) = %v, want %v", tt.input, tt.n, result, tt.expected)
		}
	}
}

func TestTextTruncation(t *testing.T) {
	// Test that truncation happens at the right length
	maxLen := database.MaxTextLength

	// Test content under limit - should not be truncated
	shortContent := strings.Repeat("a", 100)
	if len(shortContent) > maxLen {
		t.Error("Short content should not exceed max length")
	}

	// Test content over limit - verify truncation logic
	longContent := strings.Repeat("x", maxLen+1000)
	if len(longContent) <= maxLen {
		t.Error("Long content should exceed max length for this test")
	}

	// Simulate truncation logic from handleTextClipboard
	truncated := longContent
	if len(truncated) > maxLen {
		truncated = truncated[:maxLen] + "\n... (truncated)"
	}

	// Verify truncated content is at expected length
	expectedLen := maxLen + len("\n... (truncated)")
	if len(truncated) != expectedLen {
		t.Errorf("Truncated length = %d, want %d", len(truncated), expectedLen)
	}

	// Verify it ends with truncation marker
	if !strings.HasSuffix(truncated, "\n... (truncated)") {
		t.Error("Truncated content should end with truncation marker")
	}
}

func TestTextTruncation_ExactLimit(t *testing.T) {
	maxLen := database.MaxTextLength

	// Test content exactly at limit - should not be truncated
	exactContent := strings.Repeat("b", maxLen)

	truncated := exactContent
	if len(truncated) > maxLen {
		truncated = truncated[:maxLen] + "\n... (truncated)"
	}

	// Should remain unchanged (not truncated)
	if len(truncated) != maxLen {
		t.Errorf("Content at exact limit should not be truncated, got length %d", len(truncated))
	}
}

func TestURLRegex(t *testing.T) {
	validURLs := []string{
		"https://example.com",
		"http://test.org",
		"https://sub.domain.com/path",
		"www.google.com",
	}

	invalidURLs := []string{
		"ftp://server.com",
		"example.com", // missing protocol/www
		"not a url",
	}

	for _, url := range validURLs {
		if !urlRegex.MatchString(url) {
			t.Errorf("urlRegex should match %q", url)
		}
	}

	for _, url := range invalidURLs {
		if urlRegex.MatchString(url) {
			t.Errorf("urlRegex should not match %q", url)
		}
	}
}

func TestEmailRegex(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.org",
		"user+tag@gmail.com",
	}

	invalidEmails := []string{
		"notanemail",
		"@nodomain.com",
		"noat.com",
		"spaces in@email.com",
	}

	for _, email := range validEmails {
		if !emailRegex.MatchString(email) {
			t.Errorf("emailRegex should match %q", email)
		}
	}

	for _, email := range invalidEmails {
		if emailRegex.MatchString(email) {
			t.Errorf("emailRegex should not match %q", email)
		}
	}
}

func TestPhoneRegex(t *testing.T) {
	validPhones := []string{
		"1234567",
		"612 132 8036",
		"+1-234-567-8901",
		"(123) 456-7890",
	}

	invalidPhones := []string{
		"123", // too short
		"abcdefghij", // no digits
		"12345678901234567890123456789", // too long (>20)
	}

	for _, phone := range validPhones {
		if !phoneRegex.MatchString(phone) {
			t.Errorf("phoneRegex should match %q", phone)
		}
	}

	for _, phone := range invalidPhones {
		if phoneRegex.MatchString(phone) {
			t.Errorf("phoneRegex should not match %q", phone)
		}
	}
}
