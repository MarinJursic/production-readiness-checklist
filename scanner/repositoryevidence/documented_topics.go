package repositoryevidence

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	ArchitectureTopicsCollectorID = "prc.collect.prc.36.002.c1@0.1"
	ConventionsTopicsCollectorID  = "prc.collect.prc.36.005.c1@0.1"

	architectureTopicsControlID = "PRC-36-002"
	conventionsTopicsControlID  = "PRC-36-005"
	architectureTopicsFactKey   = "prc_36_002_c1.documented_topic_content"
	conventionsTopicsFactKey    = "prc_36_005_c1.documented_convention_content"
	architectureTopicsParameter = "prc_36_002_c1.required_documented_topic_content_keys"
	conventionsTopicsParameter  = "prc_36_005_c1.required_documented_convention_content_keys"
)

var (
	architectureTopics = []string{"architecture", "data-flow", "dependencies", "deployment", "recovery"}
	conventionTopics   = []string{"coding", "release", "review", "testing"}
)

// DocumentedTopicsProvider proves only a narrow positive case: every required
// topic has an exact Markdown heading followed by meaningful, bounded content
// in the sealed repository inventory. Unsupported prose, aliases, malformed
// Markdown, and placeholders remain incomplete evidence rather than failures.
type DocumentedTopicsProvider struct {
	inventory   model.Inventory
	id          string
	controlID   string
	factKey     string
	parameterID string
	topics      []string
}

func NewDocumentedTopicsProviders(item model.Inventory) ([]*DocumentedTopicsProvider, error) {
	if item.Root == "" || item.Digest == "" {
		return nil, fmt.Errorf("documented-topic collectors require a sealed inventory")
	}
	return []*DocumentedTopicsProvider{
		{item, ArchitectureTopicsCollectorID, architectureTopicsControlID, architectureTopicsFactKey, architectureTopicsParameter, append([]string(nil), architectureTopics...)},
		{item, ConventionsTopicsCollectorID, conventionsTopicsControlID, conventionsTopicsFactKey, conventionsTopicsParameter, append([]string(nil), conventionTopics...)},
	}, nil
}

func (provider *DocumentedTopicsProvider) ID() string { return provider.id }

func (provider *DocumentedTopicsProvider) Authority() controlprogram.Authority {
	return controlprogram.AuthorityRepository
}

func (provider *DocumentedTopicsProvider) Collect(ctx context.Context, request controlruntime.Request) (controlprogram.Evidence, error) {
	if provider == nil || request.Template.ControlID != provider.controlID ||
		request.Template.CollectorContract.CollectorID != provider.id {
		return controlprogram.Evidence{}, fmt.Errorf("documented-topic collector received an unsupported template")
	}
	if err := ctx.Err(); err != nil {
		return controlprogram.Evidence{}, err
	}
	values := make(map[string]string, len(provider.topics))
	for _, topic := range provider.topics {
		values[topic] = ""
	}
	readFiles, readBytes := 0, int64(0)
	for _, record := range documentationRecords(provider.inventory) {
		if topicValuesComplete(values) {
			break
		}
		if err := ctx.Err(); err != nil {
			return controlprogram.Evidence{}, err
		}
		if record.Size > maximumDocumentationFileBytes || readFiles >= maximumDocumentationFiles ||
			readBytes+record.Size > maximumDocumentationBytes {
			continue
		}
		data, err := workspaceinventory.ReadVerifiedFile(provider.inventory, record.Path, maximumDocumentationFileBytes)
		if err != nil {
			return controlprogram.Evidence{}, fmt.Errorf("read inventoried documentation %s: %w", record.Path, err)
		}
		readFiles++
		readBytes += record.Size
		if !validDocumentationText(data) {
			continue
		}
		sections := exactMarkdownTopicSections(data, provider.topics)
		for topic, content := range sections {
			if values[topic] == "" {
				values[topic] = boundedTopicValue(record.Path, content)
			}
		}
	}
	complete := topicValuesComplete(values)
	facts := map[string]controlprogram.Fact{
		provider.factKey: {Type: controlprogram.FactStringMap, Complete: complete, Values: values},
	}
	return controlruntime.NewApplicableEvidence(
		request,
		"repository-documented-topics-"+controlprogram.ProgramSHA256(request.Program)[:16],
		facts,
		complete,
	)
}

func documentedTopicsBinding(item model.Inventory, template controlprogramcatalog.Template) (controlprogramcatalog.RuntimeBinding, bool) {
	if item.Digest == "" || template.ClauseOrdinal != 1 {
		return controlprogramcatalog.RuntimeBinding{}, false
	}
	var parameter string
	var topics []string
	switch template.CollectorContract.CollectorID {
	case ArchitectureTopicsCollectorID:
		if template.ControlID != architectureTopicsControlID {
			return controlprogramcatalog.RuntimeBinding{}, false
		}
		parameter, topics = architectureTopicsParameter, architectureTopics
	case ConventionsTopicsCollectorID:
		if template.ControlID != conventionsTopicsControlID {
			return controlprogramcatalog.RuntimeBinding{}, false
		}
		parameter, topics = conventionsTopicsParameter, conventionTopics
	default:
		return controlprogramcatalog.RuntimeBinding{}, false
	}
	subject := "repository@sha256:" + item.Digest
	return controlprogramcatalog.RuntimeBinding{
		SubjectID: subject, Subjects: []string{subject}, InventorySHA256: item.Digest,
		AllowNotApplicable: false, ApplicabilityProofContractSHA256: applicabilityContractSHA256,
		MaximumEvidenceAgeSeconds: 300,
		ScannerInventoryParameters: map[string]controlprogram.Parameter{
			parameter: {Type: controlprogram.FactStringSet, Strings: append([]string(nil), topics...)},
		},
	}, true
}

func documentationRecords(item model.Inventory) []model.FileRecord {
	records := make([]model.FileRecord, 0)
	for _, file := range item.Files {
		extension := strings.ToLower(filepath.Ext(file.Path))
		if extension != ".md" && extension != ".markdown" && extension != ".mdx" {
			continue
		}
		records = append(records, file)
	}
	// Root documentation and shallower paths are the most explicit public
	// entry points. Stable depth-first priority lets a positive existence proof
	// finish early without reading an unrelated large documentation corpus.
	sort.Slice(records, func(left, right int) bool {
		leftDepth := strings.Count(records[left].Path, "/")
		rightDepth := strings.Count(records[right].Path, "/")
		if leftDepth != rightDepth {
			return leftDepth < rightDepth
		}
		return records[left].Path < records[right].Path
	})
	return records
}

func topicValuesComplete(values map[string]string) bool {
	if len(values) == 0 {
		return false
	}
	for _, value := range values {
		if value == "" {
			return false
		}
	}
	return true
}

type markdownHeading struct {
	line  int
	level int
	topic string
}

func exactMarkdownTopicSections(data []byte, topics []string) map[string]string {
	allowed := make(map[string]string, len(topics))
	for _, topic := range topics {
		allowed[strings.ReplaceAll(topic, "-", " ")] = topic
	}
	lines := bytes.Split(data, []byte{'\n'})
	headings := make([]markdownHeading, 0)
	fence := byte(0)
	for index, line := range lines {
		trimmed := bytes.TrimLeft(line, " \t")
		if len(line)-len(trimmed) <= 3 && (bytes.HasPrefix(trimmed, []byte("```")) || bytes.HasPrefix(trimmed, []byte("~~~"))) {
			marker := trimmed[0]
			if fence == 0 {
				fence = marker
			} else if fence == marker {
				fence = 0
			}
			continue
		}
		if fence != 0 {
			continue
		}
		level, title, ok := atxHeading(line)
		if !ok {
			continue
		}
		if topic, exists := allowed[normalizeHeading(title)]; exists {
			headings = append(headings, markdownHeading{line: index, level: level, topic: topic})
		} else {
			headings = append(headings, markdownHeading{line: index, level: level})
		}
	}
	result := map[string]string{}
	for index, heading := range headings {
		if heading.topic == "" || result[heading.topic] != "" {
			continue
		}
		end := len(lines)
		for _, next := range headings[index+1:] {
			if next.level <= heading.level {
				end = next.line
				break
			}
		}
		content := visibleMarkdownContent(bytes.Join(lines[heading.line+1:end], []byte{'\n'}))
		if meaningfulTopicContent(content) {
			result[heading.topic] = string(bytes.TrimSpace(content))
		}
	}
	return result
}

func atxHeading(line []byte) (int, string, bool) {
	trimmed := bytes.TrimLeft(line, " \t")
	if len(line)-len(trimmed) > 3 {
		return 0, "", false
	}
	level := 0
	for level < len(trimmed) && level < 6 && trimmed[level] == '#' {
		level++
	}
	if level == 0 || level == len(trimmed) || trimmed[level] != ' ' && trimmed[level] != '\t' {
		return 0, "", false
	}
	title := strings.TrimSpace(string(trimmed[level:]))
	title = strings.TrimSpace(strings.TrimRight(title, "#"))
	return level, title, title != ""
}

func normalizeHeading(value string) string {
	words := strings.FieldsFunc(strings.ToLower(value), func(character rune) bool {
		return !unicode.IsLetter(character) && !unicode.IsDigit(character)
	})
	return strings.Join(words, " ")
}

func meaningfulTopicContent(data []byte) bool {
	words := strings.FieldsFunc(strings.ToLower(string(data)), func(character rune) bool {
		return !unicode.IsLetter(character) && !unicode.IsDigit(character)
	})
	if len(words) < 3 {
		return false
	}
	plain := strings.Join(words, " ")
	for _, prefix := range []string{"tbd", "todo", "wip", "placeholder", "coming soon", "not documented", "not yet documented", "to be determined", "to be documented"} {
		if plain == prefix || strings.HasPrefix(plain, prefix+" ") {
			return false
		}
	}
	switch plain {
	case "to be determined", "to be documented", "coming soon", "not applicable", "not yet documented":
		return false
	}
	return len([]rune(plain)) >= 20
}

func validDocumentationText(data []byte) bool {
	if !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0 {
		return false
	}
	for _, value := range string(data) {
		if value < 0x20 && value != '\n' && value != '\r' && value != '\t' || value == 0x7f {
			return false
		}
	}
	return true
}

func visibleMarkdownContent(data []byte) []byte {
	withoutComments := make([]byte, 0, len(data))
	for len(data) > 0 {
		start := bytes.Index(data, []byte("<!--"))
		if start < 0 {
			withoutComments = append(withoutComments, data...)
			break
		}
		withoutComments = append(withoutComments, data[:start]...)
		end := bytes.Index(data[start+4:], []byte("-->"))
		if end < 0 {
			break
		}
		data = data[start+4+end+3:]
	}
	lines := bytes.Split(withoutComments, []byte{'\n'})
	visible := make([][]byte, 0, len(lines))
	fence := byte(0)
	for _, line := range lines {
		trimmed := bytes.TrimLeft(line, " \t")
		if len(line)-len(trimmed) <= 3 && (bytes.HasPrefix(trimmed, []byte("```")) || bytes.HasPrefix(trimmed, []byte("~~~"))) {
			marker := trimmed[0]
			if fence == 0 {
				fence = marker
			} else if fence == marker {
				fence = 0
			}
			continue
		}
		if fence != 0 {
			continue
		}
		if _, _, heading := atxHeading(line); heading {
			continue
		}
		visible = append(visible, line)
	}
	return bytes.Join(visible, []byte{'\n'})
}

func boundedTopicValue(path, content string) string {
	value := path + ": " + strings.Join(strings.Fields(content), " ")
	if len(value) <= controlprogram.MaxStringBytes {
		return value
	}
	value = value[:controlprogram.MaxStringBytes]
	for !utf8.ValidString(value) {
		value = value[:len(value)-1]
	}
	return value
}
