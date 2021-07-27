package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

const (
	Root = "play.ground"
)

//go:embed scripts
var static embed.FS

func main() {
	fmt.Println("starting ... ")

	if err := Run(); err != nil {
		fmt.Println(" ┬───► load failed")
		return
	}
}

func Run() error {
	overlay, err := loadFiles(static)
	if err != nil {
		return err
	}

	config := &load.Config{
		Module:  Root,
		Overlay: overlay,
	}

	insts := load.Instances([]string{"./main.cue"}, config)

	ctx := cuecontext.New()
	for _, instIdx := range insts {
		inst := ctx.BuildInstance(instIdx)

		val := inst.Value()
		fmt.Printf("%v\n", val)
	}

	return nil
}

func loadFiles(fsys fs.FS) (map[string]load.Source, error) {
	overlay := make(map[string]load.Source)
	err := fs.WalkDir(fsys, ".", func(dpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d == nil || d.Type().IsDir() || !d.Type().IsRegular() {
			return nil
		}

		if strings.Contains(dpath, ".cue") {

			f, err := fsys.Open(dpath)
			if err != nil {
				return err
			}

			buf, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			ovpath := path.Join("/"+Root, dpath)

			fmt.Println("dpath ", dpath)
			fmt.Println("ovpath ", ovpath)

			overlay[ovpath] = load.FromBytes(buf)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return overlay, nil
}
