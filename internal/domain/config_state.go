package domain

import "fmt"

// ConfigState represents configuration lifecycle states
// Generated from TypeSpec specification - DO NOT MODIFY MANUALLY
type ConfigState string

const (
	ConfigStateDraft      ConfigState = "draft"      // Draft
	ConfigStateValid      ConfigState = "valid"      // Valid
	ConfigStateInvalid    ConfigState = "invalid"    // Invalid
	ConfigStateProcessing ConfigState = "processing" // Processing
	ConfigStateGenerated  ConfigState = "generated"  // Generated
)

// ConfigState metadata - generated from TypeSpec invariants
type configStateMeta struct {
	description       string
	isFinal           bool
	allowsValidation  bool
	allowsGeneration  bool
}

var configStateMetaMap = map[ConfigState]configStateMeta{
	ConfigStateDraft: {
		description:      "Configuration is being created or modified",
		isFinal:          false,
		allowsValidation:  true,
		allowsGeneration:  false,
	},
	ConfigStateValid: {
		description:      "Configuration is valid and ready",
		isFinal:          false,
		allowsValidation:  true,
		allowsGeneration:  true,
	},
	ConfigStateInvalid: {
		description:      "Configuration has validation errors",
		isFinal:          false,
		allowsValidation:  true,
		allowsGeneration:  false,
	},
	ConfigStateProcessing: {
		description:      "Configuration is being processed",
		isFinal:          false,
		allowsValidation:  false,
		allowsGeneration:  false,
	},
	ConfigStateGenerated: {
		description:      "Configuration has been generated successfully",
		isFinal:          true,
		allowsValidation:  true,
		allowsGeneration:  false,
	},
}

// IsValid returns true if ConfigState is valid
func (cs ConfigState) IsValid() bool {
	_, exists := configStateMetaMap[cs]
	return exists
}

// String returns human-readable display name
func (cs ConfigState) String() string {
	switch cs {
	case ConfigStateDraft:
		return "Draft"
	case ConfigStateValid:
		return "Valid"
	case ConfigStateInvalid:
		return "Invalid"
	case ConfigStateProcessing:
		return "Processing"
	case ConfigStateGenerated:
		return "Generated"
	default:
		return string(cs)
	}
}

// Description returns the description for this state
func (cs ConfigState) Description() string {
	if meta, exists := configStateMetaMap[cs]; exists {
		return meta.description
	}
	return ""
}

// IsFinal returns true if this is a final state
func (cs ConfigState) IsFinal() bool {
	if meta, exists := configStateMetaMap[cs]; exists {
		return meta.isFinal
	}
	return false
}

// AllowsValidation returns true if validation is allowed in this state
func (cs ConfigState) AllowsValidation() bool {
	if meta, exists := configStateMetaMap[cs]; exists {
		return meta.allowsValidation
	}
	return false
}

// AllowsGeneration returns true if generation is allowed in this state
func (cs ConfigState) AllowsGeneration() bool {
	if meta, exists := configStateMetaMap[cs]; exists {
		return meta.allowsGeneration
	}
	return false
}

// ValidateConfigState validates a configuration state
func ValidateConfigState(state ConfigState) error {
	if !state.IsValid() {
		return NewValidationError(
			ErrInvalidConfigState,
			"Invalid configuration state",
			fmt.Sprintf("'%s' is not a valid configuration state", state),
		)
	}
	return nil
}

// GetInitialConfigState returns the initial state for new configurations
func GetInitialConfigState() ConfigState {
	return ConfigStateDraft
}

// GetAllConfigStates returns all available configuration states
func GetAllConfigStates() []ConfigState {
	return []ConfigState{
		ConfigStateDraft, ConfigStateValid, ConfigStateInvalid,
		ConfigStateProcessing, ConfigStateGenerated,
	}
}