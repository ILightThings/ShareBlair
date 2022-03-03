package smbprotocol

import (
	"fmt"
	"io/fs"

	"github.com/hirochachacha/go-smb2"
	"github.com/ilightthings/shareblair/lib/options"
)

const (
	CONTINUE             = 0
	MAX_DEPTH_STOP       = 1
	PERMISSION_DENY_STOP = 2
	OTHER                = 3
	NO_MORE_FOLDERS      = 4
)

type Share struct {
	ShareName     string
	Hidden        bool // Anything with $ after the name is hidden
	SMBConnection *smb2.Session
	Mount         *smb2.Share
	Mounted       bool
	UserFlags     *options.UserFlags
	UserRead      bool
	UserWrite     bool
	GuestRead     bool
	GuestWrite    bool
	ListOfFolders []folder_A
	ListOfFiles   []file_A
}

type folder_A struct {
	Depth           int
	Name            string
	Path            string
	FullPath        string
	HumanPath       string
	ListOfFolders   []folder_A
	ListOfFiles     []file_A
	ReadAccess      bool
	WriteAccess     bool
	NumberOfFiles   int
	NumberOfFolders int
	NumberOfItems   int
	Stop_reason     int
}

type file_A struct {
	Name       string
	Path       string
	FullPath   string
	HumanPath  string
	FolderPath string
	FilePath   string
	FileName   string
	Size       int
}

func (s *Share) InitializeShare(q *smb2.Session, f *options.UserFlags) error {
	s.SMBConnection = q
	s.UserFlags = f
	var err error
	if s.UserFlags.Verbose {
		fmt.Printf("Attempting to mount %s\n", s.ShareName)
	}
	s.Mount, err = s.SMBConnection.Mount(s.ShareName)
	if err != nil {
		if s.UserFlags.Verbose {
			fmt.Printf("Failed to mount %s\n", s.ShareName)
		}
		return err
	}
	s.Mounted = true
	if s.UserFlags.Verbose {
		fmt.Printf("Successfully mounted %s\n", s.ShareName)
	}
	s.isHidden()
	return nil

}

func (s *Share) UnmountShare() {
	if s.Mount != nil {
		s.Mount.Umount()
	}
	s.Mounted = false
}

func (s *Share) isHidden() {
	if s.ShareName[len(s.ShareName)-1:] == "$" {
		s.Hidden = true
	}
}

func (s *Share) ListFilesRoot(host string) error {
	list, err := s.Mount.ReadDir("")
	if err != nil {
		return err
	}
	s.UserRead = true
	s.ListOfFiles, s.ListOfFolders = sortFiles(list, "", 0, s, host)

	return nil
}

func sortFiles(osfile []fs.FileInfo, CurrentPath string, depth int, s *Share, host string) ([]file_A, []folder_A) {
	var folders []folder_A
	var files []file_A

	for _, x := range osfile {
		if x.IsDir() { // Define Folder Properties Here

			var newfolder folder_A
			newfolder.Name = x.Name()
			if CurrentPath == "" {
				newfolder.Path = x.Name()
			} else {
				newfolder.Path = fmt.Sprintf("%s\\%s", CurrentPath, x.Name())
			}
			newfolder.FullPath = fmt.Sprintf("\\%s\\%s", s.ShareName, newfolder.Path)
			newfolder.HumanPath = fmt.Sprintf("\\\\%s%s", host, newfolder.FullPath)
			newfolder.Depth = depth

			folders = append(folders, newfolder)
		} else { // Define File Properties Here
			var newfile file_A
			newfile.Name = x.Name()
			if CurrentPath == "" {
				newfile.Path = x.Name()
			} else {
				newfile.Path = fmt.Sprintf("%s\\%s", CurrentPath, x.Name()) // folder1\(name of folder)
			}

			newfile.FullPath = fmt.Sprintf("\\%s\\%s", s.ShareName, newfile.Path)
			newfile.HumanPath = fmt.Sprintf("\\\\%s%s", host, newfile.FullPath)
			newfile.FolderPath = fmt.Sprintf("\\\\%s\\%s\\%s", host, s.ShareName, CurrentPath)
			files = append(files, newfile)
		}
	}
	return files, folders
}

func walkDirFn(currentFolder *folder_A, depth int, s *Share, host string) error {
	if depth >= s.UserFlags.MaxDepth {
		currentFolder.Stop_reason = MAX_DEPTH_STOP
		return nil
	}
	FolderFiles, err := s.Mount.ReadDir(currentFolder.Path)
	if err != nil {
		currentFolder.Stop_reason = PERMISSION_DENY_STOP
		return err
	}
	currentFolder.ListOfFiles, currentFolder.ListOfFolders = sortFiles(FolderFiles, currentFolder.Path, depth+1, s, host)
	if len(currentFolder.ListOfFolders) == 0 {
		currentFolder.Stop_reason = NO_MORE_FOLDERS
	}

	for x, _ := range currentFolder.ListOfFolders {
		walkDirFn(&currentFolder.ListOfFolders[x], depth+1, s, host)

	}

	return nil

}

func (s *Share) DirWalk(host string) error {
	TopLevelFoldersFile, err := s.Mount.ReadDir("")
	if err != nil {
		s.UserRead = false
		return err
	} else {
		s.UserRead = true
	}
	s.ListOfFiles, s.ListOfFolders = sortFiles(TopLevelFoldersFile, "", 0, s, host)
	for x, _ := range s.ListOfFolders {
		walkDirFn(&s.ListOfFolders[x], 0, s, host)
	}
	return nil

}

// https://cs.opensource.google/go/go/+/refs/tags/go1.17.7:src/path/filepath/path.go;drc=refs%2Ftags%2Fgo1.17.7;l=385
// Basically instead of scanning each folder to completeion one by one, I can one folder until I get to a folder, start a new process to scan that folder, and when a folder has no more child folders, return to previous folder scan.
