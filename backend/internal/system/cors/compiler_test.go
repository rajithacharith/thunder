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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	yaml "gopkg.in/yaml.v3"
)

type CompilerTestSuite struct {
	suite.Suite
}

func TestCompilerTestSuite(t *testing.T) {
	suite.Run(t, new(CompilerTestSuite))
}

func (suite *CompilerTestSuite) TestCompileLiteralEntry() {
	rule, err := Compile(LiteralEntry{Value: "https://example.com"})
	suite.Require().NoError(err)
	assert.Equal(suite.T(), RuleLiteral, rule.Kind())
}

func (suite *CompilerTestSuite) TestCompileRegexEntry() {
	rule, err := Compile(RegexEntry{Pattern: `^https://.+$`})
	suite.Require().NoError(err)
	assert.Equal(suite.T(), RuleRegex, rule.Kind())
}

func (suite *CompilerTestSuite) TestCompileLiteralWhitespaceTrimmed() {
	rule, err := Compile(LiteralEntry{Value: "  https://example.com  "})
	suite.Require().NoError(err)
	m := NewMatcher([]OriginRule{rule})
	parsed, err := ParseOrigin("https://example.com")
	suite.Require().NoError(err)
	allow, _ := m.Match(parsed)
	assert.True(suite.T(), allow)
}

func (suite *CompilerTestSuite) TestCompileNullLiteral() {
	rule, err := Compile(LiteralEntry{Value: "null"})
	suite.Require().NoError(err)
	m := NewMatcher([]OriginRule{rule})
	allowNull, _ := m.Match(ParseResult{Raw: "null", IsNull: true})
	assert.True(suite.T(), allowNull)
	parsed, err := ParseOrigin("https://example.com")
	suite.Require().NoError(err)
	allowOrigin, _ := m.Match(parsed)
	assert.False(suite.T(), allowOrigin)
}

func (suite *CompilerTestSuite) TestCompileWildcardLiteralRejected() {
	_, err := Compile(LiteralEntry{Value: "*"})
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrWildcardLiteral))
}

func (suite *CompilerTestSuite) TestIsRegexAnchored() {
	cases := []struct {
		pattern  string
		anchored bool
	}{
		{`^https://example\.com$`, true},
		{`\Ahttps://example\.com\z`, true},
		{`^https://example\.com\z`, true},
		{`\Ahttps://example\.com$`, true},
		{`https://example\.com$`, false},
		{`^https://example\.com`, false},
		{`https://example\.com`, false},
		{`.*\.example\.com`, false},
	}
	for _, c := range cases {
		assert.Equal(suite.T(), c.anchored, IsRegexAnchored(c.pattern), c.pattern)
	}
}

func (suite *CompilerTestSuite) TestCompileEmptyLiteralRejected() {
	_, err := Compile(LiteralEntry{Value: "   "})
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrEmptyEntry))
}

func (suite *CompilerTestSuite) TestCompileInvalidLiteralRejected() {
	_, err := Compile(LiteralEntry{Value: "not-a-url"})
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidLiteral))
}

func (suite *CompilerTestSuite) TestCompileEmptyRegexRejected() {
	_, err := Compile(RegexEntry{Pattern: ""})
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrEmptyEntry))
}

func (suite *CompilerTestSuite) TestCompileInvalidRegexRejected() {
	_, err := Compile(RegexEntry{Pattern: "([unterminated"})
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidRegex))
}

type unknownEntry struct{}

func (unknownEntry) isOriginEntry() {}

func (suite *CompilerTestSuite) TestCompileUnknownEntryTypeRejected() {
	_, err := Compile(unknownEntry{})
	suite.Require().Error(err)
}

func (suite *CompilerTestSuite) TestCompileAllEmptyInput() {
	rules, err := CompileAll(nil)
	suite.Require().NoError(err)
	assert.Nil(suite.T(), rules)
}

func (suite *CompilerTestSuite) TestCompileAllPreservesOrder() {
	rules, err := CompileAll([]Entry{
		LiteralEntry{Value: "https://a.com"},
		RegexEntry{Pattern: `^https://b\.com$`},
		LiteralEntry{Value: "https://c.com"},
	})
	suite.Require().NoError(err)
	suite.Require().Len(rules, 3)
	assert.Equal(suite.T(), RuleLiteral, rules[0].Kind())
	assert.Equal(suite.T(), RuleRegex, rules[1].Kind())
	assert.Equal(suite.T(), RuleLiteral, rules[2].Kind())
}

func (suite *CompilerTestSuite) TestCompileAllFailsFastWithIndex() {
	_, err := CompileAll([]Entry{
		LiteralEntry{Value: "https://ok.com"},
		RegexEntry{Pattern: "([bad"},
	})
	suite.Require().Error(err)
	assert.Contains(suite.T(), err.Error(), "allowed_origins[1]")
}

func (suite *CompilerTestSuite) TestUnmarshalYAMLLiteralEntries() {
	doc := []byte(`
- https://example.com
- https://other.com
`)
	var entries OriginEntries
	suite.Require().NoError(yaml.Unmarshal(doc, &entries))
	suite.Require().Len(entries, 2)
	assert.IsType(suite.T(), LiteralEntry{}, entries[0])
	assert.IsType(suite.T(), LiteralEntry{}, entries[1])
}

func (suite *CompilerTestSuite) TestUnmarshalYAMLRegexEntries() {
	doc := []byte(`
- regex: '^https://[a-z]+\.example\.com$'
`)
	var entries OriginEntries
	suite.Require().NoError(yaml.Unmarshal(doc, &entries))
	suite.Require().Len(entries, 1)
	r, ok := entries[0].(RegexEntry)
	suite.Require().True(ok)
	assert.Equal(suite.T(), `^https://[a-z]+\.example\.com$`, r.Pattern)
}

func (suite *CompilerTestSuite) TestUnmarshalYAMLMixedEntries() {
	doc := []byte(`
- https://example.com
- regex: '^https://.+\.staging\.example\.com$'
- "null"
`)
	var entries OriginEntries
	suite.Require().NoError(yaml.Unmarshal(doc, &entries))
	suite.Require().Len(entries, 3)
	assert.IsType(suite.T(), LiteralEntry{}, entries[0])
	assert.IsType(suite.T(), RegexEntry{}, entries[1])
	assert.IsType(suite.T(), LiteralEntry{}, entries[2])
}

func (suite *CompilerTestSuite) TestUnmarshalYAMLNonSequenceRejected() {
	doc := []byte(`foo: bar`)
	var entries OriginEntries
	err := yaml.Unmarshal(doc, &entries)
	suite.Require().Error(err)
}

func (suite *CompilerTestSuite) TestUnmarshalYAMLRegexMissingFieldRejected() {
	doc := []byte(`
- pattern: foo
`)
	var entries OriginEntries
	err := yaml.Unmarshal(doc, &entries)
	suite.Require().Error(err)
}

func (suite *CompilerTestSuite) TestUnmarshalYAMLUnsupportedNodeRejected() {
	doc := []byte(`
- - nested
`)
	var entries OriginEntries
	err := yaml.Unmarshal(doc, &entries)
	suite.Require().Error(err)
}
