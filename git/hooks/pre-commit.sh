#!/bin/sh
# Copyright 2012 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# git gofmt pre-commit hook
#
# To use, store as .git/hooks/pre-commit inside your repository and make sure
# it has execute permissions.
#
# This script does not handle file names that contain spaces.

echo 'Searching for Go files...'
gofiles=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
[ -z "$gofiles" ] && echo 'No Go files found.' && exit 0

echo 'Testing Go...'
gotest=$(go test -v $(go list ./... | grep -v /vendor/))
if [[ $(echo "${gotest}" | grep ^FAIL) != '' ]]; then
  echo "${gotest}"
  echo ''
  echo 'Some Go tests failed, please fix them.'
  exit 1
fi

echo 'Checking Go formatting.'
unformatted=$(gofmt -l $gofiles)
[ -z "$unformatted" ] && echo 'Go formatting is fine.' && exit 0

echo >&2 "Go files must be formatted with gofmt. Please run:"
for fn in $unformatted; do
	echo >&2 "  gofmt -w $PWD/$fn"
done

exit 1

