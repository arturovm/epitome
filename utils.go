package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"text/tabwriter"
)

func WriteJSONError(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error": "` + errorMessage + `"}`))
}

func TestContentType(rawHeader *string, target string) bool {
	parts := strings.Split(*rawHeader, ";")
	if len(parts) > 0 {
		return strings.Trim(parts[0], " \t") == target
	}
	return false
}

func LogRequest(req *http.Request, route string) {
	log.SetOutput(os.Stdout)
	reqB, _ := httputil.DumpRequest(req, verboseModeBody)
	log.Print("Received request at " + "'" + route + "'\n\n" + string(reqB) + "\n\n")
	log.SetOutput(os.Stderr)
}

func LogMap(m map[string]string, w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 8, 1, '\t', 0)
	for i, k := range m {
		fmt.Fprintln(tw, i+":\t"+k)
	}
	tw.Flush()
}

func LogParsedValues(m map[string]string, w io.Writer) {
	fmt.Fprintln(w, "### Parsed values ###\n")
	LogMap(m, w)
	fmt.Fprint(w, "\n\n")
}
