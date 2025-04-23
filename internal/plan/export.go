package plan

type Exporter struct {
	OutputDir string
}

func NewExporter(outputDir string) *Exporter {
	return &Exporter{
		OutputDir: outputDir,
	}
}

func (e *Exporter) Export(plan *Plan) error {
	// Implement the export logic here
	// This could involve writing the plan to a file in the specified output directory
	// For example, you might want to write it as a JSON or YAML file

	// Placeholder for actual implementation
	return nil
}

//func toSnakeCase(s string) string {
//	s = strings.ToLower(strings.TrimSpace(s))
//	var result []rune
//	for i, r := range s {
//		if i > 0 && r >= 'A' && r <= 'Z' {
//			result = append(result, '_')
//		}
//		result = append(result, r)
//	}
//
//	return s
//}
