package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (int, bool, error) {
	prestr := string(data)
	if !strings.Contains(prestr, CRLF) {
		return 0, false, nil
	}
	if strings.Index(prestr, CRLF) == 0 {
		return 0, true, nil
	}

	positionColon := strings.Index(prestr, ":")

	if prestr[positionColon-1] == ' ' {
		return 0, false, errors.New("Invalid white space between field-name and the colon")
	}

	prestr = strings.ReplaceAll(prestr, "\r\n", "")
	prestr = strings.Replace(prestr, " ", "", -1)

	str := strings.SplitN(prestr, ":", 2)

	h[str[0]] = str[1]

	return len(prestr) + 3, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
