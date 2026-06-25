package xui

import (
	"golang.design/x/clipboard"
)

import (
	"context"
	"log/slog"
)

var clipboardAvailable = false

type ClipboardFormat = clipboard.Format

const (
	ClipboardText  = clipboard.FmtText
	ClipboardImage = clipboard.FmtImage
)

func ReadClipboard(form ClipboardFormat) []byte {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Read(form)
}

func WriteClipboard(form ClipboardFormat, data []byte) <-chan struct{} {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Write(form, data)
}

func WatchClipboard(ctx context.Context, form ClipboardFormat) <-chan []byte {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Watch(ctx, form)
}

func init() {
	err := clipboard.Init()
	if err != nil {
		slog.Error("clipboard not available", "err", err)
	}
	clipboardAvailable = err == nil
}
