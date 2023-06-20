#!/usr/bin/env bash

export GOPRIVATE="go.nunchi.studio/*"

go work use -r ./

rm -rf go.sum
go mod tidy

integrations=$( cat ./ecosystem.json | jq -r '.integrations[] | select(.soon|not) | .id' | cat )
for mod in $integrations; do
  cd ./integration/$mod

  rm -rf go.sum
  go mod tidy

  cd ../../
done

go work sync
