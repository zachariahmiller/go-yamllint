package runner

import (
	"context"
	"crypto/rand"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"

	"github.com/zachariahmiller/go-yamllint/internal/pysite"
)

func Run(name string, cmdArgs []string, wasm []byte, stdin io.Reader, stdout io.Writer, stderr io.Writer, cwd string) int {
	ctx := context.Background()

	rtCfg := wazero.NewRuntimeConfig()
	uc, err := os.UserCacheDir()
	if err == nil {
		cache, err := wazero.NewCompilationCacheWithDir(filepath.Join(uc, "com.github.wasilibs"))
		if err == nil {
			rtCfg = rtCfg.WithCompilationCache(cache)
		}
	}
	rt := wazero.NewRuntimeWithConfig(ctx, rtCfg)

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	args := []string{"python", ".venv/bin/yamllint"}
	args = append(args, cmdArgs...)

	libDir, _ := fs.Sub(pysite.Site, "lib")
	venvDir, _ := fs.Sub(pysite.Site, ".venv")

	cfg := wazero.NewModuleConfig().
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithStderr(stderr).
		WithStdout(stdout).
		WithStdin(stdin).
		WithRandSource(rand.Reader).
		WithArgs(args...).
		WithFSConfig(wazero.NewFSConfig().
			WithFSMount(libDir, "lib").
			WithFSMount(venvDir, ".venv").
			WithDirMount(cwd, "/")).
		WithEnv("PYTHONPATH", ".venv/lib/python3.12/site-packages").
		WithEnv("PYTHONDONTWRITEBYTECODE", "1")
	for _, env := range os.Environ() {
		k, v, _ := strings.Cut(env, "=")
		cfg = cfg.WithEnv(k, v)
	}

	_, err = rt.InstantiateWithConfig(ctx, wasm, cfg)
	if err != nil {
		if sErr, ok := err.(*sys.ExitError); ok {
			return int(sErr.ExitCode())
		}
		log.Fatal(err)
	}
	return 0
}
