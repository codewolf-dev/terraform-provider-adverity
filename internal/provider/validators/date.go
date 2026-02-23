// Copyright codewolf.dev 2025, 0
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var dateYYYYMMDDRegex = regexp.MustCompile(
	`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`,
)

type dateYYYYMMDDValidator struct{}

func DateYYYYMMDD() validator.String {
	return dateYYYYMMDDValidator{}
}

func (v dateYYYYMMDDValidator) Description(_ context.Context) string {
	return "Date in YYYY-MM-DD format"
}

func (v dateYYYYMMDDValidator) MarkdownDescription(_ context.Context) string {
	return "Date in **YYYY-MM-DD** format"
}

func (v dateYYYYMMDDValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !dateYYYYMMDDRegex.MatchString(req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid date format",
			"Expected date in YYYY-MM-DD format (e.g. 2025-05-01).",
		)
	}
}
