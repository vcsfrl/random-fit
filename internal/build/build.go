package build

//
//import (
//	"bytes"
//	"fmt"
//	"github.com/google/uuid"
//	"github.com/vcsfrl/random-fit/internal/combination"
//	"github.com/vcsfrl/random-fit/internal/tmp/platform/random"
//	"go.starlark.net/starlark"
//	"go.starlark.net/syntax"
//	"html/template"
//	"path/filepath"
//	"time"
//)
//
//var ErrBuilding = fmt.Errorf("error building combination")
//var ErrBuildingScript = fmt.Errorf("%w: error starlark script", ErrBuilding)
//
//type Builder struct {
//	thread      *starlark.Thread
//	builderFunc starlark.Value
//	definition  *combination.CombinationDefinition
//}
//
//func NewBuilder(definition *combination.CombinationDefinition) (*Builder, error) {
//	builder := &Builder{definition: definition}
//	err := builder.start()
//	if err != nil {
//		return nil, err
//	}
//
//	return builder, nil
//
//}
//
//func (bd *Builder) Build() (*combination.Combination, error) {
//	// Run the Starlark script from the definition to create a new combination.
//	combinationData, err := starlark.Call(bd.thread, bd.builderFunc, nil, nil)
//	if err != nil {
//		return nil, ErrBuildingScript
//	}
//
//	// Build the template from the definition.
//	base := filepath.Base(bd.definition.GoTemplate)
//	templateData := template.Must(template.New(base).ParseFiles(bd.definition.GoTemplate))
//
//	output := new(bytes.Buffer)
//	if err := templateData.Execute(output, combinationData); err != nil {
//		return nil, err
//	}
//
//	// Build the combination from the template and the data from the Starlark script.
//	combination := &combination.Combination{
//		UUID:       uuid.New(),
//		Definition: bd.definition,
//		Data:       combinationData,
//		Output:     output,
//	}
//
//	return combination, nil
//}
