package content

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/iximiuz/labctl/internal/api"
	"github.com/iximiuz/labctl/internal/content"
	"github.com/iximiuz/labctl/internal/labcli"
)

type listOptions struct {
	kind content.ContentKind
}

func newListCommand(cli labcli.CLI) *cobra.Command {
	var opts listOptions

	cmd := &cobra.Command{
		Use:     "list [--kind challenge|tutorial|course]",
		Aliases: []string{"ls"},
		Short:   "List authored content, possibly filtered by kind.",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return labcli.WrapStatusError(runListContent(cmd.Context(), cli, &opts))
		},
	}

	flags := cmd.Flags()

	flags.Var(
		&opts.kind,
		"kind",
		`Content kind to filter by - one of challenge, tutorial, course (an empty string means all)`,
	)

	return cmd
}

type AuthoredContent struct {
	Challenges []api.Challenge `json:"challenges" yaml:"challenges"`
	Tutorials  []api.Tutorial  `json:"tutorials" yaml:"tutorials"`
	Courses    []api.Course    `json:"courses"    yaml:"courses"`
}

func runListContent(ctx context.Context, cli labcli.CLI, opts *listOptions) error {
	var authored AuthoredContent

	if opts.kind == "" || opts.kind == content.KindChallenge {
		challenges, err := cli.Client().ListAuthoredChallenges(ctx)
		if err != nil {
			return fmt.Errorf("cannot list authored challenges: %w", err)
		}

		authored.Challenges = challenges
	}

	if opts.kind == "" || opts.kind == content.KindTutorial {
		tutorials, err := cli.Client().ListAuthoredTutorials(ctx)
		if err != nil {
			return fmt.Errorf("cannot list authored tutorials: %w", err)
		}

		authored.Tutorials = tutorials
	}

	if opts.kind == "" || opts.kind == content.KindCourse {
		courses, err := cli.Client().ListAuthoredCourses(ctx)
		if err != nil {
			return fmt.Errorf("cannot list authored courses: %w", err)
		}

		authored.Courses = courses
	}

	if err := yaml.NewEncoder(cli.OutputStream()).Encode(authored); err != nil {
		return err
	}

	return nil
}
