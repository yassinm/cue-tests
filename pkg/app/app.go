package app

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"play.ground/pkg/scripts"
)

const (
	Root = "play.ground"
)

type (
	Resolver interface {
		Resolve(path string) ([]byte, error)
	}
)

func Run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	overlay, err := loadFiles(scripts.StaticFs, cwd)
	if err != nil {
		return err
	}

	config := &load.Config{
		Module:  Root,
		Overlay: overlay,
	}

	insts := load.Instances([]string{"./local/main.cue"}, config)

	ctx := cuecontext.New()
	for _, instIdx := range insts {
		inst := ctx.BuildInstance(instIdx)

		val := inst.Value()
		fmt.Printf("%v\n", val)
	}

	return nil
}

func loadFiles(fsys fs.FS, rDir string) (map[string]load.Source, error) {
	overlay := make(map[string]load.Source)

	abs := func(dpath string) string {
		return path.Join(rDir, dpath)
	}

	err := fs.WalkDir(fsys, ".", func(dpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d == nil || d.Type().IsDir() || !d.Type().IsRegular() {
			return nil
		}

		ovpath := abs(path.Join("scripts", dpath))
		if !strings.HasSuffix(dpath, ".cue") {
			return nil
		}

		f, err := fsys.Open(dpath)
		if err != nil {
			return err
		}

		buf, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		fmt.Println("dpath ", dpath)
		fmt.Println("ovpath ", ovpath)
		fmt.Println("")

		overlay[ovpath] = load.FromBytes(buf)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return overlay, nil
}
