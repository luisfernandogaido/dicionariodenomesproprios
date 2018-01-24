package main

import (
	"net/http"
	"os"
	"io/ioutil"
	"github.com/pkg/errors"
	"fmt"
	"log"
	"strconv"
)

const site = "https://www.dicionariodenomesproprios.com.br"

var client = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func main() {
	os.Setenv("HTTP_PROXY", "http://proxycorsp:8060")
	p, err := paginasNomesMasculinos()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
	p, err = paginasNomesFemininos()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
}

func get(endpoint string) (string, error) {
	res, err := client.Get(site + endpoint)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New(res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func paginasNomes(endpoint string) (int, error) {
	p, q, m := 0, 1, 0
	var err error
	for err == nil {
		q *= 2
		_, err = get(endpoint + strconv.Itoa(q) + "/")
		if err != nil && err.Error() != "301 Moved Permanently" {
			return 0, errors.New(endpoint + ": " + err.Error())
		}
	}
	p = q / 2
	for {
		m = (p + q) / 2
		if m == p || m == q {
			return m, nil
		}
		_, err = get(endpoint + strconv.Itoa(m) + "/")
		if err == nil {
			p = m
		} else {
			if err.Error() != "301 Moved Permanently" {
				return 0, errors.New(endpoint + ": " + err.Error())
			}
			q = m
		}
	}
	return 0, nil
}

func paginasNomesMasculinos() (int, error) {
	return paginasNomes("/nomes-masculinos/")
}
func paginasNomesFemininos() (int, error) {
	return paginasNomes("/nomes-femininos/")
}
