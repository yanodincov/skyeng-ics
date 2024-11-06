package httphelper

import (
	"net/http"
)

const cookieDeletedValue = "deleted"

func MergeCookies(httpRes *http.Response, cookies []*http.Cookie) []*http.Cookie {
	cookieByName := make(map[string]*http.Cookie)
	for _, cookie := range cookies {
		cookieByName[cookie.Name] = cookie
	}

	for _, cookie := range httpRes.Cookies() {
		if cookie.Value == cookieDeletedValue {
			delete(cookieByName, cookie.Name)
		} else {
			cookieByName[cookie.Name] = cookie
		}
	}

	res := make([]*http.Cookie, 0, len(cookieByName))
	for _, cookie := range cookieByName {
		res = append(res, cookie)
	}

	return res
}
