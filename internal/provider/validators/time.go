// Copyright (c) codewolf.dev
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var timeHHMMSSRegex = regexp.MustCompile(
	`^([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`,
)

type timeHHMMSSValidator struct{}

func TimeHHMMSS() validator.String {
	return timeHHMMSSValidator{}
}

func (v timeHHMMSSValidator) Description(_ context.Context) string {
	return "Time in HH:MM:SS 24-hour format"
}

func (v timeHHMMSSValidator) MarkdownDescription(_ context.Context) string {
	return "Time in **HH:MM:SS** 24-hour format"
}

func (v timeHHMMSSValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !timeHHMMSSRegex.MatchString(req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid time format",
			"Expected time in HH:MM:SS format (e.g. 16:00:00).",
		)
	}
}
