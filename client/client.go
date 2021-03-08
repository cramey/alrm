package client

import (
	"git.binarythought.com/cdramey/alrm/api"
	"git.binarythought.com/cdramey/alrm/config"
	"net/http"
	"net/url"
)

func Shutdown(cfg *config.Config) error {
	return doCommand(cfg, "shutdown")
}

func Restart(cfg *config.Config) error {
	return doCommand(cfg, "restart")
}

func doCommand(cfg *config.Config, cm string) error {
	aurl, err := url.Parse("http://" + cfg.Listen + "/api")
	if err != nil {
		return err
	}

	cmd := api.NewCommand(cm)
	err = cmd.Sign(cfg.APIKey)
	if err != nil {
		return err
	}
	cjson, err := cmd.JSON()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("cmd", string(cjson))

	resp, err := http.PostForm(aurl.String(), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
