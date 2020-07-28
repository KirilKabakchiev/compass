package scopesdecorator

import (
	"fmt"
	"os"

	"github.com/kyma-incubator/compass/components/director/hack/plugins"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphqlOperationType string

const (
	directiveArgumentPrefix                      = "graphql"
	Query                   GraphqlOperationType = "query"
	Mutation                GraphqlOperationType = "mutation"
	directiveName                                = "hasScopes"
	directiveArg                                 = "path"
)

var _ plugin.ConfigMutator = &scopesDecoratorPlugin{}

func NewPlugin(schemaFileName string) *scopesDecoratorPlugin {
	return &scopesDecoratorPlugin{schemaFileName: schemaFileName}
}

type scopesDecoratorPlugin struct {
	schemaFileName string
}

func (p *scopesDecoratorPlugin) Name() string {
	return "scopes_decorator"
}

func (p *scopesDecoratorPlugin) MutateConfig(cfg *config.Config) error {
	fmt.Printf("[%s] Mutate Configuration\n", p.Name())

	if err := cfg.LoadSchema(); err != nil {
		return err
	}

	if cfg.Schema.Query != nil {
		for _, f := range cfg.Schema.Query.Fields {
			p.ensureDirective(f, Query)
		}
	}
	if cfg.Schema.Mutation != nil {
		for _, f := range cfg.Schema.Mutation.Fields {
			p.ensureDirective(f, Mutation)
		}
	}

	schemaFile, err := os.Create(p.schemaFileName)
	if err != nil {
		return err
	}

	f := plugins.NewFormatter(schemaFile)
	f.FormatSchema(cfg.Schema)
	return schemaFile.Close()
}

func (p *scopesDecoratorPlugin) ensureDirective(f *ast.FieldDefinition, opType GraphqlOperationType) {
	d := p.getDirective(f)
	if d == nil {
		f.Directives = append(f.Directives, &ast.Directive{
			Name:      directiveName,
			Arguments: p.getDirectiveArguments(opType, f.Name),
		})
	} else {
		d.Name = directiveName
		d.Arguments = p.getDirectiveArguments(opType, f.Name)
	}
}

func (p *scopesDecoratorPlugin) getDirective(def *ast.FieldDefinition) *ast.Directive {
	for _, d := range def.Directives {
		if d.Name == directiveName {
			return d
		}
	}
	return nil
}

func (p *scopesDecoratorPlugin) getDirectiveArguments(opType GraphqlOperationType, opName string) ast.ArgumentList {
	var args ast.ArgumentList
	path := fmt.Sprintf("%s.%s.%s", directiveArgumentPrefix, opType, opName)
	args = append(args, &ast.Argument{Name: directiveArg, Value: &ast.Value{Raw: path, Kind: ast.StringValue}})
	return args
}
