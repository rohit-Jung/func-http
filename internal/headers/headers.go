package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const CRLF = "\r\n"

var (
	errDidNotFoundCRLF           = fmt.Errorf("ERORR: Did not found CRLF")
	errFieldLineKeyHasWhiteSpace = fmt.Errorf("ERROR: Field line key shouldn't contain whitespace")
	errMalformedFieldLine        = fmt.Errorf("ERROR: Got malformed Field Line")
)

func parseSingleFieldLine(fieldLine []byte) (string, string, error) {
	fieldLineParts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(fieldLineParts) != 2 {
		return "", "", errMalformedFieldLine
	}

	cleanedFieldName := string(bytes.TrimSpace(fieldLineParts[0]))

	// it has white space ?
	if cleanedFieldName != string(fieldLineParts[0]) {
		return "", "", errFieldLineKeyHasWhiteSpace
	}

	cleanedFieldVal := string(bytes.TrimSpace(fieldLineParts[1]))
	return cleanedFieldName, cleanedFieldVal, nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	bytesRead := 0
	doneParsing := false

	for {
		indexOfCrlf := bytes.Index(data[bytesRead:], []byte(CRLF))
		if indexOfCrlf == -1 {
			return bytesRead, doneParsing, errDidNotFoundCRLF
		}

		// it found it at start itself
		if indexOfCrlf == 0 {
			doneParsing = true
			break
		}

		fieldLine := data[bytesRead : bytesRead+indexOfCrlf]
		fieldName, fieldValue, err := parseSingleFieldLine(fieldLine)
		if err != nil {
			return bytesRead, doneParsing, err
		}

		h[fieldName] = fieldValue
		bytesRead += indexOfCrlf + len([]byte(CRLF))
	}

	return bytesRead, doneParsing, nil
}
