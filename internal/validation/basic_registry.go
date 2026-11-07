package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// ValidateDockerImageName validates Docker image name.
func ValidateDockerImageName(name string) error {
	if name == "" {
		return newValidationError(
			domain.ErrInvalidDockerImage,
			"Docker image name is required",
			"Docker image name cannot be empty",
			"docker_image",
			"Specify a Docker image name",
		)
	}

	parts := strings.SplitN(name, ":", 2)
	imageName := parts[0]

	if len(imageName) > 255 {
		return newValidationError(
			domain.ErrInvalidDockerImage,
			"Docker image name too long",
			"Docker image name must be 255 characters or less",
			"docker_image",
			"Use a shorter image name",
		)
	}

	if !dockerImagePattern.MatchString(imageName) {
		return newValidationError(
			domain.ErrInvalidDockerImage,
			"Invalid Docker image name format",
			"Docker image name must contain only lowercase letters, numbers, dots, hyphens, and underscores",
			"docker_image",
			"Use lowercase alphanumeric characters with valid separators",
		)
	}

	if len(parts) == 2 {
		tag := parts[1]

		err := validateDockerTag(tag)
		if err != nil {
			return err
		}
	}

	return validateReservedDockerImageName(imageName, name)
}

var dockerTagPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func validateDockerTag(tag string) error {
	if len(tag) > 128 {
		return newValidationError(
			domain.ErrInvalidDockerImage,
			"Docker image tag too long",
			"Docker image tag must be 128 characters or less",
			"docker_image",
			"Use a shorter tag",
		)
	}

	if !strings.HasPrefix(tag, "v") && !dockerTagPattern.MatchString(tag) {
		return newValidationError(
			domain.ErrInvalidDockerImage,
			"Invalid Docker image tag format",
			"Docker image tag must start with 'v' for version tags or contain only valid characters",
			"docker_image",
			"Use semantic versioning (e.g., v1.0.0)",
		)
	}

	return nil
}

var reservedDockerImageNames = []string{
	"latest", "stable", "production", "master", "main", "develop",
	"test", "temp", "debug", "dev", "staging",
}

func validateReservedDockerImageName(imageName, fullName string) error {
	imageNameLower := strings.ToLower(imageName)
	for _, reserved := range reservedDockerImageNames {
		if strings.Contains(imageNameLower, reserved) && imageNameLower != reserved {
			continue
		}

		if imageNameLower == reserved {
			return newValidationError(
				domain.ErrInvalidDockerImage,
				"Reserved Docker image name",
				fmt.Sprintf("'%s' is a reserved Docker image name", fullName),
				"docker_image",
				"Choose a different image name",
			)
		}
	}

	return nil
}

// ValidateDockerRegistry validates Docker registry URL.
func ValidateDockerRegistry(registry string) error {
	if registry == "" {
		return nil // Registry can be empty (use default)
	}

	if len(registry) > 253 {
		return newValidationError(
			domain.ErrInvalidDockerRegistry,
			"Docker registry URL too long",
			"Docker registry URL must be 253 characters or less",
			"docker_registry",
			"Use a shorter registry URL",
		)
	}

	if !dockerRegistryPattern.MatchString(registry) {
		return newValidationError(
			domain.ErrInvalidDockerRegistry,
			"Invalid Docker registry URL format",
			"Docker registry URL must be a valid domain name",
			"docker_registry",
			"Use a valid domain name (e.g., registry.example.com)",
		)
	}

	if strings.Contains(registry, "..") {
		return newValidationError(
			domain.ErrInvalidDockerRegistry,
			"Invalid Docker registry URL format",
			"Docker registry URL cannot contain consecutive dots",
			"docker_registry",
			"Use a valid domain name (e.g., registry.example.com)",
		)
	}

	knownRegistries := []string{
		"docker.io", "ghcr.io", "registry.gitlab.com", "gcr.io",
		"azurecr.io", "quay.io", "mcr.microsoft.com",
	}

	for _, known := range knownRegistries {
		if strings.EqualFold(registry, known) {
			return nil
		}
	}

	return nil
}

// ValidateVersion validates version string.
