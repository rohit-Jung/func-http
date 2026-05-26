package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (string, bool) {
	lowerCasedKey := strings.ToLower(key)
	if value, ok := h[lowerCasedKey]; ok {
		return value, ok
	}

	return "", false
}

func (h Headers) Set(fieldName string, fieldValue string) {
	name := strings.ToLower(fieldName)
	if val, ok := h[name]; ok {
		h[name] = val + "," + fieldValue
	} else {
		h[name] = fieldValue
	}
}

func (h Headers) Replace(fieldName string, fieldValue string) {
	name := strings.ToLower(fieldName)
	h[name] = fieldValue
}

func (h Headers) Delete(fieldName string) {
	name := strings.ToLower(fieldName)
	delete(h, name)
}

func (h Headers) GetIntVal(key string, defaultVal int) int {
	val, exists := h.Get(key)
	if !exists {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return intVal
}

const CRLF = "\r\n"

var (
	errFieldLineKeyHasWhiteSpace = fmt.Errorf("ERROR: Field line key shouldn't contain whitespace")
	errMalformedFieldLine        = fmt.Errorf("ERROR: Got malformed Field Line")
	errInvalidCharactersFound    = fmt.Errorf("ERROR: Invalid Characters found in field name")
)

const validCharactersPattern = "^[a-zA-Z0-9!#$%&'*+\\-.\\^_`|~]*$"

func parseSingleFieldLine(fieldLine []byte) (string, string, error) {
	fieldLineParts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(fieldLineParts) != 2 {
		return "", "", errMalformedFieldLine
	}

	cleanedFieldName := string(bytes.TrimSpace(fieldLineParts[0]))
	isValidFieldName, err := regexp.Match(validCharactersPattern, []byte(cleanedFieldName))
	if err != nil {
		return "", "", err
	}

	if !isValidFieldName {
		return "", "", errInvalidCharactersFound
	}

	// it has white space ?
	if cleanedFieldName != string(fieldLineParts[0]) {
		return "", "", errFieldLineKeyHasWhiteSpace
	}

	cleanedFieldVal := string(bytes.TrimSpace(fieldLineParts[1]))
	return strings.ToLower(cleanedFieldName), cleanedFieldVal, nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	bytesRead := 0
	doneParsing := false

	for {
		indexOfCrlf := bytes.Index(data[bytesRead:], []byte(CRLF))
		if indexOfCrlf == -1 {
			break
		}

		// it found it at start itself
		if indexOfCrlf == 0 {
			bytesRead += len(CRLF)
			doneParsing = true
			break
		}

		fieldLine := data[bytesRead : bytesRead+indexOfCrlf]
		fieldName, fieldValue, err := parseSingleFieldLine(fieldLine)
		if err != nil {
			return bytesRead, doneParsing, err
		}

		h.Set(fieldName, fieldValue)
		bytesRead += indexOfCrlf + len([]byte(CRLF))
	}

	return bytesRead, doneParsing, nil
}
