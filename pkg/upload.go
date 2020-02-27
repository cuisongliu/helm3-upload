package pkg

import (
	"bytes"
	"crypto/tls"
	"github.com/wonderivan/logger"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type HarborCfg struct {
	URL      string
	UserName string
	Password string
}

// GetUserAgent returns a user agent for user with an HTTP client
func userAgent() string {
	return "Helm/" + strings.TrimPrefix("v3.0", "v")
}

func (cfg *HarborCfg) UploadChartPackage(chartPackagePath string, force bool) ([]byte, error) {
	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}
	parsedURL.RawPath = path.Join("api", parsedURL.RawPath, "charts")
	parsedURL.Path = path.Join("api", parsedURL.Path, "charts")
	logger.Debug("fetch url:", parsedURL.Path)
	indexURL := parsedURL.String()

	req, err := http.NewRequest("POST", indexURL, nil)
	if err != nil {
		return nil, err
	}
	// Add ?force to request querystring to force an upload if chart version already exists
	if force {
		req.URL.RawQuery = "force"
	}

	err = setUploadChartPackageRequestBody(req, chartPackagePath)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent())
	req.SetBasicAuth(cfg.UserName, cfg.Password)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)
	resp.Body.Close()
	return ioutil.ReadAll(buf)
}

func setUploadChartPackageRequestBody(req *http.Request, chartPackagePath string) error {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	defer w.Close()
	fw, err := w.CreateFormFile("chart", chartPackagePath)
	if err != nil {
		return err
	}
	w.FormDataContentType()
	fd, err := os.Open(chartPackagePath)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Body = ioutil.NopCloser(&body)
	return nil
}
