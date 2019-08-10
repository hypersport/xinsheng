package main

import (
	"net/http"
	xinsheng "xinsheng/cmd/src"
)

func main() {
	http.HandleFunc("/", xinsheng.IndexHandler)
	xinsheng.Error.Println(http.ListenAndServe(":8888", nil))
}
