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

func (colp *TreePlugin) printTreeToFile(proc *datastore.Process, depth int, outputFile *os.File) {
	if proc == nil {
		return
	}

	procFullInfo := proc.Args
	if len(procFullInfo) == 0 {
		procFullInfo = proc.FullPath
	}
	nodeValue := fmt.Sprintf("%-4d %s - %s", proc.PID, proc.ImageName, procFullInfo)
	fmt.Fprintf(outputFile, "%s%s\n", strings.Repeat(" ", depth*4), nodeValue)

	var subProcGroup datastore.ProcessByPID
	for _, child := range datastore.PIDToProcess {
		if child.ParentProc != nil && child.ParentProc.PID == proc.PID {
			subProcGroup = append(subProcGroup, child)
		}
	}

	sort.Sort(subProcGroup)
	for _, subProc := range subProcGroup {
		colp.printTreeToFile(subProc, depth+1, outputFile)
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
