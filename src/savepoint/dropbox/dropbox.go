package dropbox

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Savepoint(filePath string, identity map[string]string) error {
	accessToken := identity["accessToken"]
	folderPath := identity["folderPath"]

	if !strings.HasPrefix(folderPath, "/") {
		folderPath = "/" + folderPath
	}

	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	uploadPath := folderPath + strings.Replace(filePath, "./tmp/", "", 1)

	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("unable to read tarball: " + err.Error())
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", file)
	if err != nil {
		return errors.New("cannot create HTTP request: " + err.Error())
	}

	req.Header.Add("Authorization", "Bearer " + accessToken)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Dropbox-API-Arg", "{\"path\":\"" + uploadPath + "\",\"mode\":\"overwrite\",\"autorename\":true}")

	client := &http.Client {}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("cannot send HTTP request: " + err.Error())
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)
	if strings.Contains(str, "\"error\": {") {
		errorSummary := strings.Split(str, "{\"error_summary\": \"")[1]
		errorSummary = strings.Split(errorSummary, "\", \"error\": {")[0]
		return errors.New("failed to send a backup to dropbox: " + errorSummary)
	}

	return nil
}