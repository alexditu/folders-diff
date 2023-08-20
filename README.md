# folders-diff
Recursively traverse two directories and print the files from 1st directory that are not present in the 2nd directory

Created this short program to improve my golang knwoledge while also helping a friend to double check which files weren't backed up from his old hard drives :).

## Usage
```
fdiff - Recursively traverse folderA and folderB and prints the files that are in folderA and not in folderB (set difference folderA - folderB)

Usage:
  fdiff <folderA> <folderB> [flags]

Flags:
  -h, --help      help for fdiff
  -V, --verbose   enable verbose logging
  -v, --version   version for fdiff
```

## How to Build
Run `build.sh` script, it will create binaries for windows and linux amd64.
