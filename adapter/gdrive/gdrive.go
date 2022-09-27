package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	PFOLDER = "1pbyWUBtfsbLTy-_zQHklX81El8IduV1l"
	UFOLDER = "1HOaIg3yPnwGxLu1fbDAwBeP6dHxd-ACz"
)

type service struct {
	ds *drive.Service
}

func newDevService() service {
	ctx := context.Background()
	b, err := os.ReadFile("cmd/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	clt := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(clt))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return service{ds: srv}
}

func NewService() service {
	_, ok := os.LookupEnv("IS_DEV")
	if ok {
		return newDevService()
	}
	ctx := context.Background()
	b, ok := os.LookupEnv("DRIVE_SECRET")
	if !ok {
		log.Fatalf("Unable get client secret to config: DRIVE_SECRET not found")
	}
	config, err := google.ConfigFromJSON([]byte(b), drive.DriveFileScope)

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	clt := getClient(config)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(clt))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return service{ds: srv}
}

func (s service) SaveImg(h *multipart.FileHeader, name, p string) (string, error) {
	img, err := h.Open()
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	f := drive.File{Parents: []string{p}, Name: name}
	r, err := s.ds.Files.Create(&f).Media(img).IncludePermissionsForView("published").Do()
	if err != nil {
		return "", fmt.Errorf("create Media: %w", err)
	}
	return r.Id, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "./cmd/token.json"
	_, ok := os.LookupEnv("IS_DEV")
	if !ok {
		tokFile = ""
	}
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		tok := &oauth2.Token{}
		err = json.NewDecoder(f).Decode(tok)
		return tok, err
	}
	f, ok := os.LookupEnv("DRIVE_TOKEN")
	if !ok {
		return nil, fmt.Errorf("DRIVE_TOKEN not found")
	}
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(f), tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
