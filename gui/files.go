package gui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/salaleser/vk-poster/logger"

	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/file_dlg/roots"
)

var (
	fileColor      = gxui.Color{R: 0.7, G: 0.8, B: 1.0, A: 1}
	directoryColor = gxui.Color{R: 0.8, G: 1.0, B: 0.7, A: 1}
)

// filesAt returns a list of all immediate files in the given directory path.
func filesAt(path string) []string {
	files := []string{}
	filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err == nil && path != subpath {
			files = append(files, subpath)
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return files
}

// filesAdapter is an implementation of the gxui.ListAdapter interface.
// The AdapterItems returned by this adapter are absolute file path strings.
type filesAdapter struct {
	gxui.AdapterBase
	files []string // The absolute file paths
}

// SetFiles assigns the specified list of absolute-path files to this adapter.
func (a *filesAdapter) SetFiles(files []string) {
	a.files = files
	a.DataChanged(false)
}

func (a *filesAdapter) Count() int {
	return len(a.files)
}

func (a *filesAdapter) ItemAt(index int) gxui.AdapterItem {
	return a.files[index]
}

func (a *filesAdapter) ItemIndex(item gxui.AdapterItem) int {
	path := item.(string)
	for i, f := range a.files {
		if f == path {
			return i
		}
	}
	return -1 // Not found
}

func (a *filesAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	path := a.files[index]
	_, name := filepath.Split(path)
	label := theme.CreateLabel()
	label.SetText(name)
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		label.SetColor(directoryColor)
	} else {
		label.SetColor(fileColor)
	}
	return label
}

func (a *filesAdapter) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}

// directory implements the gxui.TreeNode interface to represent a directory
// node in a file-system.
type directory struct {
	path    string   // The absolute path of this directory.
	subdirs []string // The absolute paths of all immediate sub-directories.
}

// directoryAt returns a directory structure populated with the immediate
// subdirectories at the given path.
func directoryAt(path string) directory {
	directory := directory{path: path}
	filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
		if err == nil && path != subpath {
			if info.IsDir() {
				directory.subdirs = append(directory.subdirs, subpath)
				return filepath.SkipDir
			}
		}
		return nil
	})
	return directory
}

// Count implements gxui.TreeNodeContainer.
func (d directory) Count() int {
	return len(d.subdirs)
}

// NodeAt implements gxui.TreeNodeContainer.
func (d directory) NodeAt(index int) gxui.TreeNode {
	return directoryAt(d.subdirs[index])
}

// ItemIndex implements gxui.TreeNodeContainer.
func (d directory) ItemIndex(item gxui.AdapterItem) int {
	path := item.(string)
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}
	for i, subpath := range d.subdirs {
		subpath += string(filepath.Separator)
		if strings.HasPrefix(path, subpath) {
			return i
		}
	}
	return -1
}

// Item implements gxui.TreeNode.
func (d directory) Item() gxui.AdapterItem {
	return d.path
}

// Create implements gxui.TreeNode.
func (d directory) Create(theme gxui.Theme) gxui.Control {
	_, name := filepath.Split(d.path)
	if name == "" {
		name = d.path
	}
	l := theme.CreateLabel()
	l.SetText(name)
	l.SetColor(directoryColor)
	return l
}

// directoryAdapter is an implementation of the gxui.TreeAdapter interface.
// The AdapterItems returned by this adapter are absolute file path strings.
type directoryAdapter struct {
	gxui.AdapterBase
	directory
}

func (a directoryAdapter) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}

// Override directory.Create so that the full root is shown, unaltered.
func (a directoryAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	l := theme.CreateLabel()
	l.SetText(a.subdirs[index])
	l.SetColor(directoryColor)
	return l
}

func filePicker(theme gxui.Theme) gxui.Control {
	// fullpath is the textbox at the top of the window holding the current
	// selection's absolute file path.
	fullpath := theme.CreateTextBox()
	fullpath.SetDesiredWidth(math.MaxSize.W)

	// directories is the Tree of directories on the left of the window.
	// It uses the directoryAdapter to show the entire system's directory
	// hierarchy.
	directories := theme.CreateTree()
	directories.SetAdapter(&directoryAdapter{
		directory: directory{
			subdirs: roots.Roots(),
		},
	})

	// filesAdapter is the adapter used to show the currently selected directory's
	// content. The adapter has its data changed whenever the selected directory
	// changes.
	filesAdapter := &filesAdapter{}

	// files is the List of files in the selected directory to the right of the
	// window.
	files := theme.CreateList()
	files.SetAdapter(filesAdapter)

	upload := theme.CreateButton()
	upload.SetText("Выбрать файл")
	upload.OnClick(func(gxui.MouseEvent) {
		if files.Selected() != nil {
			filenameText = files.Selected().(string)
		}
	})

	// If the user hits the enter key while the fullpath control has focus,
	// attempt to select the directory.
	fullpath.OnKeyDown(func(ev gxui.KeyboardEvent) {
		if ev.Key == gxui.KeyEnter || ev.Key == gxui.KeyKpEnter {
			path := fullpath.Text()
			if directories.Select(path) {
				directories.Show(path)
			}
		}
	})

	// When the directory selection changes, update the files list
	directories.OnSelectionChanged(func(item gxui.AdapterItem) {
		dir := item.(string)
		filesAdapter.SetFiles(filesAt(dir))
		fullpath.SetText(dir)
	})

	// When the file selection changes, update the fullpath text
	files.OnSelectionChanged(func(item gxui.AdapterItem) {
		fullpath.SetText(item.(string))
		filenameText = item.(string)
	})

	// When the user double-clicks a directory in the file list, select it in the
	// directories tree view.
	files.OnDoubleClick(func(gxui.MouseEvent) {
		if path, ok := files.Selected().(string); ok {
			if fi, err := os.Stat(path); err == nil && fi.IsDir() {
				if directories.Select(path) {
					directories.Show(path)
				}
			} else {
				logger.Append("Пока не поддерживается выбор файла двойным кликом")
			}
		}
	})

	// Start with the CWD selected and visible.
	if cwd, err := os.Getwd(); err == nil {
		if directories.Select(cwd) {
			directories.Show(directories.Selected())
		}
	}

	splitter := theme.CreateSplitterLayout()
	splitter.SetOrientation(gxui.Horizontal)
	splitter.AddChild(directories)
	splitter.AddChild(files)

	subLayout := theme.CreateLinearLayout()
	subLayout.SetDirection(gxui.TopToBottom)
	subLayout.AddChild(fullpath)
	subLayout.AddChild(splitter)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.BottomToTop)
	layout.AddChild(subLayout)

	return layout
}
