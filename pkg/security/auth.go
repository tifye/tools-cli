package security

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Tifufu/tools-cli/pkg"
	"github.com/charmbracelet/log"
)

type tadAuthRequest struct {
	AppId string `json:"app_id"`
	Url   string `json:"url"`
}

const tadAuthUrl = "https://tad.azurewebsites.net"
const authServerAddr = "127.0.0.1:3001"
const encryptionKey string = "{7f8d534a-bf20-4e69-bbf8-54f4a9378f23}"

func AuthenticateUser(ctx context.Context, appId string) (*UserProfile, error) {
	userChan := make(chan string)

	server := &http.Server{Addr: authServerAddr}
	http.HandleFunc("GET /auth", func(wri http.ResponseWriter, req *http.Request) {
		userChan <- req.URL.Query().Get("user")
	})
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			close(userChan)
			log.Fatal("ListenAndServer()", "err", err)
		}
	}()

	loginQuery := tadAuthRequest{
		AppId: appId,
		Url:   "http://" + authServerAddr + "/auth",
	}
	loginQueryJson, err := json.Marshal(loginQuery)
	if err != nil {
		return nil, err
	}
	url := tadAuthUrl + "?state=" + base64.StdEncoding.EncodeToString(loginQueryJson)
	pkg.OpenURL(url)

	var userId string
	select {
	case userId = <-userChan:
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout while waiting for authentication callback %w", ctx.Err())
	}

	if err := server.Shutdown(ctx); err != nil {
		return nil, err
	}

	userProfileChan := make(chan *UserProfile)
	errChan := make(chan error)
	go resolveUserId(userId, userProfileChan, errChan)
	select {
	case userProfile := <-userProfileChan:
		return userProfile, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout while resolving user profile %w", ctx.Err())
	}
}

func resolveUserId(userId string, userProfileChan chan<- *UserProfile, errChan chan<- error) {
	res, err := http.Get(tadAuthUrl + "/resolve?id=" + userId)
	if err != nil {
		errChan <- err
		return
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusConflict:
		errChan <- fmt.Errorf("the authentication key was already consumed %d", res.StatusCode)
		return
	default:
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			errChan <- err
			return
		}
		errChan <- fmt.Errorf("unexpected status code during auth: %d, %s", res.StatusCode, string(bytes))
		return
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		errChan <- err
	}

	var userProfile UserProfile
	if err := json.Unmarshal(bytes, &userProfile); err != nil {
		errChan <- err
	}
	userProfileChan <- &userProfile
}

func DecryptUserProfile(data []byte) (*UserProfile, error) {
	decrypted, err := Decrypt(data, []byte(encryptionKey))
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	var userProfile UserProfile
	if err := json.Unmarshal(decrypted, &userProfile); err != nil {
		return nil, fmt.Errorf("error unmarshalling decrypted data: %w", err)
	}
	return &userProfile, nil
}

func EncryptUserProfile(user *UserProfile) ([]byte, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return Encrypt(userJson, []byte(encryptionKey))
}
