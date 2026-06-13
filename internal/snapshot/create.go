package snapshot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/kernelstub/cognitor/internal/util"
)

type CreateOptions struct {
	Name   string
	Path   string
	Source string
	Force  bool
}

type CreateResult struct {
	Path         string
	CopiedFiles  int
	CreatedFiles int
	SkippedFiles int
}

func Create(ctx context.Context, opts CreateOptions) (CreateResult, error) {
	result := CreateResult{Path: opts.Path}
	if opts.Name == "" {
		return result, fmt.Errorf("snapshot name is required")
	}
	if opts.Path == "" {
		return result, fmt.Errorf("snapshot path is required")
	}
	if err := os.MkdirAll(opts.Path, 0o755); err != nil {
		return result, err
	}
	created, err := ensureJSON(filepath.Join(opts.Path, "services.json"), []map[string]string{}, opts.Force)
	if err != nil {
		return result, err
	}
	result.CreatedFiles += created
	created, err = ensureJSON(filepath.Join(opts.Path, "registry.json"), []map[string]string{}, opts.Force)
	if err != nil {
		return result, err
	}
	result.CreatedFiles += created
	created, err = ensureReadme(opts, opts.Force)
	if err != nil {
		return result, err
	}
	result.CreatedFiles += created
	if opts.Source != "" {
		copied, skipped, err := copySnapshotInputs(ctx, opts.Source, opts.Path, opts.Force)
		if err != nil {
			return result, err
		}
		result.CopiedFiles = copied
		result.SkippedFiles = skipped
	}
	return result, nil
}

func ensureJSON(path string, value any, force bool) (int, error) {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return 0, nil
		}
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return 0, err
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return 0, err
	}
	return 1, nil
}

func ensureReadme(opts CreateOptions, force bool) (int, error) {
	path := filepath.Join(opts.Path, "SNAPSHOT.md")
	if !force {
		if _, err := os.Stat(path); err == nil {
			return 0, nil
		}
	}
	body := fmt.Sprintf("# %s\n\nPlace Windows binaries and optional sidecar analysis files in this directory, then run `cognitor scan --snapshot %s --path %s --out %s.db`.\n", opts.Name, opts.Name, opts.Path, opts.Name)
	if opts.Source != "" {
		body += fmt.Sprintf("\nCreated from source `%s`.\n", opts.Source)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return 0, err
	}
	return 1, nil
}

func copySnapshotInputs(ctx context.Context, source string, dest string, force bool) (int, int, error) {
	var files []string
	err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if util.IsBinaryLike(path) || util.IsArtifactLike(path) || isSidecar(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	sort.Strings(files)
	var copied, skipped int
	var errs []error
	for _, src := range files {
		if err := ctx.Err(); err != nil {
			return copied, skipped, err
		}
		rel, err := filepath.Rel(source, src)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		dst := filepath.Join(dest, rel)
		ok, err := copyFile(src, dst, force)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if ok {
			copied++
		} else {
			skipped++
		}
	}
	return copied, skipped, errors.Join(errs...)
}

func isSidecar(path string) bool {
	for _, suffix := range []string{".analysis.json", ".symbols.json", ".version", ".manifest", ".manifest.json"} {
		if len(path) >= len(suffix) && path[len(path)-len(suffix):] == suffix {
			return true
		}
	}
	return false
}

func copyFile(src string, dst string, force bool) (bool, error) {
	if !force {
		if _, err := os.Stat(dst); err == nil {
			return false, nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return false, err
	}
	in, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return false, err
	}
	return true, out.Close()
}
