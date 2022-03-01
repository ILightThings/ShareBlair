package smbprotocol

type Share struct {
	ShareName  string
	Hidden     bool // Anything with $ after the name is hidden
	UserRead   bool
	UserWrite  bool
	GuestRead  bool
	GuestWrite bool
	Folders    []folder
	Files      []file
}
