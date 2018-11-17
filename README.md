# gcli

_git-style cli_ is a tiny CLI tool that makes it easy to create a program that
uses specialized subcommand binaries to perform specific tasks.

You probably shouldn't use this, it was written as a hack to simplify creating
GitHub Actions in the [endocrimes/actions](https://github.com/endocrimes/actions)
repository.

## Usage

Download the correct binary for your platform from the latest release, then
rename it to be the name of your desired program.

Add any sub command binaries to `$PATH` with a name in the style of `<cli>-<subcommand>`.
The sub command binaries must accept a `--short-help` flag that will return a
single line of help for printing in the generated Usage output.

