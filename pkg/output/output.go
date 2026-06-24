package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// Destination resolves the writer to use for command output.
func Destination(path string, force bool) (*os.File, bool, error) {
	if strings.TrimSpace(path) == "" {
		return os.Stdout, false, nil
	}

	if _, err := os.Stat(path); err == nil {
		if !force {
			return nil, false, fmt.Errorf("file %s already exists; use --force to overwrite", path)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, false, fmt.Errorf("stat %s: %w", path, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, false, fmt.Errorf("create file: %w", err)
	}

	return file, true, nil
}

func WriteJSON(writer io.Writer, payload interface{}) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func WriteYAML(writer io.Writer, payload interface{}) error {
	encoder := yaml.NewEncoder(writer)
	defer encoder.Close()
	if err := encoder.Encode(payload); err != nil {
		return fmt.Errorf("encode yaml: %w", err)
	}
	return nil
}

func WriteTOML(writer io.Writer, payload interface{}) error {
	encoder := toml.NewEncoder(writer)
	if err := encoder.Encode(payload); err != nil {
		return fmt.Errorf("encode toml: %w", err)
	}
	return nil
}

// WritePayload writes a payload in one of the structured formats (json, yaml, toml).
func WritePayload(writer io.Writer, payload interface{}, format string) error {
	switch format {
	case "json":
		return WriteJSON(writer, payload)
	case "yaml":
		return WriteYAML(writer, payload)
	case "toml":
		return WriteTOML(writer, payload)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteApplyResults(writer io.Writer, results []admin.ApplyResult, format string) error {
	switch format {
	case "json", "yaml", "toml":
		return WritePayload(writer, results, format)
	case "table", "":
		return writeApplyResultsTable(writer, results)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteComparisonReport(writer io.Writer, report manifest.ComparisonReport, format string) error {
	switch format {
	case "json", "yaml", "toml":
		return WritePayload(writer, report, format)
	case "table", "":
		return writeComparisonReportTable(writer, report)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func SanitizeResources(resources []manifest.Resource, excludeFields []string) []manifest.Resource {
	exclude := excludeFieldSet(excludeFields)
	sanitized := make([]manifest.Resource, 0, len(resources))

	for _, resource := range resources {
		clone := manifest.Resource{
			Type:       resource.Type,
			Realm:      resource.Realm,
			ParentType: resource.ParentType,
			Delete:     resource.Delete,
		}

		if len(resource.Data) > 0 {
			filtered := filterValue(resource.Data, exclude)
			if mapped, ok := filtered.(map[string]interface{}); ok {
				clone.Data = mapped
			}
		}

		sanitized = append(sanitized, clone)
	}

	return sanitized
}

func WriteResourceTable(writer io.Writer, resources []manifest.Resource, detailed bool) error {
	if detailed {
		rows := make([][]string, len(resources))
		for idx, resource := range resources {
			rows[idx] = []string{tableDisplayName(resource), resource.Type, objectDetails(resource)}
		}
		return writeTabTable(writer, "NAME\tTYPE\tDETAILS", rows)
	}

	rows := make([][]string, len(resources))
	for idx, resource := range resources {
		rows[idx] = []string{tableDisplayName(resource), resource.Type}
	}
	return writeTabTable(writer, "NAME\tTYPE", rows)
}

func WriteRelationshipTable(writer io.Writer, relationships []manifest.RelationshipOperation) error {
	rows := make([][]string, len(relationships))
	for idx, relationship := range relationships {
		rows[idx] = []string{relationship.Kind, relationship.Method, relationship.Path}
	}
	return writeTabTable(writer, "KIND\tMETHOD\tPATH", rows)
}

func writeApplyResultsTable(writer io.Writer, results []admin.ApplyResult) error {
	rows := make([][]string, 0, len(results))
	for _, result := range results {
		rows = append(rows, []string{result.Resource, result.Realm, result.Name, result.Action, strconv.Itoa(result.Status)})
		if result.Error != "" {
			rows = append(rows, []string{"", "", "", "ERROR: " + result.Error, ""})
		}
	}
	return writeTabTable(writer, "RESOURCE\tREALM\tNAME\tACTION\tSTATUS", rows)
}

func writeTabTable(writer io.Writer, header string, rows [][]string) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	defer table.Flush()

	if _, err := fmt.Fprintln(table, header); err != nil {
		return fmt.Errorf("write table header: %w", err)
	}

	format := strings.Repeat("%s\t", len(strings.Split(header, "\t"))) + "\n"
	for _, row := range rows {
		if _, err := fmt.Fprintf(table, format, stringSliceToAny(row)...); err != nil {
			return fmt.Errorf("write table row: %w", err)
		}
	}
	return nil
}

func stringSliceToAny(values []string) []any {
	result := make([]any, len(values))
	for idx := range values {
		result[idx] = values[idx]
	}
	return result
}

func writeComparisonReportTable(writer io.Writer, report manifest.ComparisonReport) error {
	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	defer table.Flush()

	status := "match"
	if !report.Match {
		status = "different"
	}
	if _, err := fmt.Fprintln(table, "STATUS\tMISSING_RESOURCES\tMISMATCHED_RESOURCES\tUNEXPECTED_RESOURCES\tMISSING_RELATIONSHIPS\tUNEXPECTED_RELATIONSHIPS"); err != nil {
		return fmt.Errorf("write comparison header: %w", err)
	}
	if _, err := fmt.Fprintf(table, "%s\t%d\t%d\t%d\t%d\t%d\n", status, len(report.MissingResources), len(report.MismatchedResources), len(report.UnexpectedResources), len(report.MissingRelationships), len(report.UnexpectedRelationships)); err != nil {
		return fmt.Errorf("write comparison row: %w", err)
	}

	if err := writeReportSection(writer, "Missing Resources", report.MissingResources, func(w io.Writer) error { return WriteResourceTable(w, report.MissingResources, true) }); err != nil {
		return err
	}
	if err := writeReportSection(writer, "Mismatched Resources", report.MismatchedResources, func(w io.Writer) error { return writeMismatchedResourceTable(w, report.MismatchedResources) }); err != nil {
		return err
	}
	if err := writeReportSection(writer, "Unexpected Resources", report.UnexpectedResources, func(w io.Writer) error { return WriteResourceTable(w, report.UnexpectedResources, true) }); err != nil {
		return err
	}
	if err := writeReportSection(writer, "Missing Relationships", report.MissingRelationships, func(w io.Writer) error { return WriteRelationshipTable(w, report.MissingRelationships) }); err != nil {
		return err
	}
	if err := writeReportSection(writer, "Unexpected Relationships", report.UnexpectedRelationships, func(w io.Writer) error { return WriteRelationshipTable(w, report.UnexpectedRelationships) }); err != nil {
		return err
	}

	return nil
}

func writeReportSection[T any](writer io.Writer, title string, items []T, write func(io.Writer) error) error {
	if len(items) == 0 {
		return nil
	}
	if _, err := fmt.Fprintf(writer, "\n%s\n", title); err != nil {
		return err
	}
	return write(writer)
}

func writeMismatchedResourceTable(writer io.Writer, mismatches []manifest.ResourceMismatch) error {
	mismatchTable := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	defer mismatchTable.Flush()
	if _, err := fmt.Fprintln(mismatchTable, "RESOURCE\tEXPECTED\tACTUAL"); err != nil {
		return fmt.Errorf("write mismatched resource header: %w", err)
	}
	for _, mismatch := range mismatches {
		if _, err := fmt.Fprintf(mismatchTable, "%s\t%s\t%s\n", tableDisplayName(mismatch.Expected), objectDetails(mismatch.Expected), objectDetails(mismatch.Actual)); err != nil {
			return fmt.Errorf("write mismatched resource row: %w", err)
		}
	}
	return nil
}

func excludeFieldSet(excludeFields []string) map[string]struct{} {
	exclude := make(map[string]struct{}, len(excludeFields))
	for _, field := range excludeFields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		exclude[field] = struct{}{}
	}
	return exclude
}

func filterValue(raw interface{}, exclude map[string]struct{}) interface{} {
	switch typed := raw.(type) {
	case map[string]interface{}:
		return filterMap(typed, exclude)
	case []interface{}:
		filtered := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			if itemMap, ok := item.(map[string]interface{}); ok {
				filtered = append(filtered, filterMap(itemMap, exclude))
				continue
			}
			filtered = append(filtered, item)
		}
		return filtered
	default:
		return raw
	}
}

func filterMap(values map[string]interface{}, exclude map[string]struct{}) map[string]interface{} {
	result := make(map[string]interface{}, len(values))
	for key, value := range values {
		if _, skip := exclude[key]; skip {
			continue
		}
		result[key] = value
	}
	return result
}

func objectDetails(resource manifest.Resource) string {
	if len(resource.Data) == 0 {
		return ""
	}

	encoded, err := json.Marshal(resource.Data)
	if err != nil {
		return ""
	}
	return string(encoded)
}

func tableDisplayName(resource manifest.Resource) string {
	name := resource.DisplayName()
	if resource.Type == "realm" {
		return name
	}

	realm := strings.TrimSpace(resource.Realm)
	if realm == "" {
		return name
	}
	if name == "" {
		return realm
	}
	return fmt.Sprintf("%s/%s", realm, name)
}

// WriteResourcesToDir writes one file per resource into dir, plus an optional
// relationships file. Each resource is wrapped as a single-resource manifest so
// it round-trips through apply.
func WriteResourcesToDir(dir string, resources []manifest.Resource,
	relationships []manifest.RelationshipOperation, format string, force bool) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}
	ext, err := formatExt(format)
	if err != nil {
		return err
	}
	used := make(map[string]struct{})
	for i, r := range resources {
		name := safeFileName(r.Realm, r.Type, r.Identifier())
		if name == "" {
			name = safeFileName(r.Realm, r.Type, r.DisplayName())
		}
		if name == "" {
			name = fmt.Sprintf("resource-%d", i)
		}
		unique := uniqueFileName(name, used)
		used[unique] = struct{}{}
		path := filepath.Join(dir, unique+"."+ext)
		payload := map[string]interface{}{"resources": []manifest.Resource{r}}
		if err := writeStructuredFile(path, payload, format, force); err != nil {
			return err
		}
	}
	if len(relationships) > 0 {
		path := filepath.Join(dir, "relationships."+ext)
		payload := map[string]interface{}{"relationships": relationships}
		if err := writeStructuredFile(path, payload, format, force); err != nil {
			return err
		}
	}
	return nil
}
func writeStructuredFile(path string, payload interface{}, format string, force bool) error {
	file, shouldClose, err := Destination(path, force)
	if err != nil {
		return err
	}
	if shouldClose {
		defer file.Close()
	}
	return WritePayload(file, payload, format)
}
func formatExt(format string) (string, error) {
	switch format {
	case "json":
		return "json", nil
	case "yaml":
		return "yaml", nil
	case "toml":
		return "toml", nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

var unsafeFileName = regexp.MustCompile(`[^a-z0-9._-]+`)

func safeFileName(realm, resourceType, identifier string) string {
	parts := make([]string, 0, 3)
	for _, part := range []string{realm, resourceType, identifier} {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	joined := strings.Join(parts, "__")
	joined = strings.ToLower(joined)
	joined = unsafeFileName.ReplaceAllString(joined, "-")
	joined = strings.Trim(joined, "-")
	// collapse repeated dashes
	for strings.Contains(joined, "--") {
		joined = strings.ReplaceAll(joined, "--", "-")
	}
	return joined
}
func uniqueFileName(base string, used map[string]struct{}) string {
	if _, exists := used[base]; !exists {
		return base
	}
	// deterministic: sort existing keys so collision suffix is stable
	keys := make([]string, 0, len(used))
	for k := range used {
		if strings.HasPrefix(k, base+"-") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	n := 1
	for {
		candidate := base + "-" + strconv.Itoa(n)
		if _, exists := used[candidate]; !exists {
			return candidate
		}
		n++
	}
}
