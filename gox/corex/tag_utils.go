package corex

import (
	"strings"
	"unicode"
)

// TagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type TagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
// 把tag分割为名称和功能
func ParseTagToNameOption(tag string) (string, TagOptions) {
	if len(tag) == 0 {
		return "", TagOptions("")
	}
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, TagOptions(opt)
}

// 把tag分割为各个功能
func ParseTagToOptions(tag string) []TagOptions {
	if len(tag) == 0 {
		return nil
	}
	var subTags []TagOptions
	s := string(tag)
	for s != "" {
		var subOption string
		subOption, s, _ = strings.Cut(s, ",")
		subTags = append(subTags, TagOptions(subOption))
	}
	return subTags
}

// 把tag分割为小功能块的子subTag
func ParseToSubTag(tag string) []string {
	if len(tag) == 0 {
		return nil
	}
	var subTags []string
	s := string(tag)
	for s != "" {
		var subTag string
		subTag, s, _ = strings.Cut(s, ";")
		subTags = append(subTags, subTag)

	}
	return subTags
}

// 判断tag是否合法
func TagIsValid(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o TagOptions) Contains(optionKey string) bool {
	if len(o) == 0 || len(optionKey) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		//if name == optionKey {
		//	return true
		//}
		if name == optionKey || strings.HasPrefix(name, optionKey+"=") || strings.HasPrefix(name, optionKey+":") || strings.HasPrefix(name, optionKey+".") {
			return true
		}
	}
	return false
}

// 返回一个控制指令的具体值，可以以=:.三种符号分割控制质量和控制值
func (o TagOptions) OptionValue(optionKey string) string {
	oLen := len(o)
	keyLen := len(optionKey)
	if oLen <= 0 || keyLen <= 0 {
		return ""
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionKey {
			return ""
		} else if strings.HasPrefix(name, optionKey+"=") || strings.HasPrefix(name, optionKey+":") || strings.HasPrefix(name, optionKey+".") {
			return name[keyLen+1:]
		}
	}
	return ""
}
