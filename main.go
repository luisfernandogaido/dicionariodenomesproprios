package main

import (
	"net/http"
	"io/ioutil"
	"github.com/pkg/errors"
	"strconv"
	"log"
	"regexp"
	"time"
	"strings"
	"fmt"
)

const site = "https://www.dicionariodenomesproprios.com.br"

var client = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func main() {
	//os.Setenv("HTTP_PROXY", "http://proxycorsp:8060")
	pnm, err := paginasNomesMasculinos()
	if err != nil {
		log.Fatal(err)
	}
	pnf, err := paginasNomesFemininos()
	if err != nil {
		log.Fatal(err)
	}
	nomes, err := getNomesMasculinos(pnm)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./nomes-masculinos.txt", []byte(strings.Join(nomes, "\r\n")), 0644)
	if err != nil {
		log.Fatal(err)
	}
	nomes, err = getNomesFemininos(pnf)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./nomes-femininos.txt", []byte(strings.Join(nomes, "\r\n")), 0644)
	if err != nil {
		log.Fatal(err)
	}
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

func getNomes(endpoint string, paginas int) ([]string, error) {
	re, err := regexp.Compile(`<a class="lista-nome" href="[^"]+">([^<]+)</a>`)
	if err != nil {
		return nil, err
	}
	nomes := make([]string, 0)
	for i := 1; i <= paginas; i++ {
		fmt.Printf("\n--- Página " + strconv.Itoa(i) + " ---\n\n")
		html, err := get(endpoint + strconv.Itoa(i) + "/")
		if err != nil {
			return nil, err
		}
		matches := re.FindAllStringSubmatch(html, -1)
		for _, m := range matches {
			nomes = append(nomes, m[1])
			fmt.Println(m[1])
		}
		time.Sleep(time.Millisecond * 125)
	}
	return nomes, nil
}

func getNomesMasculinos(paginas int) ([]string, error) {
	return getNomes("/nomes-masculinos/", paginas)
}

func getNomesFemininos(paginas int) ([]string, error) {
	return getNomes("/nomes-femininos/", paginas)
}
