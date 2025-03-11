# What is this Project?
One of the most annoying things to deal with in network administration is having to convert
between IP addresses and CIDR block/space notation.

As part of my learning/upskilling in Go as well as cloud certifications, I needed to learn
how to interpret CIDR ranges and convert between notations.  Instead of just reading about it,
I decided to make my life difficult and write a CLI tool and explore some of the Go packages that
exist for command-line processing.

# Features
Right now, this tool only converts CIDR blocks `Address/Network Prefix` to IP address ranges.

I would like to also:
- [ ] Convert IP address ranges to CIDR blocks
- [ ] Test if an IP exists within a given CIDR block
- [ ] Test if an IP exists within a provided IP address range

# Installation
## Build from source
### Prerequisites
To build from source, you will need the following available on your host:
- `git`
- `go` (1.23+)
### Building
You can build this project from source by downloading a release tag artifact or cloning this repo:
```
git clone git@github.com:taylrfnt/cidr-cli.git weather
```
After cloning the repo, delete the git contents:
```
cd cidr-cli && rm -rf .git
```
Now, run a Go build:
```
go build
```
You can execute `cidr-cli` directly:
```
./cidr-cli
```
or you can add it to your `$PATH` for easier use:
```
export PATH=${PATH}:/path/to/cidr-cli/directory
```

