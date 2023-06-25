package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/o98k-ok/cheatgpt/internal/core"
	"github.com/o98k-ok/cheatgpt/internal/entity"
	"github.com/o98k-ok/lazy/v2/alfred"
)

func Entry() {
	cli := alfred.NewApp("cheat gpt")
	cli.Bind("cheat", func(s []string) {
		envs, err := alfred.FlowVariables()
		if err != nil {
			alfred.ErrItems("cannot get some configs", err).Show()
			return
		}

		if _, ok := envs[entity.API_KEY]; !ok {
			alfred.ErrItems("missing API_KEY", errors.New("cannot get API_KEY config")).Show()
			return
		}

		gpt := core.NewGPT(envs[entity.API_KEY])
		if v, ok := envs[entity.API_HOST]; ok {
			gpt.ApiHost = v
		}

		content := strings.Join(s, " ")
		start := time.Now().UnixMilli()
		req := core.NewRequest([]entity.Message{
			{
				Role:    "assistant",
				Content: content,
			},
		})

		if v, ok := envs[entity.MAX_TOKEN]; ok {
			tk, err := strconv.Atoi(v)
			if err == nil {
				req.MaxTokens = tk
			}
		}
		result, err := gpt.Ask(req)

		cost := time.Now().UnixMilli() - start
		if err != nil {
			alfred.ErrItems("ask cheat gpt error", err).Show()
			return
		}

		finialResult := fmt.Sprintf("# %s:\nQ:%s?\n%s\n\n", time.Now().Format("2006-01-02 15:01:01"), content, result.Choices[0].Message.Content)
		items := alfred.NewItems()
		items.Append(alfred.NewItem("问题描述", content, finialResult))
		items.Append(alfred.NewItem("AI答复", result.Choices[0].Message.Content, finialResult))
		items.Append(alfred.NewItem("生成耗时", fmt.Sprintf("耗时%fs", float32(cost)/1000.0), finialResult))
		items.Show()
	})

	cli.Run(os.Args)
}
