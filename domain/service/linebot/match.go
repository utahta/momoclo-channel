package linebot

import (
	"regexp"
)

var (
	reMatchOn      = regexp.MustCompile("^(おん|オン|on)$")
	reMatchOff     = regexp.MustCompile("^(おふ|オフ|off)$")
	reMatchMomota  = regexp.MustCompile("百田|[もモ][もモ][たタ]|[夏かカ][菜なナ][子こコ]")
	reMatchAriyasu = regexp.MustCompile("有安|[あア][りリ][やヤ][すス]|[もモ][もモ][かカ]|杏果")
	reMatchTamai   = regexp.MustCompile("玉井|[たタ][まマ][いイ]|[しシ][おオ][りリ][んン]?|詩織|玉さん|[たタ][まマ]さん")
	reMatchSasaki  = regexp.MustCompile("佐々木|[さサ][さサ][きキ]|[あア][やヤ][かカ]|彩夏|[あア]ー[りリ][んン]")
	reMatchTakagi  = regexp.MustCompile("高城|[たタ][かカ][ぎギ]|[れレ][にニ]")
)

// MatchOn return true if text match on
func MatchOn(text string) bool {
	return reMatchOn.MatchString(text)
}

// MatchOff return true if text match off
func MatchOff(text string) bool {
	return reMatchOff.MatchString(text)
}

// FindMemberName returns member name if text match member name or nickname
func FindMemberName(text string) string {
	if reMatchMomota.MatchString(text) {
		return "百田夏菜子"
	}

	if reMatchAriyasu.MatchString(text) {
		return "有安杏果"
	}

	if reMatchTamai.MatchString(text) {
		return "玉井詩織"
	}

	if reMatchSasaki.MatchString(text) {
		return "佐々木彩夏"
	}

	if reMatchTakagi.MatchString(text) {
		return "高城れに"
	}
	return ""
}
