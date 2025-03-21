package headers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (int, bool, error) {
	fmt.Println(string(data))
	prestr := string(data)
	if !strings.Contains(prestr, CRLF) {
		return 0, false, nil
	}
	if strings.Index(prestr, CRLF) == 0 {
		return 0, true, nil
	}

	prestr = strings.Split(prestr, "\r\n")[0]
	positionColon := strings.Index(prestr, ":")

	if positionColon == -1 {
		return 0, false, nil
	} else if prestr[positionColon-1] == ' ' {
		return 0, false, errors.New("Invalid white space between field-name and the colon")
	}

	prestr = strings.Replace(prestr, " ", "", -1)

	str := strings.SplitN(prestr, ":", 2)

	if b, err := regexp.MatchString("^[!#$%&'*+-/\\.\\^_`|~A-Za-z0-9.]*$", str[0]); !b || err != nil {
		if err != nil {
			return 0, false, err
		}

		return 0, false, errors.New("The fild name is not good, " + str[0])
	}

	_, exist := h[strings.ToLower(str[0])]

	if exist {
		h[strings.ToLower(str[0])] += ", " + str[1]
	} else {

		h[strings.ToLower(str[0])] = str[1]
	}

	return len(prestr) + 3, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
