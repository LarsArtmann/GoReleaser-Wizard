package validation

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/domain"
)

// newValidationError creates a new validation error with field and suggestion.
func newValidationError(
	code domain.ErrorCode,
	message, details, field, suggestion string,
) error {
	return domain.NewValidationError(code, message, details).
		WithField(field).
		WithSuggestion(suggestion)
}

// Validation patterns.
var (
	// Project name validation: alphanumeric, hyphens, underscores, dots, 1-50 chars.
	projectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,49}$`)

	// Binary name validation: must start with letter, then alphanumeric, hyphens, underscores, 1-30 chars.
	binaryNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{0,29}$`)

	// Main path validation: relative path only, no absolute paths, no invalid characters.
	mainPathPattern = regexp.MustCompile(`^(\.|\./|\.\./)?[^\\:*?"<>|]+$`)

	// Project description validation: 1-500 printable chars including newlines.
	projectDescriptionPattern = regexp.MustCompile(`^[[:print:]\n\r]{1,500}$`)

	// Docker image name validation: lowercase, alphanum, separators, max 255.
	dockerImagePattern = regexp.MustCompile(`^[a-z0-9]+(?:[._-][a-z0-9]+)*$`)

	// Docker registry URL validation - allows domains, localhost, and IP addresses with optional paths.
	dockerRegistryPattern = regexp.MustCompile(
		`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9.-]*[a-zA-Z0-9])?|[lL]ocalhost|127\.0\.0\.1)(?::[0-9]+)?(?:/[a-zA-Z0-9._-]+)*/?$`,
	)
)

// Reserved names shared across validators.
var reservedWindowsDeviceNames = []string{
	"con", "prn", "aux", "nul", "com1", "com2", "com3",
	"com4", "com5", "com6", "com7", "com8", "com9", "lpt1",
	"lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9",
}

// reservedProjectNames includes Windows devices and common reserved words.
var reservedProjectNames = []string{
	"go", "test", "temp", "tmp", "default", "admin", "root", "system",
	"null", "undefined", "none", "false", "true",
}

// reservedBinaryNames adds common commands to reserved names.
var reservedBinaryNames = []string{
	"test", "temp", "tmp", "debug", "release", "build", "install",
	"setup", "config", "init", "run", "start", "stop", "restart",
}

// ValidateProjectName validates project name.
func ValidateProjectName(name string) error {
	if name == "" {
		return newValidationError(
			domain.ErrInvalidProject,
			"Project name is required",
			"Project name cannot be empty",
			"project_name",
			"Choose a valid project name",
		)
	}

	if len(name) > 50 {
		return newValidationError(
			domain.ErrInvalidProject,
			"Project name too long",
			"Project name must be 50 characters or less",
			"project_name",
			"Use a shorter project name",
		)
	}

	if !projectNamePattern.MatchString(name) {
		return newValidationError(
			domain.ErrInvalidProject,
			"Invalid project name format",
			"Project name must contain only letters, numbers, hyphens, and underscores",
			"project_name",
			"Use alphanumeric characters with hyphens and underscores only",
		)
	}

	if strings.HasSuffix(name, "-") || strings.HasSuffix(name, ".") {
		return newValidationError(
			domain.ErrInvalidProject,
			"Invalid project name",
			"Project name cannot end with hyphen or dot",
			"project_name",
			"Remove trailing hyphens or dots",
		)
	}

	if strings.Contains(name, "--") || strings.Contains(name, "..") {
		return newValidationError(
			domain.ErrInvalidProject,
			"Invalid project name",
			"Project name cannot contain consecutive hyphens or dots",
			"project_name",
			"Use single hyphens or dots",
		)
	}

	lowerName := strings.ToLower(name)
	if slices.Contains(reservedProjectNames, lowerName) ||
		slices.Contains(reservedWindowsDeviceNames, lowerName) {
		return newValidationError(
			domain.ErrInvalidProject,
			"Reserved project name",
			fmt.Sprintf("'%s' is a reserved name", name),
			"project_name",
			"Choose a different project name",
		)
	}

	return nil
}

// ValidateBinaryName validates binary name.
func ValidateBinaryName(name string) error {
	if name == "" {
		return newValidationError(
			domain.ErrInvalidBinary,
			"Binary name is required",
			"Binary name cannot be empty",
			"binary_name",
			"Choose a valid binary name",
		)
	}

	if len(name) > 30 {
		return newValidationError(
			domain.ErrInvalidBinary,
			"Binary name too long",
			"Binary name must be 30 characters or less",
			"binary_name",
			"Use a shorter binary name",
		)
	}

	if !binaryNamePattern.MatchString(name) {
		return newValidationError(
			domain.ErrInvalidBinary,
			"Invalid binary name format",
			"Binary name must contain only letters, numbers, hyphens, and underscores",
			"binary_name",
			"Use alphanumeric characters with hyphens and underscores only",
		)
	}

	lowerName := strings.ToLower(name)

	reservedNames := make([]string, 0, len(reservedWindowsDeviceNames)+len(reservedBinaryNames))
	reservedNames = append(reservedNames, reservedWindowsDeviceNames...)

	reservedNames = append(reservedNames, reservedBinaryNames...)
	if slices.Contains(reservedNames, lowerName) {
		return newValidationError(
			domain.ErrInvalidBinary,
			"Reserved binary name",
			fmt.Sprintf("'%s' is a reserved binary name", name),
			"binary_name",
			"Choose a different binary name",
		)
	}

	return nil
}

// ValidateMainPath validates main path.
func ValidateMainPath(path string) error {
	if path == "" {
		return newValidationError(
			domain.ErrInvalidMainPath,
			"Main path is required",
			"Main path cannot be empty",
			"main_path",
			"Specify the main package path (e.g., '.', './cmd/app')",
		)
	}

	normalizedPath := strings.ReplaceAll(path, "\\", "/")

	if strings.HasPrefix(normalizedPath, "/") {
		return newValidationError(
			domain.ErrInvalidMainPath,
			"Absolute paths not allowed",
			"Main path must be a relative path, not an absolute path",
			"main_path",
			"Use relative paths like '.', './cmd/app', or './main'",
		)
	}

	if len(normalizedPath) > 200 {
		return newValidationError(
			domain.ErrInvalidMainPath,
			"Main path too long",
			"Main path must be 200 characters or less",
			"main_path",
			"Use a shorter path",
		)
	}

	if !mainPathPattern.MatchString(normalizedPath) {
		return newValidationError(
			domain.ErrInvalidMainPath,
			"Invalid main path format",
			"Main path must be a valid relative path",
			"main_path",
			"Use relative paths like '.', './cmd/app', or './main'",
		)
	}

	err := validatePathShellMetachars(normalizedPath)
	if err != nil {
		return err
	}

	if strings.Contains(normalizedPath, "..") {
		return newValidationError(
			domain.ErrInvalidMainPath,
			"Path traversal not allowed",
			"Main path contains path traversal sequences",
			"main_path",
			"Use relative paths without '..'",
		)
	}

	return validatePathComponents(normalizedPath)
}

func validatePathShellMetachars(path string) error {
	shellMetachars := []string{
		";", "&", "|", "<", ">", "`", "$", "(", ")", "{", "}", "!", "*", "?", "~",
	}
	for _, char := range shellMetachars {
		if strings.Contains(path, char) {
			return newValidationError(
				domain.ErrInvalidMainPath,
				"Shell metacharacters not allowed",
				fmt.Sprintf("Main path contains shell metacharacter '%s'", char),
				"main_path",
				"Use safe path characters only",
			)
		}
	}

	return nil
}

var reservedDirs = []string{"etc", "bin", "usr", "sbin", "var", "sys", "proc", "dev"}

// checkReservedNames checks if a path component is in a reserved list.
func checkReservedName(component string, reservedList []string, errorMsg, suggestion string) error {
	lowerComponent := strings.ToLower(component)
	if slices.Contains(reservedList, lowerComponent) {
		return newValidationError(
			domain.ErrInvalidMainPath,
			"Reserved directory",
			fmt.Sprintf("'%s' %s", component, errorMsg),
			"main_path",
			suggestion,
		)
	}

	return nil
}

func validatePathComponents(path string) error {
	components := strings.SplitSeq(path, "/")
	for component := range components {
		if component == "" {
			continue
		}

		err := checkReservedName(
			component,
			reservedWindowsDeviceNames,
			"is not allowed in path",
			"Remove invalid components from path",
		)
		if err != nil {
			return err
		}

		err = checkReservedName(
			component,
			reservedDirs,
			"is a reserved system directory",
			"Use project directories instead",
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateProjectDescription validates project description.
func ValidateProjectDescription(description string) error {
	if description == "" {
		return nil // Description is optional
	}

	trimmed := strings.TrimSpace(description)

	if len(trimmed) == 0 {
		return newValidationError(
			domain.ErrInvalidProjectDescription,
			"Project description cannot be empty",
			"Project description must contain text",
			"project_description",
			"Provide a meaningful description",
		)
	}

	if len(trimmed) > 255 {
		return newValidationError(
			domain.ErrInvalidProjectDescription,
			"Project description too long",
			"Project description must be 255 characters or less",
			"project_description",
			"Use a shorter description",
		)
	}

	if !projectDescriptionPattern.MatchString(trimmed) {
		return newValidationError(
			domain.ErrInvalidProjectDescription,
			"Invalid characters in project description",
			"Project description contains invalid characters",
			"project_description",
			"Use only printable characters",
		)
	}

	lowerDesc := strings.ToLower(trimmed)
	if strings.Contains(lowerDesc, "<script") || strings.Contains(lowerDesc, "javascript:") {
		return newValidationError(
			domain.ErrInvalidProjectDescription,
			"Invalid content in project description",
			"Project description contains script injection patterns",
			"project_description",
			"Remove script tags and javascript: prefixes",
		)
	}

	if strings.Contains(trimmed, "TODO") || strings.Contains(trimmed, "FIXME") {
		return newValidationError(
			domain.ErrInvalidProjectDescription,
			"Incomplete project description",
			"Project description contains TODO or FIXME markers",
			"project_description",
			"Complete the project description",
		)
	}

	if len(trimmed) > 10 && hasExcessiveWhitespace(trimmed) {
		return newValidationError(
			domain.ErrInvalidProjectDescription,
			"Suspicious whitespace pattern",
			"Project description contains too much whitespace",
			"project_description",
			"Use normal spacing in description",
		)
	}

	return nil
}

func hasExcessiveWhitespace(s string) bool {
	spaceCount := 0

	for _, r := range s {
		if r == ' ' {
			spaceCount++
		}
	}

	return spaceCount > len(s)/2
}

func ValidateVersion(version string) error {
	if version == "" {
		return newValidationError(
			domain.ErrInvalidVersion,
			"Version is required",
			"Version cannot be empty",
			"version",
			"Specify a version (e.g., v1.0.0)",
		)
	}

	semverPattern := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$`)
	gitDescribePattern := regexp.MustCompile(`^[a-f0-9]{7,40}(-dirty)?$`)

	if !semverPattern.MatchString(version) && !gitDescribePattern.MatchString(version) {
		return newValidationError(
			domain.ErrInvalidVersion,
			"Invalid version format",
			"Version must follow semantic versioning (e.g., v1.0.0) or be a git commit hash",
			"version",
			"Use semantic versioning (e.g., v1.0.0, v1.0.0-alpha.1)",
		)
	}

	return nil
}

// ValidateGitBranch validates Git branch name.
func ValidateGitBranch(branch string) error {
	if branch == "" {
		return newValidationError(
			domain.ErrInvalidGitBranch,
			"Git branch name is required",
			"Git branch name cannot be empty",
			"git_branch",
			"Specify a Git branch name",
		)
	}

	branchPattern := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9/_-]{0,252}$`)
	if !branchPattern.MatchString(branch) {
		return newValidationError(
			domain.ErrInvalidGitBranch,
			"Invalid Git branch name format",
			"Git branch name must contain only letters, numbers, hyphens, underscores, and forward slashes",
			"git_branch",
			"Use valid Git branch naming conventions",
		)
	}

	invalidBranchNames := []string{
		"HEAD", "master", "main", "develop", "feature", "release", "hotfix",
		"MERGE_HEAD", "ORIG_HEAD", "FETCH_HEAD",
	}

	lowerBranch := strings.ToLower(branch)
	if slices.Contains(invalidBranchNames, lowerBranch) {
		return newValidationError(
			domain.ErrInvalidGitBranch,
			"Reserved Git branch name",
			fmt.Sprintf("'%s' is a reserved branch name", branch),
			"git_branch",
			"Use a descriptive branch name",
		)
	}

	return nil
}

// ValidateGitTag validates Git tag name.
func ValidateGitTag(tag string) error {
	if tag == "" {
		return newValidationError(
			domain.ErrInvalidGitTag,
			"Git tag name is required",
			"Git tag name cannot be empty",
			"git_tag",
			"Specify a Git tag name",
		)
	}

	tagPattern := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,254}$`)
	if !tagPattern.MatchString(tag) {
		return newValidationError(
			domain.ErrInvalidGitTag,
			"Invalid Git tag format",
			"Git tag name must contain only letters, numbers, dots, hyphens, and underscores",
			"git_tag",
			"Use valid Git tag naming conventions",
		)
	}

	return nil
}

var buildTagPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*$`)

// ValidateBuildTags validates build tags.
func ValidateBuildTags(tags []string) error {
	if len(tags) == 0 {
		return nil // No build tags is valid
	}

	for i, tag := range tags {
		err := validateSingleBuildTag(tag, i)
		if err != nil {
			return err
		}
	}

	return checkDuplicateBuildTags(tags)
}

func validateSingleBuildTag(tag string, index int) error {
	if tag == "" {
		return newBuildTagValidationError(
			fmt.Sprintf("Build tag at index %d is empty", index),
			"Remove empty build tags",
			tag,
		)
	}

	if !buildTagPattern.MatchString(tag) {
		return newBuildTagValidationError(
			fmt.Sprintf("Build tag '%s' is invalid", tag),
			"Use alphanumeric characters and underscores only",
			tag,
		)
	}

	if len(tag) > 50 {
		return newBuildTagValidationError(
			fmt.Sprintf("Build tag '%s' exceeds 50 characters", tag),
			"Use shorter build tags",
			tag,
		)
	}

	return nil
}

func checkDuplicateBuildTags(tags []string) error {
	seen := make(map[string]bool)
	for _, tag := range tags {
		if seen[tag] {
			return newBuildTagValidationError(
				fmt.Sprintf("Build tag '%s' is specified multiple times", tag),
				"Remove duplicate build tags",
				tag,
			)
		}

		seen[tag] = true
	}

	return nil
}

func newBuildTagValidationError(details, suggestion, tag string) error {
	return domain.NewValidationError(
		domain.ErrInvalidBuildTag,
		getBuildTagErrorTitle(tag),
		details,
	).WithField("build_tags").WithSuggestion(suggestion)
}

func getBuildTagErrorTitle(tag string) string {
	if tag == "" {
		return "Build tag cannot be empty"
	}

	if !buildTagPattern.MatchString(tag) {
		return "Invalid build tag format"
	}

	return "Build tag too long"
}

// ValidatePort validates port number.
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return newValidationError(
			domain.ErrInvalidPort,
			"Invalid port number",
			"Port number must be between 1 and 65535",
			"port",
			"Use a valid port number (1-65535)",
		)
	}

	if slices.Contains(restrictedPorts, port) {
		return newValidationError(
			domain.ErrInvalidPort,
			"Restricted port number",
			fmt.Sprintf("Port %d is restricted and shouldn't be used by applications", port),
			"port",
			"Use a non-privileged port (1024-65535)",
		)
	}

	return nil
}

var restrictedPorts = []int{
	1, 7, 9, 11, 13, 17, 19, 20, 21, 22, 23, 25,
	53, 67, 68, 69, 80, 110, 123, 135, 137, 138, 139,
	143, 161, 162, 179, 389, 443, 445, 512, 513, 514,
	515, 520, 521, 522, 523, 524, 525, 526, 527, 528,
	529, 530, 531, 532, 533, 534, 535, 536, 537, 538,
	539, 540, 541, 542, 543, 544, 545, 546, 547, 548,
	549, 550, 551, 552, 553, 554, 555, 556, 557, 558,
	559, 560, 561, 562, 563, 564, 565, 566, 567, 568,
	569, 570, 571, 572, 573, 574, 575, 576, 577, 578,
	579, 580, 581, 582, 583, 584, 585, 586, 587, 588,
	589, 590, 591, 592, 593, 594, 595, 596, 597, 598,
	599, 600, 601, 602, 603, 604, 605, 606, 607, 608,
	609, 610, 611, 612, 613, 614, 615, 616, 617, 618,
	619, 620, 621, 622, 623, 624, 625, 626, 627, 628,
	629, 630, 631, 632, 633, 634, 635, 636, 637, 638,
	639, 640, 641, 642, 643, 644, 645, 646, 647, 648,
	649, 650, 651, 652, 653, 654, 655,
}
