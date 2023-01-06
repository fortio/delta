# Delta
Diff 2 sets and apply command to deltas

## Installation

go install github.com/fortio/delta@latest

## Usage

delta -new "echo NEW: " -gone "echo REMOVED: " oldFile newFile

if `old` is
```
old1
old2
gone1
old3
```

and `new` is
```
new1
old1
old2
old3
new2
```

will output
```
NEW: new1
NEW: new2
REMOVED: gone1
```
