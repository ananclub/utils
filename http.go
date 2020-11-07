package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func HttpDo(remoteUrl string, method string, tr *http.Transport, j http.CookieJar, headers map[string]string, queryValues url.Values, data []byte) (statusCode int, body []byte, err error) {

	body = nil
	uri, err := url.Parse(remoteUrl)
	if err != nil {
		return
	}
	method = strings.ToUpper(method)
	if method == "POST" && data == nil {
		data = []byte(queryValues.Encode())
	} else {
		if queryValues != nil {
			values := uri.Query()
			if values != nil {
				for k, v := range values {
					queryValues[k] = v
				}
			}
			uri.RawQuery = queryValues.Encode()
		}
	}
	var client *http.Client
	if j != nil {
		client.Jar = j
	}
	if uri.Scheme == "https" {
		if tr == nil {
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}

		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}
	//println(uri.String())
	//println(string(data))
	request, err := http.NewRequest(method, uri.String(), bytes.NewReader(data))
	if err != nil {
		return
	}
	request.Close = true
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", uri.Host)
	request.Header.Add("Referer", uri.String())
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	statusCode = response.StatusCode
	defer response.Body.Close()

	var n int

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(response.Body)
		for {
			buf := make([]byte, 1024)
			n, err = reader.Read(buf)

			if err != nil && err != io.EOF {
				return
			}

			if n == 0 {
				break
			}
			body = append(body, buf...)
		}
	default:
		body, err = ioutil.ReadAll(response.Body)

	}

	if cls := response.Header.Get("Content-Length"); cls != "" {
		cl, errP := strconv.ParseInt(cls, 10, 32)
		if errP == nil && cl >= 0 {
			if err == io.EOF && int(cl) == n {
				err = nil
			}
		}
	} else {
		if err == io.EOF {
			err = nil
		}
	}
	return
}
