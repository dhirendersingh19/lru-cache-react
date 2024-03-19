package controllers

import "lru-cache/util"

var cache = util.NewLRUcache(1024)

func GET(id string) (interface{}, error, int) {
	return nil, nil, 200
}

func SET(id string) (interface{}, error, int) {
	return nil, nil, 200
}
