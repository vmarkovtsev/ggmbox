package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/araddon/dateparse"
	"github.com/jhillyerd/enmime"
)

const (
	malformedMIMEHeaderLineErrorMessage = "malformed MIME header line: "
)

var (
	base64RE = regexp.MustCompile("^([a-zA-Z0-9+/]+\\r\\n)+[a-zA-Z0-9+/]+={0,2}$")
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("missing argument")
	}
	dir := os.Args[1]
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("could not read directory %s: %v", dir, err)
	}

	type result struct {
		path string
		env  *enmime.Envelope
		err  error
	}
	results := make(chan result, 1024)

	var wgFiles, wgProcess sync.WaitGroup
	wgFiles.Add(1)
	go func() {
		defer wgFiles.Done()
		filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				log.Printf("reading %s: %v", path, err)
				return nil
			}
			if fi.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".email") {
				return nil
			}
			wgFiles.Add(1)
			go func(path string) {
				defer wgFiles.Done()
				env, err := readEnvelope(path)
				results <- result{path, env, err}
			}(path)
			return nil
		})
	}()

	archive := map[string][]*enmime.Envelope{}
	var nenvs, nerr int
	wgProcess.Add(1)
	go func() {
		defer wgProcess.Done()
		for res := range results {
			if res.err != nil {
				log.Printf("processing %s: %v", res.path, res.err)
				nerr++
				continue
			}

			nenvs++
			topic := filepath.Base(filepath.Dir(res.path))
			archive[topic] = append(archive[topic], res.env)
		}
	}()

	wgFiles.Wait()
	close(results)
	wgProcess.Wait()
	nthreads := len(archive)

	fmt.Fprintf(os.Stderr, "Topics: %d\n", nthreads)
	fmt.Fprintf(os.Stderr, "Errors: %d\n", nerr)
	fmt.Fprintf(os.Stderr, "Average topic size: %.1f\n", float64(nenvs)/float64(nthreads))

	// Extract plain text threads and write CSV to stdout
	writePlainTextTopics(archive)
}

func readEnvelope(path string) (*enmime.Envelope, error) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", path, err)
	}
	env, err := enmime.ReadEnvelope(bytes.NewReader(text))
	for i := 0; i < 4 && err != nil; i++ {
		errStr := err.Error()
		if pos := strings.Index(errStr, malformedMIMEHeaderLineErrorMessage); pos >= 0 {
			data := strings.Replace(errStr[pos+len(malformedMIMEHeaderLineErrorMessage):], " ", "\r\n", -1)
			if base64RE.MatchString(data) {
				var decodedData []byte
				decodedData, err = base64.StdEncoding.DecodeString(data)
				if err == nil {
					text = bytes.Replace(text, []byte(data), []byte(decodedData), 1)
					env, err = enmime.ReadEnvelope(bytes.NewReader(text))
					if err == nil {
						return env, nil
					}
				} else {
					break
				}
			}
		} else {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("could not read envelope: %v", err)
	}
	return env, nil
}

func writePlainTextTopics(archive map[string][]*enmime.Envelope) {
	fmt.Println("topic_id,text")
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()
	seq := make([]string, 2)
	keys := make([]string, len(archive))
	{
		i := 0
		for key := range archive {
			keys[i] = key
			i++
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		envs := archive[key]
		seq[0] = key
		sort.Slice(envs, func(i, j int) bool {
			time1, err := dateparse.ParseAny(envs[i].GetHeader("Date"))
			if err != nil {
				panic(err)
			}
			time2, err := dateparse.ParseAny(envs[j].GetHeader("Date"))
			if err != nil {
				panic(err)
			}
			return time1.Before(time2)
		})
		thread := new(bytes.Buffer)
		for i, env := range envs {
			cleanupMessage(env.Text, thread)
			if i < len(envs)-1 {
				thread.Write([]byte("\n"))
			}
		}
		seq[1] = thread.String()
		writer.Write(seq)
	}
}

func cleanupMessage(text string, output io.Writer) {
	lines := strings.Split(text, "\n")
	first := true
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		line = line[:strings.Index(line, trimmedLine)+len(trimmedLine)]
		if strings.HasPrefix(line, ">") {
			continue
		}
		if strings.HasPrefix(line, "On ") && strings.HasSuffix(line, " wrote:") {
			continue
		}
		if !first {
			output.Write([]byte("\n"))
		} else {
			first = false
		}
		output.Write([]byte(line))
	}
}
