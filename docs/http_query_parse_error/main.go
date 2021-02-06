package main

import (
	"log"
	"net/http"
	"strconv"
)

func main () {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request){
		query := request.URL.Query()

		page, err := strconv.Atoi(query.Get("page")) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}

		log.Print("page", page)

		writer.WriteHeader(200)
		_, err = writer.Write([]byte("ok")) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	err := http.ListenAndServe(":3000", mux) ; if err != nil {
		// 服务无法启动必须 panic
		panic(err)
	}
}
