package domain

// Platform represents supported target platforms
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type Platform string

const (
	PlatformLinux   Platform = "linux"   // Linux
	PlatformDarwin  Platform = "darwin"  // macOS
	PlatformWindows Platform = "windows" // Windows
	PlatformFreeBSD Platform = "freebsd" // FreeBSD
	PlatformOpenBSD Platform = "openbsd" // OpenBSD
	PlatformNetBSD  Platform = "netbsd"  // NetBSD
)

// Platform metadata - generated from TypeSpec invariants
type platformMetadata struct {
	architectures   []Architecture
	isWindowsBased bool
	isUnixLike     bool
	supportsCGO    bool
}

var platformMetadata = map[Platform]platformMetadata{
	PlatformLinux: {
		architectures:   []Architecture{ArchitectureAMD64, ArchitectureARM64, ArchitectureARM, Architecture386},
		isWindowsBased:  false,
		isUnixLike:      true,
		supportsCGO:     true,
	},
	PlatformDarwin: {
		architectures:   []Architecture{ArchitectureAMD64, ArchitectureARM64},
		isWindowsBased:  false,
		isUnixLike:      true,
		supportsCGO:     true,
	},
	PlatformWindows: {
		architectures:   []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386},
		isWindowsBased:  true,
		isUnixLike:      false,
		supportsCGO:     true,
	},
	PlatformFreeBSD: {
		architectures:   []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386},
		isWindowsBased:  false,
		isUnixLike:      true,
		supportsCGO:     true,
	},
	PlatformOpenBSD: {
		architectures:   []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386},
		isWindowsBased:  false,
		isUnixLike:      true,
		supportsCGO:     true,
	},
	PlatformNetBSD: {
		architectures:   []Architecture{ArchitectureAMD64, ArchitectureARM64, Architecture386},
		isWindowsBased:  false,
		isUnixLike:      true,
		supportsCGO:     true,
	},
}

// IsValid returns true if Platform is valid
func (p Platform) IsValid() bool {
	_, exists := platformMetadata[p]
	return exists
}

// String returns human-readable display name
func (p Platform) String() string {
	switch p {
	case PlatformLinux:
		return "Linux"
	case PlatformDarwin:
		return "macOS"
	case PlatformWindows:
		return "Windows"
	case PlatformFreeBSD:
		return "FreeBSD"
	case PlatformOpenBSD:
		return "OpenBSD"
	case PlatformNetBSD:
		return "NetBSD"
	default:
		return string(p)
	}
}

// Architectures returns supported architectures for this platform
func (p Platform) Architectures() []Architecture {
	if meta, exists := platformMetadata[p]; exists {
		return meta.architectures
	}
	return []Architecture{ArchitectureAMD64, ArchitectureARM64}
}

// IsWindowsBased returns true if platform is Windows-based
func (p Platform) IsWindowsBased() bool {
	if meta, exists := platformMetadata[p]; exists {
		return meta.isWindowsBased
	}
	return false
}

// IsUnixLike returns true if platform is Unix-like
func (p Platform) IsUnixLike() bool {
	if meta, exists := platformMetadata[p]; exists {
		return meta.isUnixLike
	}
	return false
}

// SupportsCGO returns true if platform supports CGO
func (p Platform) SupportsCGO() bool {
	if meta, exists := platformMetadata[p]; exists {
		return meta.supportsCGO
	}
	return true
}

// ValidatePlatforms validates a slice of platforms for compatibility
func ValidatePlatforms(platforms []Platform) error {
	if len(platforms) == 0 {
		return fmt.Errorf("at least one platform is required")
	}
	
	for _, platform := range platforms {
		if !platform.IsValid() {
			return fmt.Errorf("invalid platform: %s", platform)
		}
	}
	
	return nil
}

// ValidatePlatformArchCompatibility validates that all architectures are compatible with selected platforms
func ValidatePlatformArchCompatibility(platforms []Platform, architectures []Architecture) error {
	if len(platforms) == 0 || len(architectures) == 0 {
		return fmt.Errorf("both platforms and architectures must be specified")
	}
	
	for _, platform := range platforms {
		platformArchs := platform.Architectures()
		for _, arch := range architectures {
			compatible := false
			for _, platformArch := range platformArchs {
				if arch == platformArch {
					compatible = true
					break
				}
			}
			if !compatible {
				return fmt.Errorf("architecture %s is not supported on platform %s", arch, platform)
			}
		}
	}
	
	return nil
}