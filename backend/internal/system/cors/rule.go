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

import "regexp"

// RuleKind identifies whether a compiled rule was sourced from a literal or
// a regex configuration entry. It is intended for diagnostics and logging.
type RuleKind int

const (
	// RuleLiteral indicates a rule compiled from a bare-string YAML entry.
	RuleLiteral RuleKind = iota
	// RuleRegex indicates a rule compiled from a regex YAML entry.
	RuleRegex
)

// OriginRule is the discriminated-union type produced by Compile. The matcher
// disassembles compiled rules into kind-specific data structures
// (canonical-key map for literals, regex slice for regex rules) so request-time
// matching avoids interface dispatch and per-request canonicalization. Only
// Kind() is needed at runtime; matching itself is performed by the Matcher.
type OriginRule interface {
	// Kind reports the rule's source kind.
	Kind() RuleKind
}

// LiteralRule matches a single canonicalized origin. The canonical form uses
// the lowercased scheme + lowercased host (with IDN labels Punycode-encoded
// and any trailing dot stripped); IPv6 hosts are bracketed. The port is
// preserved verbatim, so a portless origin and the same origin with an
// explicit default port (e.g. "https://example.com" vs.
// "https://example.com:443") remain distinct rules — operators that want both
// allowed must list each entry. The "null" origin is represented by isNull;
// such a rule matches only inputs whose IsNull flag is set.
type LiteralRule struct {
	canonical string
	isNull    bool
}

// Kind reports RuleLiteral.
func (r LiteralRule) Kind() RuleKind { return RuleLiteral }

// RegexRule matches the raw request Origin header against an operator-supplied
// RE2 pattern. The regex sees the raw header byte for byte after only the
// parse gate; no canonicalization or transformation is applied on the regex
// path. Operators own pattern correctness, including any anchoring required
// for full-input match (\A...\z).
type RegexRule struct {
	re *regexp.Regexp
}

// Kind reports RuleRegex.
func (r RegexRule) Kind() RuleKind { return RuleRegex }
