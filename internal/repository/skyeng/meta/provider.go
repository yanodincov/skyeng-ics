package meta

import httphelper "github.com/yanodincov/skyeng-ics/pkg/http-helper"

const (
	headerAcceptEncoding = "gzip, deflate, br, zstd"
	headerAcceptLanguage = "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7"
)

type SkyengMeta struct {
	GetCsfrHeaders       map[string]string
	LoginHeaders         map[string]string
	GetAuthCookieHeaders map[string]string
	GetScheduleHeaders   map[string]string
}

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) GenerateSkyengMeta() *SkyengMeta {
	userAgent := httphelper.GetUserAgent()
	secChUa := httphelper.GetSecCHUA(userAgent)
	secChUaPlatform := httphelper.GetSecCHUAPlatform(userAgent)

	csfrTokenHeader := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7", //nolint:lll
		"Accept-Encoding":           headerAcceptEncoding,
		"Accept-Language":           headerAcceptLanguage,
		"Cache-Control":             "no-cache",
		"Pragma":                    "no-cache",
		"Priority":                  "u=0, i",
		"Referer":                   "https://skyeng.ru/",
		"Sec-Ch-Ua-Mobile":          "?0",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "same-site",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                userAgent,
		"Sec-Ch-Ua":                 secChUa,
	}
	loginHeader := map[string]string{
		"Accept":             "*/*",
		"Accept-Encoding":    headerAcceptEncoding,
		"Accept-Language":    headerAcceptLanguage,
		"Content-Type":       "application/x-www-form-urlencoded; charset=UTF-8",
		"Origin":             "https://id.skyeng.ru",
		"Priority":           "u=1, i",
		"Referer":            "https://id.skyeng.ru/login?redirect=https%3A%2F%2Fskyeng.ru%2F",
		"Sec-Ch-Ua":          secChUa,
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": secChUaPlatform,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         userAgent,
		"X-Requested-With":   "XMLHttpRequest",
	}
	getAuthCookieHeader := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Encoding":    headerAcceptEncoding,
		"Accept-Language":    "en",
		"Origin":             "https://student.skyeng.ru",
		"Priority":           "u=1, i",
		"Referer":            "https://student.skyeng.ru/",
		"Sec-Ch-Ua":          secChUa,
		"User-Agent":         userAgent,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-site",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": secChUaPlatform,
	}
	sheduleHeader := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Encoding":    headerAcceptEncoding,
		"Accept-Language":    "en",
		"Content-Type":       "application/json",
		"Origin":             "https://student.skyeng.ru",
		"Priority":           "u=1, i",
		"Referer":            "https://student.skyeng.ru/",
		"Sec-Ch-Ua":          secChUa,
		"User-Agent":         userAgent,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-site",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": secChUaPlatform,
	}

	return &SkyengMeta{
		GetCsfrHeaders:       csfrTokenHeader,
		LoginHeaders:         loginHeader,
		GetAuthCookieHeaders: getAuthCookieHeader,
		GetScheduleHeaders:   sheduleHeader,
	}
}
