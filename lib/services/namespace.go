/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"fmt"
	"regexp"

	"github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
)

// Check checks validity of all parameters and sets defaults
func (n *Namespace) CheckAndSetDefaults() error {
	if err := n.Metadata.CheckAndSetDefaults(); err != nil {
		return trace.Wrap(err)
	}
	isValid := IsValidNamespace(n.Metadata.Name)
	if !isValid {
		return trace.BadParameter("namespace %q is invalid", n.Metadata.Name)
	}

	return nil
}

const NamespaceSpecSchema = `{
  "type": "object",
  "additionalProperties": false,
  "default": {}
}`

const NamespaceSchemaTemplate = `{
  "type": "object",
  "additionalProperties": false,
  "default": {},
  "required": ["kind", "spec", "metadata"],
  "properties": {
    "kind": {"type": "string"},
    "version": {"type": "string", "default": "v1"},
    "metadata": %v,
    "spec": %v
  }
}`

// GetNamespaceSchema returns namespace schema
func GetNamespaceSchema() string {
	return fmt.Sprintf(NamespaceSchemaTemplate, MetadataSchema, NamespaceSpecSchema)
}

// MarshalNamespace marshals namespace to JSON
func MarshalNamespace(ns Namespace) ([]byte, error) {
	return utils.FastMarshal(ns)
}

// UnmarshalNamespace unmarshals role from JSON or YAML,
// sets defaults and checks the schema
func UnmarshalNamespace(data []byte) (*Namespace, error) {
	if len(data) == 0 {
		return nil, trace.BadParameter("missing namespace data")
	}

	// always skip schema validation on namespaces unmarshal
	// the namespace is always created by teleport now
	var namespace Namespace
	if err := utils.FastUnmarshal(data, &namespace); err != nil {
		return nil, trace.BadParameter(err.Error())
	}

	if err := namespace.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}

	return &namespace, nil
}

// SortedNamespaces sorts namespaces
type SortedNamespaces []Namespace

// Len returns length of a role list
func (s SortedNamespaces) Len() int {
	return len(s)
}

// Less compares roles by name
func (s SortedNamespaces) Less(i, j int) bool {
	return s[i].Metadata.Name < s[j].Metadata.Name
}

// Swap swaps two roles in a list
func (s SortedNamespaces) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func IsValidNamespace(s string) bool {
	return validNamespace.MatchString(s)
}

var validNamespace = regexp.MustCompile(`^[A-Za-z0-9]+$`)
