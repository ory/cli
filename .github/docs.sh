#!/usr/bin/env bash

set -euxo pipefail

mkdir -p docs/docs
bash <(curl -s https://raw.githubusercontent.com/ory/ci/master/src/scripts/install/git.sh)
make docs/cli

docsdir=$(mktemp -d)
git clone --depth 1 git@github.com:ory/docs.git "$docsdir"
/bin/cp -rf docs "$docsdir"
cd "$docsdir"
git add -A
git commit -a -m "autogen(docs): regenerate ory/cli reference" --allow-empty
git pull -ff
git push

