package photo

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func DownloadFile(URL, name string) (string, error) {
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
	}
	s := response.Request.URL.Path
	split := strings.Split(s, ".")
	ext := split[len(split)-1]
	fileName := fmt.Sprintf("%s.%s", name, ext)
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(file, response.Body)
	err = file.Close()
	if err != nil {
		return "", err
	}
	return fileName, nil
}

