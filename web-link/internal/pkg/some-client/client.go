package some_client

import (
	"flag"
	_ "github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
	"log"
)

type Client struct {
	// TODO: client
	UID string
	AccessToken string
}

func New() *Client {
	return &Client{}
}

func (c *Client)getLinkInfo() string  {
	// TODO: logic
	return ""
}

func (c *Client)Auth() string  {
	// TODO: logic
	return ""
}

func (c *Client)addLink(url string) string  {
	// TODO: logic
	return ""
}

func (c *Client)delLink(shorturl string) string  {
	// TODO: logic
	return ""
}

func (c *Client)showLinks() string  {
	// TODO: logic
	return ""
}

func (c *Client)openShortLink(url string) string  {
	// TODO: logic
	return ""
}

func main()  {
	log.Print("Starting the app")
	// настройка порта, настроек хранилища, таймаут при закрытии сервиса
	port := flag.String("port", "8000", "Port")
}