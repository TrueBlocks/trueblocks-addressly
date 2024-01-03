package main

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

type Command struct {
	MaxRecords int
	Subcommand string
	Address    base.Address
	Filename   string
	Format     string
	Rest       string
	Silent     bool
	Chain      string
}

func (cmd *Command) String() string {
	listCmd := `chifra {{.Subcommand}} --cache {{.Address}} --fmt {{.Format}} --chain {{.Chain}} {{.Rest}}`
	if !cmd.Silent {
		listCmd += "  2>/dev/null"
	}
	if cmd.Address.IsZero() {
		listCmd = strings.Replace(listCmd, " --cache {{.Address}}", "", -1)
	} else if cmd.Subcommand == "list" {
		listCmd = strings.Replace(listCmd, " --cache", "", -1)
	}
	if len(cmd.Filename) > 0 {
		listCmd += " --output {{.Filename}}"
	}
	tmpl, err := template.New("greeting").Parse(listCmd)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, cmd)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	return buf.String()
}
