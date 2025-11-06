package helpers

import (
	"context"
	"encoding/json"
	"github.com/febry3/gamingin/internal/dto"
	"golang.org/x/oauth2"
	"io"
)

func GetGoogleUserInfo(ctx context.Context, token *oauth2.Token, gauth *oauth2.Config) (dto.LoginWithGoogleData, error) {
	client := gauth.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return dto.LoginWithGoogleData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.LoginWithGoogleData{}, err
	}

	var userInfo dto.LoginWithGoogleData
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return dto.LoginWithGoogleData{}, err
	}
	return userInfo, nil
}
