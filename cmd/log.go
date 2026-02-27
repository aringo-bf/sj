package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Result struct {
	Method string `json:"method"`
	Status int    `json:"status"`
	Target string `json:"target"`
}

type VerboseResult struct {
	Method      string      `json:"method"`
	Preview     interface{} `json:"preview"`
	Status      int         `json:"status"`
	ContentType string      `json:"content_type"`
	Target      string      `json:"target"`
	Curl        string      `json:"curl"`
}

func printInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprint(os.Stderr, msg)
}

func printWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, "[!] %s", msg)
}

func printErr(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, "[x] %s", msg)
}

func die(format string, args ...interface{}) {
	printErr(format, args...)
	os.Exit(1)
}

func writeLog(sc int, target, method, errorMsg, response string) {
	var out io.Writer = os.Stdout
	tempResponsePreviewLength := responsePreviewLength

	if len(response) < responsePreviewLength {
		responsePreviewLength = len(response)
	}

	if outfile != "" {
		file, err := os.OpenFile(outfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			die("Output file does not exist or cannot be created")
		}
		defer file.Close()
		out = file
	}

	preview := ""
	if verbose && responsePreviewLength > 0 {
		preview = response[:responsePreviewLength]
	}

	switch sc {
	case 8899:
		if verbose {
			logVerboseJSON(specTitle, specDescription, out)
		} else {
			logJSON(specTitle, specDescription, out)
		}
	default:
		logResult(sc, target, method, errorMsg, preview, out)
	}

	responsePreviewLength = tempResponsePreviewLength
}

func logResult(sc int, target, method, errorMsg, preview string, out io.Writer) {
	symbol := "[!]"
	switch sc {
	case 200:
		symbol = "[+]"
	case 401, 403, 404:
		symbol = "[x]"
	}

	status := strconv.Itoa(sc)
	switch sc {
	case 0:
		status = "N/A"
	case 1:
		status = "SKIP"
	}

	fmt.Fprintf(out, "%s  %-6s %-4s %s\n", symbol, method, status, target)
	if preview != "" {
		fmt.Fprintf(out, "    %s\n", preview)
	}
	if errorMsg != "" && sc != 200 && sc != 1 {
		fmt.Fprintf(out, "    %s\n", errorMsg)
	}
}

func logJSON(title, description string, out io.Writer) {
	output := struct {
		APITitle    string   `json:"apiTitle"`
		Description string   `json:"description"`
		Results     []Result `json:"results"`
	}{
		APITitle:    title,
		Description: description,
		Results:     jsonResultArray,
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		die("Error marshalling JSON: %v", err)
	}
	fmt.Fprintln(out, string(data))
}

func logVerboseJSON(title, description string, out io.Writer) {
	output := struct {
		APITitle    string          `json:"apiTitle"`
		Description string          `json:"description"`
		Results     []VerboseResult `json:"results"`
	}{
		APITitle:    title,
		Description: description,
		Results:     jsonVerboseResultArray,
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		die("Error marshalling JSON: %v", err)
	}
	fmt.Fprintln(out, string(data))
}
