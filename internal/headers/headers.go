package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Get(field string) string {
	field = strings.ToLower(field)

	return h[field]
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	prestr := string(data)
	if !strings.Contains(prestr, CRLF) {
		return 0, false, nil
	}
	if strings.Index(prestr, CRLF) == 0 {
		return 2, true, nil
	}

	prestr = strings.Split(prestr, "\r\n")[0]
	n := len(prestr) + 2

	prestr = strings.Trim(prestr, " ")

	positionColon := strings.Index(prestr, ":")

	if positionColon == -1 {
		return 0, false, nil
	} else if prestr[positionColon-1] == ' ' {
		return 0, false, errors.New("Invalid white space between field-name and the colon")
	}

	str := strings.SplitN(prestr, ":", 2)

	if b, err := regexp.MatchString("^[!#$%&'*+-/\\.\\^_`|~A-Za-z0-9.]*$", str[0]); !b || err != nil {
		if err != nil {
			return 0, false, err
		}

		return 0, false, errors.New("The fild name is not good, " + str[0])
	}

	str[0] = strings.ToLower(str[0])
	str[1] = strings.ReplaceAll(str[1], " ", "")

	_, exist := h[str[0]]

	if exist {
		h[str[0]] += ", " + str[1]
	} else {

		h[str[0]] = str[1]
	}

	return n, false, nil
}

func (h Headers) Add(field, value string) {
	h[field] = value
}

func NewHeaders() Headers {
	return make(Headers)
}
