package smbprotocol

import (
	"fmt"
	"io/fs"

	"github.com/hirochachacha/go-smb2"
	"github.com/ilightthings/shareblair/lib/options"
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
	Name            string
	Path            string
	ListOfFolders   []folder_A
	ListOfFiles     []file_A
	ReadAccess      bool
	WriteAccess     bool
	NumberOfFiles   int
	NumberOfFolders int
	NumberOfItems   int
}

type file_A struct {
	Name       string
	Path       string
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

func (s *Share) ListFilesRoot() error {
	list, err := s.Mount.ReadDir("")
	if err != nil {
		return err
	}
	s.UserRead = true
	s.ListOfFiles, s.ListOfFolders = sortFiles(list, fmt.Sprintf("\\\\%s", s.ShareName))

	return nil
}

func sortFiles(osfile []fs.FileInfo, CurrentPath string) ([]file_A, []folder_A) {
	var folders []folder_A
	var files []file_A
	NewPath := fmt.Sprintf("%s", CurrentPath)
	for _, x := range osfile {
		if x.IsDir() {
			var newfolder folder_A
			newfolder.Name = x.Name()
			newfolder.Path = fmt.Sprintf("%s\\%s\\", NewPath, x.Name())
			folders = append(folders, newfolder)
		} else {
			var newfile file_A
			newfile.Name = x.Name()
			newfile.Path = fmt.Sprintf("%s\\%s", NewPath, x.Name())
			files = append(files, newfile)
		}
	}
	return files, folders
}
