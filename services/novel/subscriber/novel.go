package subscriber

import (
	"context"
	"fmt"
	"github.com/ilylx/gconv"
	"novel/repo"

	log "github.com/micro/go-micro/v2/logger"

	novel "novel/proto/novel"
)

type NovelRead struct {
	Note      repo.Notes
	Chapter   repo.Chapter
	WalletLog repo.WalletLog
}

func (e *NovelRead) Handle(ctx context.Context, msg *novel.ReadRequest) (err error) {
	log.Info("novel read handler received message: ", msg)
	chapter, err := e.Chapter.GetOne(gconv.Int(msg.NovelId), gconv.Int(msg.ChapterNum))
	if err != nil {
		return err
	}
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
	wLog, _ := e.WalletLog.GetChapterByUserIdAndChapterId(gconv.Int(msg.UserId), gconv.Int(chapter.Id))
	if chapter.IsVip == int(novel.VipType_IS_VIP) && wLog.Id == 0 {
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
