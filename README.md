#TinyPuters

A utility to manage raspberry pi and small IoT boards.

##Build

Requires go `> 1.15`

```bash
go build .
```

##Run

```bash
tinyputers
```

##Features

Currently discovers other devices on your network running tinyputers and stores them in a boltDB file, most of which is configurable.

TODO list:
- Add file based configuration
- Abstract DB
- Figure out if SSH is enabled on an rpi
- Share some details from  uname in the broadcast
- Establish a connection with found hosts
- Daemonize search/discovery
