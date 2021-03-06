package sqltc

import (
	"io/ioutil"
	"strings"
)

type SqlFile struct {
	files   []string
	Queries []string
}

type Column struct {
	Name   string
	Type   string
	IsNULL bool
}

type Columns []Column

const (
	// number type
	BIT = iota
	TINYINT
	BOOL
	BOOLEAN
	SMALLINT
	MEDIUMINT
	INT
	INTEGER
	BIGINT
	DECIMAL
	DEC
	FLOAT
	DOUBLE

	// date
	DATE
	DATETIME
	TIMESTAMP
	TIME
	YEAR

	// string type
	CHAR
	VARCHAR
	BINARY
	VARBINARY
	TINYBLOB
	TINYTEXT
	BLOB
	TEXT
	MEDIUMBLOB
	MEDIUMTEXT
	LONGBLOB
	LONGTEXT
	ENUM
	SET
)

var tokens = [...]string{
	BIT:       "BIT",
	TINYINT:   "TINYINT",
	BOOL:      "BOOL",
	BOOLEAN:   "BOOLEAN",
	SMALLINT:  "SMALLINT",
	MEDIUMINT: "MEDIUMINT",
	INT:       "INT",
	INTEGER:   "INTEGGER",
	BIGINT:    "BIGINT",
	DECIMAL:   "DECIMAL",
	DEC:       "DEC",
	FLOAT:     "FLOAT",
	DOUBLE:    "DOUBLE",

	// date
	DATE:      "DATE",
	DATETIME:  "DATETIME",
	TIMESTAMP: "TIMESTAMP",
	TIME:      "TIME",
	YEAR:      "YEAR",

	// string type
	CHAR:       "CHAR",
	VARCHAR:    "VARCHAR",
	BINARY:     "BINARY",
	VARBINARY:  "VARBINARY",
	TINYBLOB:   "TINYBLOB",
	TINYTEXT:   "TINYTEXT",
	BLOB:       "BLOB",
	TEXT:       "TEXT",
	MEDIUMBLOB: "MEDIUMBLOB",
	MEDIUMTEXT: "MEDIUMTEXT",
	LONGBLOB:   "LONGBLOB",
	LONGTEXT:   "LONGTEXT",
	ENUM:       "ENUM",
	SET:        "SET",
}

func Convert(query string) Columns {
	columns := Columns{}
	sc := strings.Split(query, "  ")
	for _, v := range sc {
		column := Column{}
		column.IsNULL = true
		column.Name = strings.Split(v, " ")[0]
		for _, token := range tokens {
			if !strings.Contains(v, token) {
				continue
			}
			column.Type = token
			if strings.Contains(v, "NOT NULL") {
				column.IsNULL = false
			}
		}
		if column.Type == "" {
			continue
		}
		columns = append(columns, column)
	}
	return columns
}

func (s *SqlFile) Directory(dir string) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if name[len(name)-3:] != "sql" {
			continue
		}

		if err := s.File(dir + "/" + name); err != nil {
			return err
		}
	}

	return nil
}

func (s *SqlFile) Files(files ...string) error {

	for _, file := range files {
		if err := s.File(file); err != nil {
			return err
		}
	}

	return nil
}

func (s *SqlFile) File(file string) error {

	queries, err := load(file)
	if err != nil {
		return err
	}

	s.files = append(s.files, file)
	s.Queries = append(s.Queries, queries...)

	return nil
}

func load(path string) ([]string, error) {

	lines, err := lineFromFile(path)
	if err != nil {
		return nil, err
	}

	var ls []string
	for _, line := range lines {
		l := excludeComment(line)
		ls = append(ls, l)
	}

	line := strings.Join(ls, "")
	queries := strings.Split(line, ";")
	queries = queries[:len(queries)-1]
	return queries, nil
}

func lineFromFile(path string) (ls []string, err error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return ls, err
	}

	ls = strings.Split(string(f), "\n")
	return ls, nil
}

func excludeComment(line string) string {
	d := "\""
	s := "'"
	c := "--"

	var nc string
	ck := line
	mx := len(line) + 1

	for {
		if len(ck) == 0 {
			return nc
		}

		di := strings.Index(ck, d)
		si := strings.Index(ck, s)
		ci := strings.Index(ck, c)

		if di < 0 {
			di = mx
		}
		if si < 0 {
			si = mx
		}
		if ci < 0 {
			ci = mx
		}

		var ei int

		if di < si && di < ci {
			nc += ck[:di+1]
			ck = ck[di+1:]
			ei = strings.Index(ck, d)
		} else if si < di && si < ci {
			nc += ck[:si+1]
			ck = ck[si+1:]
			ei = strings.Index(ck, s)
		} else if ci < di && ci < si {
			return nc + ck[:ci]
		} else {
			return nc + ck
		}

		nc += ck[:ei+1]
		ck = ck[ei+1:]
	}
}
