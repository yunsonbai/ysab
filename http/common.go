package http

import (
	"github.com/yunsonbai/ysab/conf"
	"github.com/yunsonbai/ysab/summary"
)

var (
	config = conf.Conf
)

func Request(url, data, method string, headers map[string]string, fast int) summary.Res {
	if fast == 0 {
		return do(url, method, data, headers)
	}
	return fastDo(url, method, data, headers)
}
