.PHONY: init tests release

#
# Install and replace required dependencies.
#
init:
	bash ./scripts/mod-tidy.sh

#
# Run tests across all Go modules/packages.
#
tests:
	bash ./scripts/tests.sh

#
# Commit, tag, release, and update dependencies when releasing a new version.
#
release:
	bash ./scripts/release.sh
