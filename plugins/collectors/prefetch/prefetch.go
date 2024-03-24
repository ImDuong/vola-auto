package prefetch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	PrefetchPlugin struct {
	}
)

func (colp *PrefetchPlugin) GetName() string {
	return "PREFETCH COLLECTION PLUGIN"
}

func (colp *PrefetchPlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(config.Default.OutputFolder, "prefetch")
}

// dump prefetch files
func (colp *PrefetchPlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionPath(), 0755)
	if err != nil {
		return err
	}

	filePlg := collectors.FilesPlugin{}
	prefetchFiles, err := filePlg.FindFilesByRegex("\\.pf$")
	if err != nil {
		return err
	}

	dumpFilesPool := pond.New(20, 100)
	var aggregatedError error
	var aggregateErrorMutex sync.Mutex
	for i := range prefetchFiles {
		copiedIdx := i
		dumpFilesPool.Submit(func() {
			err := filePlg.DumpFile(prefetchFiles[copiedIdx], colp.GetArtifactsCollectionPath())
			if err != nil {
				aggregateErrorMutex.Lock()
				aggregatedError = fmt.Errorf("%w;%w", aggregatedError, err)
				aggregateErrorMutex.Unlock()
			}
		})
	}
	dumpFilesPool.StopAndWait()
	if aggregatedError != nil {
		return aggregatedError
	}

	return filepath.Walk(colp.GetArtifactsCollectionPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if it's a regular file and ends with ".pf.dat"
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".pf.dat") {
			newName := strings.TrimSuffix(info.Name(), ".dat")
			err := os.Rename(path, filepath.Join(filepath.Dir(path), newName))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
