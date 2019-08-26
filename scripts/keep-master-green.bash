#!/bin/bash
set -euo pipefail

reset=$(git rev-parse --abbrev-ref HEAD)
echo "Current HEAD: $reset"
master="remotes/origin/master"
if ! git merge-base --is-ancestor "$master" "$reset"; then
    echo "\`$master\` is not an ancestor of \`$reset\`; this script assumes that the current HEAD is directly based on current \`$master\`."
    exit 1
fi

if ! scripts/git-verify-no-diff.bash; then
    echo "Working directory should be clean before building all commits."
    exit 1
fi

dir=$(mktemp -d)
echo "Directory for build log files: $dir"
cmd="make all"
echo "Command used to build: \`$cmd\`"
range="$master..$reset"
echo "Testing the following commits:"
git --no-pager log --reverse --pretty='format:%H - %s' "$range"
echo "" # The above `git` command appears to not always output the last newline, which is done with this `echo` instead.
commits=$(git rev-list --reverse "$range")
for commit in $commits; do
    echo "$commit: checkout"
    git checkout --quiet "$commit"
    logfile="$dir/$commit.log"
    errfile="$dir/$commit.err"
    echo "$commit: build"
    set +e
    if ! $cmd 1>"$logfile" 2>"$errfile"; then
        echo "$commit: build failed"
        echo "Output of \`$cmd\` on standard out:"
        cat "$logfile"
        echo ""
        echo "Output of \`$cmd\` on standard err:"
        cat "$errfile"
        exit 1
    fi
    set -e
    echo "$commit: build passed"
done
echo "All commits from \`$master\` to \`$reset\` build successfully!"
echo "Removing directory for build log files."
rm -rf "$dir"
echo "Resetting to previous HEAD."
git checkout "$reset"
