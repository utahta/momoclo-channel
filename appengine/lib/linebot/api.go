package linebot

import (
	"golang.org/x/net/context"
)

func ReplyText(ctx context.Context, replyToken, text string) error {
	client, err := New(ctx)
	if err != nil {
		return err
	}

	if err := client.ReplyText(replyToken, text); err != nil {
		return err
	}
	return nil
}

func ReplyImage(ctx context.Context, replyToken, originalContentURL, previewImageURL string) error {
	client, err := New(ctx)
	if err != nil {
		return err
	}

	if err := client.ReplyImage(replyToken, originalContentURL, previewImageURL); err != nil {
		return err
	}
	return nil
}
