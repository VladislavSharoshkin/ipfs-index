package internal

import (
	"context"
	shell "github.com/ipfs/go-ipfs-api"
	"log"
	"strings"
)

type Internal struct {
	sh *shell.Shell
}

func NewInternal() *Internal {
	sh := shell.NewShell("127.0.0.1:5001")

	subscribe, err := sh.PubSubSubscribe("ipfs-index")
	if err != nil {
		return nil
	}

	go func() {
		for {
			next, err := subscribe.Next()
			if err != nil {
				return
			}
			query := string(next.Data)
			log.Println("New query:", query)

			ls, err := sh.FilesLs(context.Background(), "/", shell.FilesLs.Stat(true))
			if err != nil {
				return
			}

			for _, file := range ls {

				log.Println("File:", file.Name, file.Hash)
				if strings.Contains(file.Name, query) {
					err := sh.PubSubPublish(query, file.Name)
					if err != nil {
						return
					}
				}
			}
		}
	}()

	select {}
}
