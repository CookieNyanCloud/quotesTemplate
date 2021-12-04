package photo

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(URL, name string) (*os.File, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}
	fileName := fmt.Sprintf("%s.jpg", name)
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, response.Body)
	return file, nil
}
