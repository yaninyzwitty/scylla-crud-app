package configuration

import (
	"os"
	"testing"
)

// TestNewConfig tests the NewConfig function

func TestNewConfig(t *testing.T) {
	// Set up environment variables for the test
	os.Setenv("PORT", "9090")
	os.Setenv("HOSTS", "localhost")

	// Call NewConfig
	config, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() returned an error: %v", err)
	}

	// Check the PORT value
	if config.PORT != "9090" {
		t.Errorf("Expected PORT to be '9090', got '%s'", config.PORT)
	}

	// Check the HOSTS value
	if config.HOSTS != "localhost" {
		t.Errorf("Expected HOSTS to be 'localhost', got '%s'", config.HOSTS)
	}

	// Clean up
	os.Unsetenv("PORT")
	os.Unsetenv("HOSTS")
}

// TestNewConfigWithFallback tests the NewConfig function with fallback values

func TestNewConfigWithFallback(t *testing.T) {
	// here we make sure that env variables not set
	os.Unsetenv("PORT")
	os.Unsetenv("HOSTS")

	// WE CALL THE CONFIG FUNC
	config, err := NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() returned error: %v", err)
	}

	// Check the default PORT value
	if config.PORT != "8080" {
		t.Errorf("Expected PORT to be 8080, got %s", config.PORT)
	}

	// Check the default HOSTS value

	if config.HOSTS != "127.0.0.1" {
		t.Errorf("Expected HOSTS to be localhost, got %s", config.HOSTS)
	}

}
