// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apex/log"
	lint "github.com/apigee/registry/cmd/registry/plugins/linter"
	"github.com/apigee/registry/rpc"
	"gopkg.in/yaml.v3"
)

type DescriptionField struct {
	StartLineNumber   int
	StartColumnNumber int
	EndLineNumber     int
	EndColumnNumber   int
	Description       string
}

// sampleOpenApiLinterRunner implements the LinterRunner interface for the sample OpenAPI linter.
type sampleOpenApiLinterRunner struct{}

const descriptionLessThan1000CharsRuleId = "description-less-than-1000-chars"
const descriptionContainsNoTagsRuleId = "description-contains-no-tags"

func (linter *sampleOpenApiLinterRunner) Run(req *rpc.LinterRequest) (*rpc.LinterResponse, error) {
	lintFiles := make([]*rpc.LintFile, 0)

	// Traverse the files in the directory
	err := filepath.Walk(req.GetSpecDirectory(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Async API and Open API specs are YAML files
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}

		// Execute the linter.
		problems, err := lintFile(path, req.GetRuleIds())
		if err != nil {
			return err
		}

		// Formulate the response.
		lintFiles = append(lintFiles, &rpc.LintFile{
			FilePath: path,
			Problems: problems,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &rpc.LinterResponse{
		Lint: &rpc.Lint{
			Name:  "registry-lint-openapi-sample",
			Files: lintFiles,
		},
	}, nil
}

func lintFile(specPath string, ruleIds []string) ([]*rpc.LintProblem, error) {
	specFile, err := ioutil.ReadFile(specPath)
	if err != nil {
		return nil, err
	}

	var parsedNode yaml.Node
	err = yaml.Unmarshal(specFile, &parsedNode)
	if err != nil {
		log.Fatalf("Unmarshal node: %v", err)
	}

	problems := make([]*rpc.LintProblem, 0)
	for _, ruleId := range ruleIds {
		lintProblems, err := lintWithRule(&parsedNode, ruleId)
		if err != nil {
			log.Errorf("Error while linting %s", err)
			continue
		}

		problems = append(problems, lintProblems...)
	}

	return problems, nil
}

func getDescriptionsFromSpec(node *yaml.Node) []*DescriptionField {
	results := make([]*DescriptionField, 0)
	getDescriptionsFromSpecHelper(node, &results)
	return results
}

func getDescriptionsFromSpecHelper(node *yaml.Node, results *[]*DescriptionField) {
	contentLength := len(node.Content)
	for i := 0; i < contentLength; i++ {
		child := node.Content[i]
		if i-1 >= 0 && node.Content[i-1].Kind == yaml.ScalarNode && node.Content[i-1].Value == "description" {
			endLineNumber := child.Line + 1
			endLineColumn := 0
			if i+1 < contentLength {
				endLineNumber = node.Content[i+1].Line
				endLineColumn = node.Content[i+1].Column
			}

			*results = append(*results, &DescriptionField{
				StartLineNumber:   child.Line,
				StartColumnNumber: child.Column,
				EndLineNumber:     endLineNumber,
				EndColumnNumber:   endLineColumn,
				Description:       child.Value,
			})
		}

		if child.Kind != yaml.ScalarNode {
			getDescriptionsFromSpecHelper(child, results)
		}
	}
}

func enforceDescriptionLessThan1000Chars(descriptions *[]*DescriptionField) []*rpc.LintProblem {
	problems := make([]*rpc.LintProblem, 0)
	for _, description := range *descriptions {
		if len(description.Description) >= 1000 {
			problems = append(problems, &rpc.LintProblem{
				Message:    "Description field should be less than 1000 chars.",
				RuleId:     descriptionLessThan1000CharsRuleId,
				Suggestion: "Ensure that your description field is less than 1000 chars in length.",
				Location: &rpc.LintLocation{
					StartPosition: &rpc.LintPosition{
						LineNumber:   int32(description.StartLineNumber),
						ColumnNumber: int32(description.StartColumnNumber),
					},
					EndPosition: &rpc.LintPosition{
						LineNumber:   int32(description.EndLineNumber),
						ColumnNumber: int32(description.EndColumnNumber),
					},
				},
			})
		}
	}
	return problems
}

func enforceDescriptionContainsNoTagsRuleId(descriptions *[]*DescriptionField) []*rpc.LintProblem {
	problems := make([]*rpc.LintProblem, 0)
	for _, description := range *descriptions {
		r, err := regexp.Compile(".*<[^>]*>.*")
		if err != nil {
			continue
		}
		if r.MatchString(description.Description) {
			problems = append(problems, &rpc.LintProblem{
				Message:    "Description field should not contain any tags.",
				RuleId:     descriptionContainsNoTagsRuleId,
				Suggestion: "Ensure that your description field does not contain any tags (regex <[^>]*>)",
				Location: &rpc.LintLocation{
					StartPosition: &rpc.LintPosition{
						LineNumber:   int32(description.StartLineNumber),
						ColumnNumber: int32(description.StartColumnNumber),
					},
					EndPosition: &rpc.LintPosition{
						LineNumber:   int32(description.EndLineNumber),
						ColumnNumber: int32(description.EndColumnNumber),
					},
				},
			})
		}
	}
	return problems
}

func lintWithRule(node *yaml.Node, ruleId string) ([]*rpc.LintProblem, error) {
	descriptions := getDescriptionsFromSpec(node)

	if ruleId == descriptionLessThan1000CharsRuleId {
		return enforceDescriptionLessThan1000Chars(&descriptions), nil
	} else if ruleId == descriptionContainsNoTagsRuleId {
		return enforceDescriptionContainsNoTagsRuleId(&descriptions), nil
	}

	return nil, fmt.Errorf("unsupported rule id %s", ruleId)
}

func main() {
	lint.Main(&sampleOpenApiLinterRunner{})
}
