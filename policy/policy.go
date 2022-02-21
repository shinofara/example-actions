package policy

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/google/go-github/v33/github"
	ghc "github.com/shinofara/example-actions/github"

	"github.com/open-policy-agent/opa/rego"
)

var PolicyDir = "./policy"

func New(op ...OptionInterface) *Checker {
	return &Checker{
		option: op,
	}
}

type Checker struct {
	option []OptionInterface
}

type OptionInterface interface {
	Check(ctx context.Context, repo *github.Repository, vs []Violation) ([]Violation, error)
}

type Base struct {
	name string
	ghc  ghc.Client
}

type Access struct {
	Base
}

type Option struct {
	Base
}

type Protection struct {
	Base
}

type Result struct {
	Repository string
	Violations []Violation
}

type Violation struct {
	Name    string
	Message string
}

func NewAccess(client ghc.Client) OptionInterface {
	return &Access{
		Base: Base{
			name: "access",
			ghc:  client,
		},
	}
}

func (p *Base) check(ctx context.Context, vs []Violation, data interface{}) ([]Violation, error) {
	results, err := check(ctx, p.name, data, "x = allow; z = check")
	if err != nil {
		return vs, err
	}
	msg, ok := results[0].Bindings["z"].([]interface{})
	if !ok {
		return vs, fmt.Errorf("failed to cast to bool, got %v", results)
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok {
		return vs, fmt.Errorf("failed to cast to bool, got %v", results)
	}

	if !result {
		for _, m := range msg {
			if mm, ok := m.(string); ok {
				vs = append(vs, Violation{
					Name:    p.name,
					Message: mm,
				})
			}
		}

		return vs, nil
	}

	return vs, nil
}

func check(ctx context.Context, policyName string, data interface{}, query string) (rego.ResultSet, error) {
	buf, err := ioutil.ReadFile(path.Join(PolicyDir, "/repository/", policyName+".rego"))
	if err != nil {
		return nil, err
	}

	r, err := rego.New(
		rego.Query(query),
		rego.Package("repository."+policyName),
		rego.Module(policyName+".rego", string(buf)),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	resultSet, err := r.Eval(ctx, rego.EvalInput(data))
	switch {
	case err != nil:
		return nil, err
	case len(resultSet) == 0:
		return nil, errors.New("results is 0")
	}

	return resultSet, nil
}

func (p *Access) Check(ctx context.Context, repo *github.Repository, vs []Violation) ([]Violation, error) {
	teams, err := p.ghc.GetListTeamsByRepo(ctx, *repo.Name)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"data": teams,
	}

	return p.check(ctx, vs, data)
}

func NewProtection(client ghc.Client) OptionInterface {
	return &Protection{
		Base: Base{
			name: "protection",
			ghc:  client,
		},
	}
}

func (p *Protection) Check(ctx context.Context, repo *github.Repository, vs []Violation) ([]Violation, error) {
	protection := p.ghc.GetBranchProtection(ctx, *repo.Name)
	return p.check(ctx, vs, protection)
}

func NewOption(client ghc.Client) OptionInterface {
	return &Option{
		Base: Base{
			name: "option",
			ghc:  client,
		},
	}
}

func (p *Option) Check(ctx context.Context, repo *github.Repository, vs []Violation) ([]Violation, error) {
	return p.check(ctx, vs, repo)
}

// Check policyに則っているか確認
func (re *Checker) Do(ctx context.Context, repo *github.Repository) (*Result, error) {
	result := &Result{
		Repository: *repo.Name,
	}

	ok, err := checkOptout(ctx, repo)
	switch {
	case err != nil:
		return nil, err
	case ok:
		return nil, nil
	}

	for _, p := range re.option {
		result.Violations, err = p.Check(ctx, repo, result.Violations)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func checkOptout(ctx context.Context, data interface{}) (bool, error) {
	results, err := check(ctx, "optout", data, "x = allow")
	if err != nil {
		return false, err
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok {
		return false, fmt.Errorf("failed to cast to bool, got %v", results)
	}

	return result, nil
}
