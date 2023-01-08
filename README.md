[![codecov](https://codecov.io/github/fortio/delta/branch/main/graph/badge.svg?token=LONYZDFQ7C)](https://codecov.io/github/fortio/delta)

# Delta
Diff 2 sets and apply command to deltas

## Installation

go install github.com/fortio/delta@latest

## Usage

delta -b "echo NEW:" -a "echo REMOVED:" oldFile newFile

if `oldFile` is
```
old1
old2
gone1
old3
```

and `newFile` is
```
new1
old1
old2
old3
new2
```

will output
```
REMOVED: gone1
NEW: new1
NEW: new2
```

See also [delta.txtar](delta.txtar) for examples (tests)
