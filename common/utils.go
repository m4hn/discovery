package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	toolsRender "github.com/devopsext/tools/render"
	"github.com/devopsext/utils"
)

func ReadFiles(pattern string) ([]string, error) {

	ret := []string{}
	err := filepath.Walk(pattern, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ret = append(ret, path)
		}
		return nil
	})
	return ret, err
}

func MergeStringMaps(maps ...map[string]string) map[string]string {

	r := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			r[k] = v
		}
	}
	return r
}

func MergeInterfacegMaps(maps ...map[string]interface{}) map[string]interface{} {

	r := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			r[k] = v
		}
	}
	return r
}

func ByteMD5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func FileMD5(path string) []byte {

	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil
	}
	return h.Sum(nil)
}

func IfDef(v, def interface{}) interface{} {
	if v == nil {
		return def
	}
	switch v.(type) {
	case string:
		if v.(string) == "" {
			return def
		}
	case int:
		if v.(int) == 0 {
			return def
		}
	}
	return v
}

func StringInArr(a string, arr []string) bool {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}

func RenderTemplate(tpl *toolsRender.TextTemplate, def string, obj interface{}) (string, error) {

	if tpl == nil {
		return def, nil
	}

	b, err := tpl.RenderObject(obj)
	if err != nil {
		return def, err
	}
	r := strings.TrimSpace(string(b))
	// simplify <no value> => empty string
	return strings.ReplaceAll(r, "<no value>", ""), nil
}

func Render(def string, obj interface{}, observability *Observability) string {

	logger := observability.Logs()
	tpl, err := toolsRender.NewTextTemplate(toolsRender.TemplateOptions{Content: def}, observability)
	if err != nil {
		logger.Error(err)
		return def
	}

	s, err := RenderTemplate(tpl, def, obj)
	if err != nil {
		logger.Error(err)
		return def
	}
	return s
}

func GetStringKeys(arr map[string]string) []string {
	var keys []string
	for k := range arr {
		keys = append(keys, k)
	}
	return keys
}

func SortStringMapByKeys(m map[string]string, keys []string) map[string]string {

	r := make(map[string]string)
	for _, k := range keys {
		r[k] = m[k]
	}
	return r
}

func GetBaseConfigKeys(arr map[string]*BaseConfig) []string {
	var keys []string
	for k := range arr {
		keys = append(keys, k)
	}
	return keys
}

func GetFileKeys(arr map[string]*File) []string {
	var keys []string
	for k := range arr {
		keys = append(keys, k)
	}
	return keys
}

func GetLabelsKeys(arr map[string]Labels) []string {
	var keys []string
	for k := range arr {
		keys = append(keys, k)
	}
	return keys
}

func StringContainsAny(s string, arr []string) bool {

	for _, v := range arr {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

func ParsePeriodFromNow(period string, t time.Time) string {

	durStr := period
	if utils.IsEmpty(durStr) {
		return ""
	}

	if durStr == "" {
		durStr = "0s"
	}

	if durStr == "0d" {
		durStr = "0h"
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return ""
	}

	from := t.Add(time.Duration(dur))
	return strconv.Itoa(int(from.Unix()))
}

func GetPrometheusDiscoveriesByInstances(names string) map[string]string {

	m := make(map[string]string)
	def := "unknown"
	arr := strings.Split(names, ",")
	if len(arr) > 0 {
		index := 0
		for _, v := range arr {

			n := fmt.Sprintf("%s%d", def, index)
			kv := strings.Split(v, "=")
			if len(kv) > 1 {
				name := strings.TrimSpace(kv[0])
				if utils.IsEmpty(name) {
					name = n
				}
				url := strings.TrimSpace(kv[1])
				if !utils.IsEmpty(url) {
					m[name] = url
				}
			} else {
				m[n] = strings.TrimSpace(kv[0])
			}
			index++
		}
	} else {
		m[def] = strings.TrimSpace(names)
	}
	return m
}
