package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type Client struct {
	Endpoint   *url.URL
	HTTPClient *http.Client
	Cookie     *http.Cookie
	Jar        *http.CookieJar
	Logger     *log.Logger
}

func NewClient(urlStr string, logger *log.Logger) (*Client, error) {
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url :%s", urlStr)
	}

	var discardLogger = log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &Client{
		Endpoint: parsedURL,
		HTTPClient: &http.Client{
			Jar: jar,
		},
		Cookie: nil,
		Logger: logger,
	}

	return client, nil
}

func (c *Client) newRequest(ctx context.Context, method string, urlpath string, body io.Reader) (*http.Request, error) {
	endpoint := *c.Endpoint
	endpoint.Path = path.Join(c.Endpoint.Path, urlpath)

	req, err := http.NewRequest(method, endpoint.String(), body)
	if err != nil {
		return nil, err
	}

	cookieStr, _ := loadCookieStr()
	cookies := strings.Join(cookieStr, ";")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookies)

	return req, nil
}

func (c *Client) GetCsrfToken(urlpath string) string {
	ctx := context.Background()
	req, _ := c.newRequest(ctx, "GET", urlpath, nil)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	csrfToken, found := doc.Find(`form input[type="hidden"]`).Attr("value")
	if !found {
		log.Fatal("error: cannot find CSRF token.")
	}
	return csrfToken
}

func (c *Client) IsLoggedIn() bool {
	ctx := context.Background()
	req, _ := c.newRequest(ctx, "GET", "/contests/abc001/submit", nil)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// If not logged in, redirected to login page.
	return string(res.Request.URL.Path) != "/login"
}

func (c *Client) Login(username, password string) {
	ctx := context.Background()
	csrfToken := c.GetCsrfToken("/login")

	values := url.Values{
		"username":   {username},
		"password":   {password},
		"csrf_token": {csrfToken},
	}

	req, _ := c.newRequest(ctx, "POST", "/login", strings.NewReader(values.Encode()))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()

	// TODO: すでにログイン済みの場合も"/login"に戻されてしまう。IsLoggedInを先に実行して回避可能。
	if string(res.Request.URL.Path) == "/login" {
		log.Fatal("Falied to login. Check your username/password.")
		return
	}

	fmt.Println("Succeffsully logged in!")

	writeLines(NewAcConfig().CookiePath, res.Cookies())
}

func (c *Client) Logout() {
	err := os.Remove(NewAcConfig().CookiePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logout")
}

func loadCookieStr() ([]string, error) {
	cookieStr, err := readLines(NewAcConfig().CookiePath)
	return cookieStr, err
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(filePath string, values []*http.Cookie) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range values {
		fmt.Fprintln(f, value) // write values to f, one per line.
	}
	return nil
}
