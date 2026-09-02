package scan

func clippySeverity(level string) Severity {
	switch level {
	case "error":
		return SeverityHigh
	case "warning":
		return SeverityWarning
	case "note", "help":
		return SeverityInfo
	default:
		return SeverityWarning
	}
}

func clippyApplicability(value *string) Applicability {
	if value == nil {
		return ApplicabilityUnspecified
	}
	switch *value {
	case "MachineApplicable":
		return ApplicabilityMachine
	case "MaybeIncorrect":
		return ApplicabilityMaybe
	case "HasPlaceholders":
		return ApplicabilityHasPlaceholders
	case "Unspecified":
		return ApplicabilityUnspecified
	default:
		return ApplicabilityManual
	}
}

func lessSafeApplicability(left, right Applicability) Applicability {
	if applicabilityRank(right) > applicabilityRank(left) {
		return right
	}
	return left
}

func applicabilityRank(value Applicability) int {
	switch value {
	case ApplicabilityMachine:
		return 1
	case ApplicabilityMaybe:
		return 2
	case ApplicabilityHasPlaceholders:
		return 3
	case ApplicabilityUnspecified:
		return 4
	case ApplicabilityManual:
		return 5
	default:
		return 6
	}
}
