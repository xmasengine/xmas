package xgal

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

// Paste gets clipped data from the clipboard.
func Paste(form ClipboardFormat) []byte {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Read(form)
}

// Board pastes to the clipboard if available.
func Board(form ClipboardFormat, data []byte) <-chan struct{} {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Write(form, data)
}

// Boarded returns a channel to watch any changed to the clipboard.
func Boarded(ctx context.Context, form ClipboardFormat) <-chan []byte {
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
