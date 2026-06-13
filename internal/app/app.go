package app

import (
	"context"
	"io"
	"os"

	"github.com/kernelstub/cognitor/internal/cli"
)

func Run(ctx context.Context, args []string) error {
	root := cli.NewRoot(os.Stdout, os.Stderr)
	root.SetArgs(args)
	return root.ExecuteContext(ctx)
}

func RunWithIO(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) error {
	root := cli.NewRoot(stdout, stderr)
	root.SetArgs(args)
	return root.ExecuteContext(ctx)
}
