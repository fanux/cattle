package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/docker/swarm/common"
)

func showFlags(c *cli.Context) {
	fmt.Println("cattle connect to: ", c.String("H"))
	fmt.Println("ENVS: ")
	for _, env := range c.StringSlice("e") {
		fmt.Println(env)
	}
	fmt.Println("labels: ")
	for _, lalel := range c.StringSlice("l") {
		fmt.Println(lalel)
	}
	fmt.Println("filters: ")
	for _, filter := range c.StringSlice("f") {
		fmt.Println(filter)
	}
	fmt.Println("numbers: ", c.Int("n"))
}

func buildRequestBody(c *cli.Context) (common.ScaleAPI, error) {
	body := common.ScaleAPI{}
	item := common.ScaleItem{}
	item.Labels = make(map[string]string)

	item.ENVs = c.StringSlice("e")
	for _, v := range c.StringSlice("l") {
		vSlice := strings.SplitN(v, "=", 2)
		if len(vSlice) == 2 {
			item.Labels[vSlice[0]] = vSlice[1]
		} else {
			log.Printf("invalid label: %s", v)
			return body, errors.New("invalid label")
		}
	}
	item.Filters = c.StringSlice("f")
	item.Number = c.Int("n")

	body.Items = append(body.Items, item)

	return body, nil
}

//SendRequest is
func SendRequest(body common.ScaleAPI, url string) error {
	return sendRequest(body, url)
}

func sendRequest(body common.ScaleAPI, url string) error {
	s, err := json.Marshal(body)
	if err != nil {
		log.Printf("encode json error: %s", err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url+"/scale", strings.NewReader(string(s)))
	if err != nil {
		log.Printf("new http request error: %s", err)
		return err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("do http request failed: %s", err)
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		log.Printf("send scale request failed: %s", err)
	} else {
		io.Copy(os.Stdout, resp.Body)
	}

	return err
}

func scale(c *cli.Context) {
	showFlags(c)
	body, err := buildRequestBody(c)
	if err != nil {
		return
	}
	sendRequest(body, c.String("H"))
}
