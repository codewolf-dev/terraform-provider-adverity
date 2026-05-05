// Copyright codewolf.dev 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// SplitImportParts splits a colon-separated import ID into exactly expectedParts parts.
// It adds a diagnostic error and returns nil if the format is invalid.
func SplitImportParts(id string, expectedParts int, format string, diagnostics *diag.Diagnostics) []string {
	parts := strings.Split(id, ":")
	if len(parts) != expectedParts {
		diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format: %s, got: %q", format, id),
		)
		return nil
	}
	return parts
}

// ParseImportPartInt parses a single string part of an import ID as an int64.
// It adds a diagnostic error and returns 0 if the value is not a valid integer.
func ParseImportPartInt(value string, fieldName string, diagnostics *diag.Diagnostics) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("%s must be a number, got: %q", fieldName, value),
		)
		return 0
	}
	return parsed
}
