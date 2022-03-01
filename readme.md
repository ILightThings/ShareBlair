# shareblair
An SMB Scanner that scans hosts for SMB shares that are accessible to users and/or guests on the network.

Currently supports a single host scanning.

To build:
```
go build -o shareblair main.go
```

Usage Example:

**Target:** Localhost

**Username:** gameandwatch

**Password:** password

**Domain:** .\ (Local system)

```
.\shareblair -t 127.0.0.1 -u gameandwatch -p password -d .\
```

```
usage: Usage: [-h|--help] -t|--target "<value>" [-u|--user "<value>"]
              [-d|--domain "<value>"] [-p|--password "<value>"] [--hash
              "<value>"] [--port <integer>]



Arguments:

  -h  --help      Print help information
  -t  --target    IP Target to scan for SMB Shares
  -u  --user      User to authenticate with
  -d  --domain    Domain to authenticate with
  -p  --password  Password to authenticate with
      --hash      Hash to authenticate with
      --port      Port to connect to. Default: 445
```


## TODO
- [x] Guest Access Check 
- [ ] Test Read Write Access for every share found
- [ ] Read Targets in Cidr Notation
- [ ] Read Targets from file
- [x] Add ability to use hash instead of password
- [ ] Add Threading
- [ ] Output format
- [ ] HTML Output similar to dumpldapdomain   