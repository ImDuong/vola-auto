package processes

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ImDuong/vola-auto/config"
	"github.com/ImDuong/vola-auto/datastore"
)

type (
	TreePlugin struct {
	}
)

const (
	ProcessCollectionFolderName = "processes"
)

func (colp *TreePlugin) GetName() string {
	return "PROCESS TREE COLLECTION PLUGIN"
}

func (colp *TreePlugin) GetArtifactsCollectionPath() string {
	return filepath.Join(colp.GetArtifactsCollectionFolderpath(), "tree.txt")
}

func (colp *TreePlugin) GetArtifactsCollectionFolderpath() string {
	return filepath.Join(config.Default.OutputFolder, ProcessCollectionFolderName)
}

func (colp *TreePlugin) printTreeToFile(process *datastore.Process, depth int, outputFile *os.File) {
	if process == nil {
		return
	}

	nodeValue := fmt.Sprintf("%d %s - %s", process.PID, process.ImageName, process.Args)
	fmt.Fprintf(outputFile, "%s%s\n", strings.Repeat(" ", depth*4), nodeValue)

	for _, child := range datastore.PIDToProcess {
		if child.ParentProc != nil && child.ParentProc.PID == process.PID {
			colp.printTreeToFile(child, depth+1, outputFile)
		}
	}
}

// print out process tree
func (colp *TreePlugin) Run() error {
	err := os.MkdirAll(colp.GetArtifactsCollectionFolderpath(), 0755)
	if err != nil {
		return err
	}

	treeFileWriter, err := os.OpenFile(colp.GetArtifactsCollectionPath(), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer treeFileWriter.Close()

	var roots datastore.ProcessByPID
	for _, process := range datastore.PIDToProcess {
		if process.ParentProc == nil {
			roots = append(roots, process)
		}
	}

	sort.Sort(roots)

	for _, root := range roots {
		colp.printTreeToFile(root, 0, treeFileWriter)
	}

	return nil
}
