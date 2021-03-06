package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

func perform() error {
	if len(os.Args) != 4 {
		return fmt.Errorf("must pass 3 arguments")
	}

	name := os.Args[1]
	source := os.Args[2]
	destination := os.Args[3]

	rows, err := readFile(source)
	if err != nil {
		return err
	}
	params, ok := paramsMap[name]
	if !ok {
		return fmt.Errorf("%q not supported", name)
	}
	return writeTempalte(destination, rows, params)
}

func main() {
	if err := perform(); err != nil {
		log.Fatal(err)
	}
}

func readFile(path string) ([][]string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close() // nolint: errcheck

	var rows [][]string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		rows = append(rows, strings.Split(line, "|"))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rows, nil
}

var paramsMap = map[string]*templateParams{
	// FIPS|USPS|STATE_NAME|GNISID
	"states": &templateParams{
		Preface: `type State struct {
            FIPS string
            USPS string
            Name string
            GNISID string
        }`,
		Var:  "States",
		Type: "[]State",
		ValFn: func(row []string) string {
			return fmt.Sprintf(`State{%q, %q, %q, %q}`,
				row[0], row[1], row[2], row[3])
		},
	},
	// USPS|GEOID|ANSICODE|NAME|ALAND|AWATER|ALAND_SQMI|AWATER_SQMI|INTPTLAT|INTPTLONG
	"counties": &templateParams{
		Preface: `type County struct {
            USPS string
            GEOID string
            ANSICode string
            Name string
        }`,
		Var:  "Counties",
		Type: "[]County",
		ValFn: func(row []string) string {
			return fmt.Sprintf(`County{%q,%q,%q,%q}`,
				row[0], row[1], row[2], row[3])
		},
	},
}

var tpl = template.Must(template.
	New("file").
	Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package generated

{{ .Params.Preface }}

var {{ .Params.Var }} = {{ .Params.Type }}{
{{- range .Rows }}
    {{ call $.Params.ValFn . }},
{{- end }}
}
`))

type templateParams struct {
	Preface string
	Var     string
	Type    string
	ValFn   func(row []string) string
}

func writeTempalte(path string, rows [][]string, params *templateParams) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	err = tpl.Execute(f, struct {
		Timestamp time.Time
		Rows      [][]string
		Params    *templateParams
	}{
		Timestamp: time.Now(),
		Rows:      rows,
		Params:    params,
	})

	if err != nil {
		// os.Remove(path) // nolint: errcheck,gas
		return err
	}
	return nil
}
