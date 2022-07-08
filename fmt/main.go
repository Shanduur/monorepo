package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	SUCCESS   = 0
	FAILURE   = 1
	FILENAME  = "^(.*)\\.(sql)$"
	INSERT    = "(INSERT INTO).*(\n.*)*(VALUES)"
	COMMENT   = "(^--.*)|([^']--.*)|(\\/\\*.*\\*\\/)"
	EMPTYLINE = "^( )*$"
)

var ErrWrongQuerry = errors.New("wrong format")

func generateSpace(size int) (out string) {
	if size <= 0 {
		return
	}

	for i := 0; i < size; i++ {
		out = out + " "
	}

	return out
}

func getQuerry(reader *bufio.Reader) (query string, comment []string, err error) {
	query, err = reader.ReadString(';')
	if err == io.EOF {
		return
	} else if err != nil {
		err = fmt.Errorf("unable to read string: %v", err)
		return
	}

	insertReg := regexp.MustCompile(INSERT)
	if insertReg.Match([]byte(query)) {
		commReg := regexp.MustCompile(COMMENT)
		bComment := commReg.FindAll([]byte(query), -1)
		for _, b := range bComment {
			c := strings.ReplaceAll(string(b), "\n", "")
			strings.ReplaceAll(c, "\r", "")

			comment = append(comment, c)
		}

		query = commReg.ReplaceAllString(query, "")
		query = strings.ReplaceAll(query, "\n", "")
		query = strings.ReplaceAll(query, "\r", "")
	} else {
		query = strings.TrimSpace(query)
		query = strings.Trim(query, "\n")
		query = strings.Trim(query, "\r")
	}

	return
}

func connectWronglySplitted(splits []string) (out string) {
	for _, s := range splits {
		if len(out) > 0 {
			out = fmt.Sprintf("%s,%s", out, s)
		} else {
			out = fmt.Sprintf("%s", s)
		}
	}

	return
}

func tryMerge(values []string) (merged []string) {
	for i := 0; i < len(values); i++ {
		if strings.Count(values[i], "'") == 1 {
			for j := i + 1; j < len(values); j++ {
				if strings.Contains(values[j], "'") {
					merged = append(merged, connectWronglySplitted(values[i:j+1]))
					i = j
					break
				}
			}
		} else {
			merged = append(merged, values[i])
		}
	}

	return
}

func formatQuerry(query string) (formatted string, err error) {
	format := "INSERT INTO %v\n" +
		"     VALUES %v\n"

	split := strings.Split(query, "VALUES")
	if len(split) < 2 {
		err = ErrWrongQuerry
		formatted = query
		return
	}

	split[0] = strings.Replace(split[0], "INSERT INTO", "", 1)
	splitCols := strings.Split(split[0], ",")
	splitVals := strings.Split(split[1], ",")

	if len(splitCols) != len(splitVals) {
		splitVals = tryMerge(splitVals)

		if len(splitCols) != len(splitVals) {
			err = ErrWrongQuerry
			formatted = "-- CHECK FORMAT!\n" + query + "\n-- CHECK FORMAT!"
			return
		}
	}

	var (
		columns string
		values  string
	)

	headerCols := strings.Split(splitCols[0], "(")

	columns += strings.TrimSpace(headerCols[0]) + "(" + strings.TrimSpace(headerCols[1])
	values += strings.TrimSpace(splitVals[0])
	columns += generateSpace(len(values) - len(columns))
	values = generateSpace(len(columns)-len(values)) + values

	for i := 1; i < len(splitCols); i++ {

		columns += ", " + strings.TrimSpace(splitCols[i])
		values += ", " + strings.TrimSpace(splitVals[i])

		columns += generateSpace(len(values) - len(columns))
		values += generateSpace(len(columns) - len(values))
	}

	formatted = fmt.Sprintf(format, columns, values)

	return
}

func format(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()
	defer file.Sync()

	scanner := bufio.NewReader(file)
	if err != nil {
		return fmt.Errorf("unable to create scanner on file: %v", err)
	}

	regEmptyLine := regexp.MustCompile(EMPTYLINE)
	var out []string
	for {
		query, comment, err := getQuerry(scanner)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("unable to get query: %v", err)
		}

		if len(comment) > 0 {
			for _, c := range comment {
				c = regEmptyLine.ReplaceAllString(c, "")
				out = append(out, fmt.Sprintf("%v\n", c))
			}
		}

		if len(query) > 0 {
			query, err = formatQuerry(query)
			if err == ErrWrongQuerry {
				out = append(out, fmt.Sprintf("%v\n", query))
				continue
			} else if err != nil {
				return fmt.Errorf("unable to format query: %v", err)
			}

			out = append(out, fmt.Sprintf("%v", query))
		}
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("unable to truncate file: %v", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("unable to seek begining of file: %v", err)
	}

	writer := bufio.NewWriter(file)
	for _, o := range out {
		_, err = writer.WriteString(o)
		if err != nil {
			return fmt.Errorf("unable to write to a file: %v", err)
		}
		writer.Flush()
	}

	return nil
}

type Status struct {
	ok  bool
	mux sync.Mutex
}

func NewStatus() *Status {
	return &Status{
		ok: true,
	}
}

func (s *Status) Failed() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.ok = false
}

func (s *Status) OK() bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.ok
}

func wrap(x func(string) error, f string, wg *sync.WaitGroup, stat *Status) {
	defer wg.Done()

	err := x(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to format file: %v\n", err)
		stat.Failed()
	}
}

func timed(x func(string) error) func(string) error {
	return func(f string) error {
		fmt.Printf("formatting: %s\n", f)

		start := time.Now()
		err := x(f)
		fmt.Printf("%s done, took: %v\n", f, time.Since(start))
		return err
	}
}

func main() {
	path := "."
	if len(os.Args) >= 2 {

		path = os.Args[1]
	}

	si, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get os stat: %v\n", err)
	}

	mode := si.Mode()

	switch {
	case mode.IsDir():
		var files []string

		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to format file: %v\n", err)
			os.Exit(FAILURE)
		}

		regFileName := regexp.MustCompile(FILENAME)

		status := NewStatus()
		var wg sync.WaitGroup
		for _, f := range files {
			if !regFileName.MatchString(f) {
				continue
			}

			wg.Add(1)
			go wrap(timed(format), f, &wg, status)
		}

		wg.Wait()
		if !status.OK() {
			os.Exit(FAILURE)
		}

	case mode.IsRegular():
		regFileName := regexp.MustCompile(FILENAME)
		if !regFileName.MatchString(path) {
			fmt.Fprintf(os.Stderr, "wrong file name\n")
			os.Exit(FAILURE)
		}

		err = timed(format)(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to format file: %v\n", err)
			os.Exit(FAILURE)
		}
	}

	os.Exit(SUCCESS)
}
