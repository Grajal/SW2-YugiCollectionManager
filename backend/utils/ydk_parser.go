package utils

import (
	"bufio"
	"io"
	"mime/multipart"
	"strings"
)

func ParseYDK(file multipart.File) ([]string, []string, []string, error) {
	var (
		mainIDs  []string
		extraIDs []string
		sideIDs  []string
		current  *[]string
	)

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, nil, nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#created") {
			if err == io.EOF {
				break
			}
			continue
		}

		switch line {
		case "#main":
			current = &mainIDs
		case "#extra":
			current = &extraIDs
		case "#side":
			current = &sideIDs
		default:
			if current == nil {
				// Alternativa: log.Printf("Ignoring line outside of section: %s", line)
				continue
			}
			*current = append(*current, line)
		}

		if err == io.EOF {
			break
		}
	}

	return mainIDs, extraIDs, sideIDs, nil
}
