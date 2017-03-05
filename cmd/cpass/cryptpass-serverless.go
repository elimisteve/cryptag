// Steve Phillips / elimisteve
// 2015.11.04

package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cryptag/cryptag/backend"
	"github.com/cryptag/cryptag/cli"
	"github.com/cryptag/cryptag/cli/color"
	"github.com/cryptag/cryptag/importer"
	"github.com/cryptag/cryptag/rowutil"
	"github.com/elimisteve/clipboard"
	shellwords "github.com/mattn/go-shellwords"
)

var (
	db backend.Backend
)

func init() {
	bkName := os.Getenv("BACKEND")

	var fs backend.Backend
	fs, err := backend.LoadOrCreateFileSystem(
		os.Getenv("BACKEND_PATH"),
		bkName,
	)
	if err != nil {
		if err != backend.ErrWrongBackendType {
			log.Fatalf("LoadOrCreateFileSystem error: %v\n", err)
		}

		// err == backend.ErrWrongBackendType

		fs, err = backend.LoadBackend("", bkName)
		if err != nil {
			log.Fatalf("Error loading Backend `%s` using LoadBackend: %v\n", bkName, err)
		}
	}

	db = fs
}

func main() {
	if len(os.Args) == 1 {
		cli.ArgFatal(allUsage)
	}

	switch os.Args[1] {
	case "create":
		if len(os.Args) < 4 {
			cli.ArgFatal(createUsage)
		}

		data := os.Args[2]
		tags := append(os.Args[3:], "app:cryptpass", "type:text")

		row, err := backend.CreateRow(db, nil, []byte(data), tags)
		if err != nil {
			log.Fatalf("Error creating then saving new row: %v", err)
		}

		color.Println(color.TextRow(row))

	case "tags":
		pairs, err := db.AllTagPairs(nil)
		if err != nil {
			log.Fatal(err)
		}

		for _, pair := range pairs {
			color.Printf("%s  %s\n", pair.Random, color.BlackOnWhite(pair.Plain()))
		}

	case "delete":
		if len(os.Args) < 3 {
			cli.ArgFatal(deleteUsage)
		}

		plaintags := append(os.Args[2:], "type:text")

		err := backend.DeleteRows(db, nil, plaintags)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Row(s) successfully deleted\n")

	case "run":
		if len(os.Args) < 3 {
			cli.ArgFatal(runUsage)
		}

		plaintags := append(os.Args[2:], "type:text", "type:command")

		rows, err := backend.RowsFromPlainTags(db, nil, plaintags)
		if err != nil {
			log.Fatal(err)
		}

		rows.Sort(rowutil.ByTagPrefix("created:", true))

		dec := rows[0].Decrypted()

		args, err := parse(string(dec))
		if err != nil {
			log.Fatalf("Error parsing command `%s`: %v\n", dec, err)
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin // Useful for `sudo ...` commands and the like
		var out bytes.Buffer
		cmd.Stdout = &out

		err = cmd.Run()
		if err != nil {
			log.Fatalf("Error running command `%s`: %v\n", dec, err)
		}

		err = clipboard.WriteAll(out.Bytes())
		if err != nil {
			log.Printf("WARNING: Error writing command output `%s`"+
				" to clipboard: %v\n", err)
		} else {
			log.Printf("Added output of first command\n\n"+
				"    $ %s\n\nto clipboard:\n\n", dec)
			color.Println(color.BlackOnCyan(string(out.Bytes())))
		}

		color.Println(color.TextRows(rows))

	case "import":
		if len(os.Args) < 3 {
			cli.ArgFatal(importUsage)
		}

		filename := os.Args[2]
		plaintags := os.Args[3:]

		rows, err := importer.KeePassCSV(filename, plaintags)
		if err != nil {
			log.Fatalf("Error importing KeePass CSV `%v`: %v", filename, err)
		}

		pairs, err := db.AllTagPairs(nil)
		if err != nil {
			log.Fatalf("Error fetching all TagPairs: %v\n", err)
		}

		for _, row := range rows {
			if _, err = backend.PopulateRowBeforeSave(db, row, pairs); err != nil {
				log.Printf("Error decrypting row %#v: %v\n", row, err)
				continue
			}
			if err := db.SaveRow(row); err != nil {
				log.Printf("Error saving row %#v: %v\n", row, err)
				continue
			}

			log.Printf("Successfully imported password for site %s\n",
				rowutil.TagWithPrefixStripped(row, "url:"))
		}

	default: // Search
		// Empty clipboard
		clipboard.WriteAll(nil)

		plaintags := append(os.Args[1:], "type:text")
		rows, err := backend.RowsFromPlainTags(db, nil, plaintags)
		if err != nil {
			log.Fatal(err)
		}

		rows.Sort(rowutil.ByTagPrefix("created:", true))

		// Add first row's contents to clipboard
		dec := rows[0].Decrypted()
		if err = clipboard.WriteAll(dec); err != nil {
			log.Printf("WARNING: Error writing first result to clipboard: %v\n", err)
		} else {
			log.Printf("Added first result `%s` to clipboard\n", dec)
		}

		color.Println(color.TextRows(rows))
	}
}

var (
	prefix = "Usage: " + filepath.Base(os.Args[0]) + " "

	createUsage = prefix + "create <password or text> <tag1> [type:password <tag3> ...]"
	tagsUsage   = prefix + "tags"
	deleteUsage = prefix + "delete <tag1> [<tag2> ...]"
	runUsage    = prefix + "run    <tag used to select command to run (commands are tagged with 'type:command')> [<tag1> ...]"
	importUsage = prefix + "import <exported-from-keepassx.csv> [<tag1> ...]"

	allUsage = strings.Join([]string{createUsage, tagsUsage, deleteUsage, runUsage, importUsage}, "\n")
)

func parse(cmd string) (args []string, err error) {
	p := shellwords.NewParser()
	p.ParseEnv = true
	p.ParseBacktick = true
	return p.Parse(cmd)
}
