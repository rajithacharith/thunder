/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cors

import (
	"fmt"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// Entry is the discriminated-union type produced by YAML decoding. It carries
// either a literal allowed-origin string or a regex pattern. Compile turns an
// Entry into the corresponding compiled OriginRule.
type Entry interface {
	isOriginEntry()
}

// LiteralEntry is the bare-string YAML form, e.g. "https://example.com" or
// the special-case "null".
type LiteralEntry struct {
	Value string
}

func (LiteralEntry) isOriginEntry() {}

// RegexEntry is the object YAML form, e.g. { regex: "\\Ahttps://..." }.
type RegexEntry struct {
	Pattern string
}

func (RegexEntry) isOriginEntry() {}

// OriginEntries is the slice wrapper that carries the heterogeneous YAML
// schema for cors.allowed_origins. Custom YAML unmarshaling on this type
// dispatches between the two entry forms.
type OriginEntries []Entry

// UnmarshalYAML decodes a YAML sequence whose elements are either scalar
// strings (LiteralEntry) or mappings of the shape { regex: "..." }
// (RegexEntry). Anything else is rejected at decode time.
func (e *OriginEntries) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("cors: allowed_origins must be a list, got %v", nodeKindString(node.Kind))
	}
	out := make(OriginEntries, 0, len(node.Content))
	for i, child := range node.Content {
		switch child.Kind {
		case yaml.ScalarNode:
			out = append(out, LiteralEntry{Value: child.Value})
		case yaml.MappingNode:
			var obj struct {
				Regex string `yaml:"regex"`
			}
			if err := child.Decode(&obj); err != nil {
				return fmt.Errorf("cors: allowed_origins[%d]: %w", i, err)
			}
			if obj.Regex == "" {
				return fmt.Errorf("cors: allowed_origins[%d]: regex object missing 'regex' field", i)
			}
			out = append(out, RegexEntry{Pattern: obj.Regex})
		default:
			return fmt.Errorf("cors: allowed_origins[%d]: entry must be a string or { regex: ... } object", i)
		}
	}
	*e = out
	return nil
}

// nodeKindString renders a yaml.Node kind for diagnostics.
func nodeKindString(k yaml.Kind) string {
	switch k {
	case yaml.DocumentNode:
		return "document"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.MappingNode:
		return "mapping"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"
	default:
		return "unknown"
	}
}

// Compile turns one Entry into a compiled OriginRule. Literal entries are
// gated through ParseOrigin and (for non-null entries) canonicalized; regex
// entries are compiled via Go's RE2 engine without any additional validation.
// Operator-supplied regex patterns are taken as-is.
func Compile(entry Entry) (OriginRule, error) {
	switch e := entry.(type) {
	case LiteralEntry:
		return compileLiteral(e.Value)
	case RegexEntry:
		return compileRegex(e.Pattern)
	default:
		return nil, fmt.Errorf("cors: unknown entry type %T", entry)
	}
}

// CompileAll compiles a slice of entries in declaration order. It fails fast
// on the first invalid entry, reporting the index and underlying cause so the
// operator can locate the bad entry in deployment.yaml. Empty input yields a
// nil slice with no error.
func CompileAll(entries []Entry) ([]OriginRule, error) {
	if len(entries) == 0 {
		return nil, nil
	}
	out := make([]OriginRule, 0, len(entries))
	for i, e := range entries {
		rule, err := Compile(e)
		if err != nil {
			return nil, fmt.Errorf("cors: allowed_origins[%d]: %w", i, err)
		}
		out = append(out, rule)
	}
	return out, nil
}

// compileLiteral builds a LiteralRule from a YAML bare-string value.
func compileLiteral(value string) (OriginRule, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, fmt.Errorf("%w: literal value is empty", ErrEmptyEntry)
	}
	if trimmed == "*" {
		return nil, fmt.Errorf("%w: list explicit origins or use a regex entry", ErrWildcardLiteral)
	}
	if trimmed == "null" {
		return LiteralRule{isNull: true}, nil
	}
	if _, err := ParseOrigin(trimmed); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidLiteral, err)
	}
	canonical, err := Canonicalize(trimmed)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidLiteral, err)
	}
	return LiteralRule{canonical: canonical}, nil
}

// compileRegex builds a RegexRule from a YAML regex object's pattern. The
// pattern is compiled by Go's RE2 engine; no further "safety" checks are
// applied — operator owns the pattern.
func compileRegex(pattern string) (OriginRule, error) {
	if pattern == "" {
		return nil, fmt.Errorf("%w: regex pattern is empty", ErrEmptyEntry)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidRegex, err)
	}
	return RegexRule{re: re}, nil
}

// IsRegexAnchored reports whether the given pattern starts with a
// start-of-input anchor (^ or \A) and ends with an end-of-input anchor ($ or
// \z). A pattern lacking either anchor permits substring matches and almost
// always allows far more origins than the operator intended; callers should
// log a warning at boot for unanchored patterns. The check is intentionally
// syntactic — alternation patterns like "(^a|^b)$" are flagged as a false
// positive, which is acceptable given the diagnostic-only intent.
func IsRegexAnchored(pattern string) bool {
	starts := strings.HasPrefix(pattern, "^") || strings.HasPrefix(pattern, `\A`)
	ends := strings.HasSuffix(pattern, "$") || strings.HasSuffix(pattern, `\z`)
	return starts && ends
}
