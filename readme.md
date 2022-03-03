# shareblair
An SMB Scanner that scans hosts for SMB shares that are accessible to users and/or guests on the network. The recursive scanner allows you as the tester to find files that may normaly get missed in the normal enumeration phase. 

Currently supports a Hostname, IP, CIDR Notation, or File input with the single tag `-t`. ShairBlare also allows NTLM hashes for "Pass The HasH".

At the moment, Shairblair outpus in json format for parseing. We are currently working on exporting in an HTML interactive format.

## Usage 
```
usage: shareblair [-h|--help] -t|--target "<value>" [-u|--user "<value>"]
                  [-d|--domain "<value>"] [-p|--password "<value>"] [--hash
                  "<value>"] [--port <integer>] [-v|--verbose] [--maxdepth
                  <integer>]



Arguments:

  -h  --help      Print help information
  -t  --target    Hostname, IP, CIDR or file of targets
  -u  --user      User to authenticate with
  -d  --domain    Domain to authenticate with
  -p  --password  Password to authenticate with
      --hash      Hash to authenticate with
      --port      Port to connect to. Default: 445
  -v  --verbose   Add verbosity. Default: false
      --maxdepth  Max Recursive Depth for Share Scanning. 0 will only scan the
                  top level folders.. Default: 5
                  
```

## Example
```
go run main.go -t 192.168.1.159 -u gameandwatch --hash 8846F7EAEE8FB117AD06BDD830B7586C
```
### Output - 192.168.1.159_SMBShares.json
```json
{
    "HostDestination": "192.168.1.159",
    "ResolvedIP": "",
    "UserFlag": {
        "Target": "192.168.1.159",
        "Threads": 0,
        "Verbose": true,
        "User": "gameandwatch",
        "Domain": "",
        "Password": "",
        "Hash": "8846F7EAEE8FB117AD06BDD830B7586C",
        "Port": 445,
        "MaxDepth": 5,
        "OutFileLocation": ""
    },
    "ConnectionTCP": {},
    "ConnectionTCP_OK": true,
    "ConnectionSMB": {},
    "ConnectionSMB_OK": true,
    "GuestOnly": false,
    "GuestAccess": false,
    "ListOfShares": [
        {
            "ShareName": "Test4",
            "Hidden": false,
            "SMBConnection": {},
            "Mount": {},
            "Mounted": false,
            "UserFlags": {
                "Target": "192.168.1.159",
                "Threads": 0,
                "Verbose": true,
                "User": "gameandwatch",
                "Domain": "",
                "Password": "",
                "Hash": "8846F7EAEE8FB117AD06BDD830B7586C",
                "Port": 445,
                "MaxDepth": 5,
                "OutFileLocation": ""
            },
            "UserRead": true,
            "UserWrite": false,
            "GuestRead": false,
            "GuestWrite": false,
            "ListOfFolders": [
                {
                    "Depth": 0,
                    "Name": "InnerFolder",
                    "HumanPath": "\\\\192.168.1.159\\Test4\\InnerFolder",
                    "ListOfFolders": [
                        {
                            "Depth": 1,
                            "Name": "Folder2",
                            "HumanPath": "\\\\192.168.1.159\\Test4\\InnerFolder\\Folder2",
                            "ListOfFolders": [
                                {
                                    "Depth": 2,
                                    "Name": "SysinternalsSuite",
                                    "HumanPath": "\\\\192.168.1.159\\Test4\\InnerFolder\\Folder2\\SysinternalsSuite",
                                    "ListOfFolders": null,
                                    "ListOfFiles": [
                                        {
                                            "Name": "ADExplorer.exe",
                                            "HumanPath": "\\\\192.168.1.159\\Test4\\InnerFolder\\Folder2\\SysinternalsSuite\\ADExplorer.exe",
                                            "FolderPath": "\\\\192.168.1.159\\Test4\\InnerFolder\\Folder2\\SysinternalsSuite",
                                            "Size": 1237896
                                        }
                                        ...SNIP...
```


---

## TODO
- [x] Guest Access Check 
- [ ] Test Read Write Access for every share found
- [x] Read Targets in Cidr Notation
- [x] Read Targets from file
- [x] Add ability to use hash instead of password
- [ ] Add Threading
- [ ] Output format
- [-] HTML Output similar to dumpldapdomain   (IN PROGRESS)