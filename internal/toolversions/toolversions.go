// Package toolversions handles reading and writing tools and versions from
// asdf's .tool-versions files
package toolversions

import (
	"os"
	"slices"
	"strings"
)

// ToolVersions represents a tool along with versions specified for it
type ToolVersions struct {
	Name     string
	Versions []string
}

// FindToolVersions looks up a tool version in a tool versions file and if found
// returns a slice of versions for it.
func FindToolVersions(filepath, toolName string) (versions []string, found bool, err error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return versions, false, err
	}

	versions, found = findToolVersionsInContent(string(content), toolName)
	return versions, found, nil
}

// GetAllToolsAndVersions returns a list of all tools and associated versions
// contained in a .tool-versions file
func GetAllToolsAndVersions(filepath string) (toolVersions []ToolVersions, err error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return toolVersions, err
	}

	toolVersions = getAllToolsAndVersionsInContent(string(content))
	return toolVersions, nil
}

// Unique takes a slice of ToolVersions and returns a slice of unique tools and
// versions.
func Unique(versions []ToolVersions) (uniques []ToolVersions) {
	for _, version := range versions {
		var found bool

		for index, unique := range uniques {
			if unique.Name == version.Name {
				// Duplicate name, check versions
				for _, versionNumber := range version.Versions {
					if !slices.Contains(unique.Versions, versionNumber) {
						unique.Versions = append(unique.Versions, versionNumber)
					}
				}

				uniques[index] = unique
				found = true
				break
			}
		}

		// None with name found, add
		if !found {
			uniques = append(uniques, version)
		}
	}

	return uniques
}

// readLines reads all the lines in a given file
// removing spaces and comments which are marked by '#'
func readLines(content string) (lines []string) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.SplitN(line, "#", 2)[0]
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return
}

func findToolVersionsInContent(content, toolName string) (versions []string, found bool) {
	toolVersions := getAllToolsAndVersionsInContent(content)
	for _, tool := range toolVersions {
		if tool.Name == toolName {
			return tool.Versions, true
		}
	}

	return versions, found
}

func getAllToolsAndVersionsInContent(content string) (toolVersions []ToolVersions) {
	for _, line := range readLines(content) {
		tokens := parseLine(line)
		newTool := ToolVersions{Name: tokens[0], Versions: tokens[1:]}
		toolVersions = append(toolVersions, newTool)
	}

	return toolVersions
}

func parseLine(line string) (tokens []string) {
	for _, token := range strings.Split(line, " ") {
		token = strings.TrimSpace(token)
		if len(token) > 0 {
			tokens = append(tokens, token)
		}
	}

	return tokens
}