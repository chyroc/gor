package gor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchReg(t *testing.T) {
	as := assert.New(t)

	{
		var noRegMatch = func(matchType matchType) {
			params, matched := matchPath("/a", "/a", matchType)
			as.Equal(params, map[string]string{})
			as.True(matched)

			params, matched = matchPath("/a/", "/a/", matchType)
			as.Equal(params, map[string]string{})
			as.True(matched)

			params, matched = matchPath("/a", "/a/", matchType)
			as.Equal(params, map[string]string{})
			as.True(matched)

			params, matched = matchPath("/a/", "/a", matchType)
			as.Equal(params, map[string]string{})
			as.True(matched)
		}
		{
			// fullmatch + no-reg
			noRegMatch(fullMatch)

			// premathc + no-reg
			noRegMatch(preMatch)

			params, matched := matchPath("/", "/a", preMatch)
			as.Equal(params, map[string]string{})
			as.True(matched)

			params, matched = matchPath("/a/", "/a/b/c", preMatch)
			as.Equal(params, map[string]string{})
			as.True(matched)

			params, matched = matchPath("/a/b", "/a/b/c", preMatch)
			as.Equal(params, map[string]string{})
			as.True(matched)

			params, matched = matchPath("/a", "/", preMatch)
			as.Equal(params, map[string]string{})
			as.False(matched)
		}
	}
	{
		var realRegMatch = func(matchtype matchType) {
			params, matched := matchPath("/:aname", "/", matchtype)
			as.Equal(params, map[string]string{})
			as.False(matched)

			params, matched = matchPath("/:aname", "/a", matchtype)
			as.Equal(map[string]string{"aname": "a"}, params)
			as.True(matched)

			params, matched = matchPath("/:ANAME", "/a/", matchtype)
			as.Equal(map[string]string{"aname": "a"}, params)
			as.True(matched)

			params, matched = matchPath("/:user/:name", "/a/b", matchtype)
			as.Equal(map[string]string{"user": "a", "name": "b"}, params)
			as.True(matched)

			params, matched = matchPath("/:user/no-param/:name", "/a/no-param/b", matchtype)
			as.Equal(map[string]string{"user": "a", "name": "b"}, params)
			as.True(matched)

			params, matched = matchPath("/:user/no-param/:name", "/a/oh-no/b", matchtype)
			as.Equal(params, map[string]string{})
			as.False(matched)
		}
		{
			// fullmatch + reg
			realRegMatch(fullMatch)

			params, matched := matchPath("/:aname", "/a/b", fullMatch)
			as.Equal(params, map[string]string{})
			as.False(matched)
		}
		{
			// fprematch + reg
			realRegMatch(preMatch)

			params, matched := matchPath("/:a", "/a/b", preMatch)
			as.Equal(map[string]string{"a": "a"}, params)
			as.True(matched)

			params, matched = matchPath("/:a", "/a/b/c/d", preMatch)
			as.Equal(map[string]string{"a": "a"}, params)
			as.True(matched)

			params, matched = matchPath("/:a/:b", "/a/b/c/d", preMatch)
			as.Equal(map[string]string{"a": "a", "b": "b"}, params)
			as.True(matched)

			params, matched = matchPath("/:a/user", "/a/user/:b", preMatch)
			as.Equal(map[string]string{"a": "a"}, params)
			as.True(matched)
		}
	}
}
