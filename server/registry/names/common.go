// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package names

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	// The format of a resource identifier.
	// This may be extended to include all characters that do not require escaping.
	// See https://aip.dev/122#resource-id-segments.
	identifier = `([A-Za-z0-9-.]+)`

	// The format of a custom revision tag.
	revisionTag = `([a-z0-9-]+)`
)

// The format of a custom resource identifier.
// User provided identifiers should be validated according to this format.
var customIdentifier = regexp.MustCompile(`^[a-z0-9-.]+$`)

// GenerateID generates a random resource ID.
func GenerateID() string {
	return uuid.New().String()[:8]
}

// validateID returns an error if the provided ID is invalid.
func validateID(id string) error {
	if !customIdentifier.MatchString(id) {
		return fmt.Errorf("invalid identifier %q: must match %q", id, customIdentifier)
	} else if _, err := uuid.Parse(id); err == nil {
		return fmt.Errorf("invalid identifier %q: must not match UUID format", id)
	} else if len(id) > 80 {
		return fmt.Errorf("invalid identifier %q: must be 80 characters or less", id)
	} else if strings.HasPrefix(id, "-") || strings.HasPrefix(id, ".") {
		return fmt.Errorf("invalid identifier %q: must begin with a number or letter", id)
	} else if strings.HasSuffix(id, "-") || strings.HasSuffix(id, ".") {
		return fmt.Errorf("invalid identifier %q: must end with a number or letter", id)
	}

	return nil
}

// Normalize is an idempotent operation for normalizing resource names and identifiers.
// Identifiers `a` and `b` should be considered equal if and only if normalize(a) == normalize(b).
func normalize(identifier string) string {
	return strings.ToLower(identifier)
}

// Location is included in resource names immediately following the project_id.
const Location = "global"

// Name is an interface that represents resource names.
type Name interface {
	String() string // all names have a string representation.
}
