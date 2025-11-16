package domain

// Architecture represents supported CPU architectures
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type Architecture string

const (
	ArchitectureAMD64   Architecture = "amd64"   // 64-bit x86
	ArchitectureARM64   Architecture = "arm64"   // 64-bit ARM
	Architecture386     Architecture = "386"     // 32-bit x86
	ArchitectureARM     Architecture = "arm"     // 32-bit ARM
	ArchitecturePPC64   Architecture = "ppc64"   // 64-bit PowerPC (big endian)
	ArchitecturePPC64LE Architecture = "ppc64le" // 64-bit PowerPC (little endian)
	ArchitectureS390X   Architecture = "s390x"   // IBM System z
	ArchitectureMIPS    Architecture = "mips"    // 32-bit MIPS (big endian)
	ArchitectureMIPSLE  Architecture = "mipsle"  // 32-bit MIPS (little endian)
)

// Architecture metadata - generated from TypeSpec invariants
type architectureMetadata struct {
	supportedByAllPlatforms bool
	is64Bit               bool
	goSupport             string
}

var architectureMetadata = map[Architecture]architectureMetadata{
	ArchitectureAMD64: {
		supportedByAllPlatforms: true,
		is64Bit:               true,
		goSupport:             "stable",
	},
	ArchitectureARM64: {
		supportedByAllPlatforms: true,
		is64Bit:               true,
		goSupport:             "stable",
	},
	Architecture386: {
		supportedByAllPlatforms: true,
		is64Bit:               false,
		goSupport:             "stable",
	},
	ArchitectureARM: {
		supportedByAllPlatforms: true,
		is64Bit:               false,
		goSupport:             "stable",
	},
	ArchitecturePPC64: {
		supportedByAllPlatforms: false,
		is64Bit:               true,
		goSupport:             "stable",
	},
	ArchitecturePPC64LE: {
		supportedByAllPlatforms: false,
		is64Bit:               true,
		goSupport:             "stable",
	},
	ArchitectureS390X: {
		supportedByAllPlatforms: false,
		is64Bit:               true,
		goSupport:             "stable",
	},
	ArchitectureMIPS: {
		supportedByAllPlatforms: false,
		is64Bit:               false,
		goSupport:             "stable",
	},
	ArchitectureMIPSLE: {
		supportedByAllPlatforms: false,
		is64Bit:               false,
		goSupport:             "stable",
	},
}

// IsValid returns true if Architecture is valid
func (a Architecture) IsValid() bool {
	_, exists := architectureMetadata[a]
	return exists
}

// String returns human-readable display name
func (a Architecture) String() string {
	switch a {
	case ArchitectureAMD64:
		return "64-bit x86"
	case ArchitectureARM64:
		return "64-bit ARM"
	case Architecture386:
		return "32-bit x86"
	case ArchitectureARM:
		return "32-bit ARM"
	case ArchitecturePPC64:
		return "64-bit PowerPC (big endian)"
	case ArchitecturePPC64LE:
		return "64-bit PowerPC (little endian)"
	case ArchitectureS390X:
		return "IBM System z"
	case ArchitectureMIPS:
		return "32-bit MIPS (big endian)"
	case ArchitectureMIPSLE:
		return "32-bit MIPS (little endian)"
	default:
		return string(a)
	}
}

// SupportedByAllPlatforms returns true if architecture is supported by all platforms
func (a Architecture) SupportedByAllPlatforms() bool {
	if meta, exists := architectureMetadata[a]; exists {
		return meta.supportedByAllPlatforms
	}
	return false
}

// Is64Bit returns true if architecture is 64-bit
func (a Architecture) Is64Bit() bool {
	if meta, exists := architectureMetadata[a]; exists {
		return meta.is64Bit
	}
	return false
}

// GoSupport returns Go support level for this architecture
func (a Architecture) GoSupport() string {
	if meta, exists := architectureMetadata[a]; exists {
		return meta.goSupport
	}
	return "unknown"
}

// ValidateArchitectures validates a slice of architectures
func ValidateArchitectures(architectures []Architecture) error {
	if len(architectures) == 0 {
		return fmt.Errorf("at least one architecture is required")
	}
	
	for _, arch := range architectures {
		if !arch.IsValid() {
			return fmt.Errorf("invalid architecture: %s", arch)
		}
	}
	
	return nil
}

// GetRecommendedArchitectures returns recommended architectures for common use cases
func GetRecommendedArchitectures() []Architecture {
	return []Architecture{ArchitectureAMD64, ArchitectureARM64}
}

// GetAllArchitectures returns all supported architectures
func GetAllArchitectures() []Architecture {
	return []Architecture{
		ArchitectureAMD64, ArchitectureARM64, Architecture386, ArchitectureARM,
		ArchitecturePPC64, ArchitecturePPC64LE, ArchitectureS390X,
		ArchitectureMIPS, ArchitectureMIPSLE,
	}
}