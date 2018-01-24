package main

import (
	"net/http"
	"os"
	"io/ioutil"
	"github.com/pkg/errors"
	"strconv"
	"fmt"
	"log"
	"regexp"
	"time"
)

const site = "https://www.dicionariodenomesproprios.com.br"

var client = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func main() {
	os.Setenv("HTTP_PROXY", "http://proxycorsp:8060")
	pnm := 100
	err := getNomesMasculinos(pnm)
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

func getNomesMasculinos(paginas int) error {
	re, err := regexp.Compile(`<a class="lista-nome" href="[^"]+">([^<]+)</a>`)
	if err != nil {
		return err
	}
	nomes := make([]string, 0)
	for i := 1; i <= paginas; i++ {
		fmt.Printf("\n--- página %v ---\n\n", i)
		html, err := get("/nomes-masculinos/" + strconv.Itoa(i) + "/")
		if err != nil {
			return err
		}
		matches := re.FindAllStringSubmatch(html, -1)
		for _, m := range matches {
			fmt.Println(m[1])
			nomes = append(nomes, m[1])
		}
		fmt.Printf("%v elementos encontrados.\n", len(nomes))
		time.Sleep(time.Second * 1)
	}
	return nil
}
