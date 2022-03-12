package report

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/ilightthings/shareblair/lib/options"
	"github.com/ilightthings/shareblair/lib/smbprotocol"
)

func GenerateReport(host *smbprotocol.Target, UserFlags *options.UserFlags) error {
	const SubEntry = `
	<style>
    table, th, td {
      border:1px solid black;
    }
    </style>
	<table>
	<tr>
	<th>UP</th>
	<th><a id={{.HumanPath}}>{{.HumanPath}}</a></th>
	</tr>
	{{range .ListOfFolders}}
	<tr>
	<td>📁</td>
	<td><a href="#{{.HumanPath}}">{{.Name}}</a></td>
	</tr>
	{{end}}
	{{range .ListOfFiles}}
	<tr>
	<td>📄</td>
	<td>{{.Name}}</td>
	</tr>
	{{end}}
	</table>
	<br>`

	//TODO Add function, if share is no access append "NO ACCESS" to share entry
	const ShareEntry = `
	
	<table>
	<tr>
	<th>UP</th>
	<th><h1><a id="{{.ShareName}}">Share {{.ShareName}}</a></h1></th>
	</tr>
	{{range .ListOfFolders}}
	<tr>
	<td>📁</td>
	<td><a href=#{{.HumanPath}}>{{.Name}}</a></td>
	</tr>
	{{end}}
	{{range .ListOfFiles}}
	<tr>
	<td>📄</td>
	<td>{{.Name}}</td>
	</tr>
	{{end}}
	</table>
	<br>`

	const header = `
	<style>
    table, th, td {
      border:1px solid black;
    }
    </style>
	<table>
	<tr><th>Share</th></tr>
	{{range .ListOfShares}}
	<tr><td><a href="#{{.ShareName}}">{{.ShareName}}</a></td></tr>
	{{end}}
	</table>
	<br>
	`

	//t, err := template.New("doc").Parse(MainTemplate)
	//if err != nil {
	//	return err
	//}

	//t.Execute(os.Stdout, s)

	var documentBytes []byte
	var headerBytes bytes.Buffer
	folderHTML, _ := template.New("table").Parse(SubEntry)
	headerHTML, _ := template.New("header").Parse(header)
	shareHTML, _ := template.New("Share").Parse(ShareEntry)
	headerHTML.Execute(&headerBytes, host)
	documentBytes = append(documentBytes, headerBytes.Bytes()...)

	for x := range host.ListOfShares {
		documentBytes = append(documentBytes, makeNewEntry(&host.ListOfShares[x], folderHTML, shareHTML, UserFlags)...)

	}
	os.WriteFile("file.html", documentBytes, 0644)

	return nil

}

func makeNewEntry(s *smbprotocol.Share, g *template.Template, h *template.Template, o *options.UserFlags) []byte {
	var buf bytes.Buffer
	var outBytes []byte
	if s.ListOfFolders == nil {
		return nil
	}
	for _, y := range s.ListOfFolders {
		makeNewEntryFunc(y, g, &buf, &outBytes)
	}

	err := h.Execute(&buf, &s)
	if err != nil {
		log.Panic(err)
	}

	var newOut []byte
	newOut = append(newOut, buf.Bytes()...)
	newOut = append(newOut, outBytes...)

	return newOut

}

func makeNewEntryFunc(f smbprotocol.Folder_A, g *template.Template, b *bytes.Buffer, bb *[]byte) {
	var appendBytes bytes.Buffer
	if len(f.ListOfFolders) != 0 {
		for _, x := range f.ListOfFolders {
			makeNewEntryFunc(x, g, b, bb)
		}
	}
	g.Execute(&appendBytes, f)

	//var newByteArray []byte
	var formerBytes []byte
	formerBytes = append(formerBytes, appendBytes.Bytes()...)
	formerBytes = append(formerBytes, *bb...)
	*bb = formerBytes

}
