#!/usr/bin/env bash

if [[ -z "$GITHUB_TOKEN" ]]; then
  echo "Environment variable `GITHUB_TOKEN` must be set and not be empty" 1>&2
  exit 1
fi

export GORELEASER_CURRENT_TAG=v0.17.0

go work use -r ./

go mod tidy

integrations=$( cat ./ecosystem.json | jq -r '.integrations[] | select(.soon|not) | .id' | cat )
for mod in $integrations; do
  cd ./integration/$mod

  go mod tidy

  ./scripts/mod-require.sh $GORELEASER_CURRENT_TAG
  ./scripts/mod-dropreplace.sh

  cd ../../
done

git commit -s -am "version: Release $GORELEASER_CURRENT_TAG"
git push origin main

for mod in $integrations; do
  git tag "integration/$mod/$GORELEASER_CURRENT_TAG"
done

git tag "$GORELEASER_CURRENT_TAG"
git push --tags

goreleaser release --rm-dist

for mod in $integrations; do
  cd ./integration/$mod

  ./scripts/mod-replace.sh

  cd ../../
done

git commit -s -am "version(scripts): Apply post-release of $GORELEASER_CURRENT_TAG"
git push origin main
