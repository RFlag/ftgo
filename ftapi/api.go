package ftapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type resultI interface {
	ok() error
}

type ResultCodeError struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
}

func (r *ResultCodeError) ok() error {
	if r.Code != 1 {
		return errors.New(r.Error)
	}
	return nil
}

type ResultList struct {
	ResultCodeError
	Pend string `json:"pend,omitempty"`
}

var private string

func init() {
	log := stdlog.New(os.Stdout, "[api] ", stdlog.LstdFlags)

	private = os.Getenv("FTGO_PRIVATE")
	log.Print("FTGO_PRIVATE: " + private)
}

func Post(api string, param interface{}, result resultI) error {
	p, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err, "序列化参数错误")
	}
	resp, err := http.Post(private+api, "application/json", bytes.NewReader(p))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(r, result)
	if err != nil {
		return err
	}
	if err := result.ok(); err != nil {
		return err
	}
	return nil
}
