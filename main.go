package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Livro struct {
	Id     int    `json:"id"`
	Titulo string `json:"titulo"`
	Autor  string `json:"autor"`
}

var Livros []Livro = []Livro{
	{
		Id:     0,
		Titulo: "livro 01",
		Autor:  "Autor 01",
	},
	{
		Id:     1,
		Titulo: "livro 02",
		Autor:  "Autor 02",
	}, {
		Id:     2,
		Titulo: "livro 03",
		Autor:  "Autor 03",
	},
}

func configurarHandler() {
	http.HandleFunc("/", rotaPrincipal())
	http.HandleFunc("/livros", middlewareLivros)

	//GET /livros/{ID}
	http.HandleFunc("/livros/", middlewareLivros)

}

func buscarLivros(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	fmt.Println(path)
	// "/livros/{id}" -> ["","livros","{id}"]
	partes := strings.Split(path, "/")

	idLivro, err := strconv.Atoi(partes[2])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.Encode(strings.Contains("Erro ao Deserializar path", err.Error()))
	}

	for _, livro := range Livros {
		if livro.Id == idLivro {
			encode := json.NewEncoder(writer)
			encode.Encode(livro)
			return
		}
	}
	writer.WriteHeader(http.StatusNotFound)

}

//
func middlewareLivros(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	path := request.URL.Path
	partes := strings.Split(path, "/")
	if len(partes) == 2 || len(partes) == 3 && partes[2] == "" {
		if request.Method == "GET" {
			listarLivros(writer, request)
		} else if request.Method == "POST" {
			cadastrarLivro(writer, request)
		}
	} else if len(partes) == 3 {
		if request.Method == "GET" {
			buscarLivros(writer, request)
		} else if request.Method == "DELETE" {
			excluirLivro(writer, request)
		} else if request.Method == "PUT" {
			alterarLiro(writer, request)
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
}

func alterarLiro(writer http.ResponseWriter, request *http.Request) {

	path := request.URL.Path
	fmt.Println(path)
	// "/livros/{id}" -> ["","livros","{id}"]
	partes := strings.Split(path, "/")

	idLivro, err := strconv.Atoi(partes[2])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.Encode(strings.Contains("Erro ao Deserializar path", err.Error()))
	}

	indiceLivro := -1

	for indece, livro := range Livros {
		if livro.Id == idLivro {
			indiceLivro = indece
			break
		}
	}

	if indiceLivro < 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	bytesBody, err := ioutil.ReadAll(request.Body)
	if err != nil {

		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.Encode(strings.Contains("Erro ao Deserializar Json", err.Error()))
	}
	var livroModificado Livro
	err = json.Unmarshal(bytesBody, &livroModificado)
	if err != nil {

		writer.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(writer)
		encoder.Encode(strings.Contains("Erro ao Serializar Json", err.Error()))
	}

	Livros[indiceLivro] = livroModificado
	json.NewEncoder(writer).Encode(livroModificado)

}

func excluirLivro(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	fmt.Println(path)
	// "/livros/{id}" -> ["","livros","{id}"]
	partes := strings.Split(path, "/")

	idLivro, err := strconv.Atoi(partes[2])

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.Encode(strings.Contains("Erro ao Deserializar path", err.Error()))
	}
	indiceLivro := -1

	for indece, livro := range Livros {
		if livro.Id == idLivro {
			indiceLivro = indece
			break
		}
	}

	if indiceLivro < 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// a := [ a, b, c, d, e] remove c --> indece:=3  --> firstPart := a[0:indece] --> lastPart := a[indece+1: len(a)] --> a = append(firstPart, lastPart...)
	firstPart := Livros[0:indiceLivro]
	lastPart := Livros[indiceLivro+1 : len(Livros)]
	Livros = append(firstPart, lastPart...)
	writer.WriteHeader(http.StatusNoContent)
}

func cadastrarLivro(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusCreated)

	bytesBody, err := ioutil.ReadAll(request.Body)
	if err != nil {

		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.Encode(strings.Contains("Erro ao Deserializar Json", err.Error()))
	}
	var novoLivro Livro
	json.Unmarshal(bytesBody, &novoLivro)

	Livros = append(Livros, novoLivro)
	novoLivro.Id = len(Livros) - 1
	encoder := json.NewEncoder(writer)
	encoder.Encode(Livros)

}

func listarLivros(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	encoder.Encode(Livros)

}

func rotaPrincipal() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello word !")
	}
}

func configurarServidor() {
	configurarHandler()
	fmt.Print("Servidor Start porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	configurarServidor()
}
