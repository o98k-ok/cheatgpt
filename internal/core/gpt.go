package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/o98k-ok/cheatgpt/internal/entity"
	np "github.com/o98k-ok/pcurl/http"
)

var (
	DefaultModel    string = "gpt-3.5-turbo"
	DefaultMaxToken int    = 2048
	DefaultHost     string = "https://api.openai.com/v1"
)

type GPT struct {
	ApiHost string
	Token   string
}

func NewGPT(apiKey string) *GPT {
	return &GPT{
		ApiHost: DefaultHost,
		Token:   apiKey,
	}
}

func NewRequest(msgs []entity.Message) *entity.GPTRequest {
	return &entity.GPTRequest{
		Model:     DefaultModel,
		Messages:  msgs,
		MaxTokens: DefaultMaxToken,
	}
}

func (g *GPT) Ask(req *entity.GPTRequest) (*entity.GPTResonse, error) {
	by := new(bytes.Buffer)
	json.NewEncoder(by).Encode(req)

	url := g.ApiHost + "/chat/completions"
	resp, err := np.NewRequest(entity.CLI, url).
		WithMethod(http.MethodPost).
		AddHeader("Authorization", fmt.Sprintf("Bearer %s", g.Token)).
		AddHeader("Content-Type", "application/json").WithData(by).Do()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ent entity.GPTResonse
	if err := json.NewDecoder(resp.Body).Decode(&ent); err != nil {
		return nil, err
	}
	return &ent, nil
}
