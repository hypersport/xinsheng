package main

import (
	xinsheng "fuck996/cmd/src"
	"net/http"
)

func main() {
	http.HandleFunc("/", xinsheng.IndexHandler)
	xinsheng.Error.Println(http.ListenAndServe(":8888", nil))
}
