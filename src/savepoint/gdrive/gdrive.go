package gdrive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func Savepoint(filePath string, identity map[string]string) error {
	ctx := context.Background()

	credentialFile := identity["credentialFile"]
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		return errors.New("failed to read credential file: " + err.Error())
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return errors.New("failed to read credential file: " + err.Error())
	}

	client, err := getClient(config)
	if err != nil {
		return errors.New("cannot create gdrive client: " + err.Error())
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return errors.New("unable to retrieve Drive client: " + err.Error())
	}

	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("unable to read tarball: " + err.Error())
	}

	fileInf, err := file.Stat()
	if err != nil {
		return errors.New("unable to retrieve tarball information: " + err.Error())
	}
	defer file.Close()

	folderId := identity["folderId"]
	f := &drive.File {
		Name: strings.Replace(filePath, "./tmp/", "", 1),
	}

	if folderId != "" {
		f.Parents = []string { folderId }
	}
	_, err = srv.Files.
		Create(f).
		// TODO: use Media instead.
		ResumableMedia(context.Background(), file, fileInf.Size(), "application/tar+gzip").
		Do()
	if err != nil {
		return errors.New("failed to send a backup to gdrive: " + err.Error())
	}

	return nil
}

func getClient(config *oauth2.Config) (*http.Client, error) {
	if _, err := os.Stat("./tokens/"); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir("./tokens/", os.ModePerm)
		if err != nil {
			return nil, errors.New("cannot create tokens directory")
		}
	}

	tokFile := "./tokens/gdrive-" + config.ClientID + ".json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}

		err = saveToken(tokFile, tok)
		if err != nil {
			return nil, err
		}
	}

	return config.Client(context.Background(), tok), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token {}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the " +
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.New("unable to read authorization code: " + err.Error())
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, errors.New("unable to retrieve token from web: " + err.Error())
	}

	return tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.New("unable to cache oauth token: " + err.Error())
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}