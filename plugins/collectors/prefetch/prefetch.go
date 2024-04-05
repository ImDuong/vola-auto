package prefetch

import (
	"os"
	"path/filepath"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/plugins/collectors"
	"github.com/alitto/pond"
)

type (
	PrefetchPlugin struct {
		WorkerPool *pond.WorkerPool
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

	filePlg := collectors.FilesPlugin{
		WorkerPool: colp.WorkerPool,
	}
	prefetchFiles, err := filePlg.FindFilesByRegex(`\.pf`)
	if err != nil {
		return err
	}

	err = filePlg.DumpFiles(prefetchFiles, colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	err = filePlg.RenameDumpedFilesExtention(".pf.dat", "", colp.GetArtifactsCollectionPath())
	if err != nil {
		return err
	}

	return nil
}
