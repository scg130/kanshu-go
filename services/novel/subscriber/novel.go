package subscriber

import (
	"context"
	"fmt"
	"novel/repo"

	log "github.com/micro/go-micro/v2/logger"

	novel "novel/proto/novel"
)

type NovelRead struct {
	Note repo.Notes
}

func (e *NovelRead) Handle(ctx context.Context, msg *novel.ReadRequest) (err error) {
	log.Info("novel read handler received message: ", msg)

	note, err := e.Note.GetNote(msg.UserId, msg.NovelId, msg.ChapterNum)
	if err != nil {
		return err
	}
	if note.IsDelete == 1 {
		fmt.Println("RecoveryNote")
		err = e.Note.RecoveryNote(msg.UserId, msg.NovelId, msg.ChapterNum)
		if err != nil {
			return err
		}
		return
	}
	if note.Id > 0 {
		return
	}
	err = e.Note.CreateNote(msg.UserId, msg.NovelId, msg.ChapterNum, msg.IsJoin)
	if err != nil {
		return err
	}
	return nil
}

func Handler(ctx context.Context, msg *novel.Message) error {
	log.Info("Function Received message: ", msg.Flag)
	return nil
}
