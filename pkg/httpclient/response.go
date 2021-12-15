package httpclient

import (
	"io/ioutil"
	"net/http"
)

func (c *Calling) judgeResponse() {
	if c.response == nil {
		panic(errorf("the calling has't been done"))
	}
}

func (c *Calling) GetResponse() *http.Response {
	c.judgeResponse()
	return c.response
}

func (c *Calling) GetRespStatusCode() int {
	c.judgeResponse()
	return c.GetResponse().StatusCode
}

func (c *Calling) GetRespStatusIsOK() bool {
	c.judgeResponse()
	return c.GetRespStatusCode() == http.StatusOK
}

func (c *Calling) GetRespHeader() http.Header {
	c.judgeResponse()
	return c.GetResponse().Header
}

func (c *Calling) GetRespBody() ([]byte, error) {
	c.judgeResponse()
	if c.respBody == nil {
		defer c.GetResponse().Body.Close()
		bs, err := ioutil.ReadAll(c.GetResponse().Body)
		if err != nil {
			return nil, errorf(err)
		}
		c.respBody = bs
	}
	return c.respBody, nil
}

func (c *Calling) GetRespBodyString() (string, error) {
	bs, err := c.GetRespBody()
	if err != nil {
		return "", errorf(err)
	}
	return string(bs), nil
}
