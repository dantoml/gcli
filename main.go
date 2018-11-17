package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// binName is the name that is currently being used by the binary.
	binName string

	// subcommandName is the command that has been requested. If the user did not
	// provide one, then it defaults to `help`
	subcommandName string

	// args are any arguments provided after the subcommand
	args []string
)

func init() {
	binName = filepath.Base(os.Args[0])

	if len(os.Args) > 1 {
		subcommandName = os.Args[1]
		args = os.Args[2:]
	} else {
		subcommandName = "help"
	}
}

func main() {
	subcommands := findSubcommands()

	// Find and execute the subcommand if possible. Doing this before checking for
	// help allows users to provide a custom help command.
	if cmd := subcommands.findNamed(subcommandName); cmd != nil {
		if err := cmd.exec(args); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n\n", err)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n", generateUsage(subcommands))
}

type commandList []*subcommand

func (c commandList) findNamed(name string) *subcommand {
	for _, cmd := range c {
		if cmd.name == name {
			return cmd
		}
	}

	return nil
}

func findBinaries(path, prefix string) fileList {
	var files fileList

	// Find all binaries in the users path :sob_sunglasses:
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}

		infos, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read directory %s: %s", dir, err)
			continue
		}

		for _, i := range infos {
			files = append(files, i)
		}
	}

	// Filter the files to executables that are matching the binary prefix.
	// We then unique these, because users may have duplicates in their path.
	return files.filesMatchingPrefix(prefix).executables().unique()
}

func findSubcommands() commandList {
	path := os.Getenv("PATH")
	prefix := fmt.Sprintf("%s-", binName)
	binaries := findBinaries(path, prefix)
	var cmds commandList

	for _, bin := range binaries {
		cmds = append(cmds, &subcommand{
			name: strings.TrimPrefix(filepath.Base(bin.Name()), prefix),
			path: bin.Name(),
		})
	}

	return cmds
}

func generateUsage(subcommands commandList) string {
	outtext := fmt.Sprintf("%s\n\nUsage:\n", binName)

	for _, c := range subcommands {
		str, err := c.usage()
		if err != nil {
			outtext += fmt.Sprintf("  %s - help not available: %s\n", c.name, err)
		} else {
			outtext += fmt.Sprintf("  %s - %s\n", c.name, str)
		}
	}

	return outtext
}
