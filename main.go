package main

import (
	"encoding/json" // pacote para trabalhar com JSON
	"fmt"
	"io/ioutil"
	"log"
	"net/http" // pacote para criar um servidor HTTP
	"strconv"  // pacote para converter string em número
	"strings"  // pacote para trabalhar com strings
)

// Struct para representar um livro
type Livro struct {
	ID     int    `json:"id"`     // A tag json é usada para definir o nome da chave no JSON
	Titulo string `json:"titulo"` // A tag json é usada para definir o nome da chave no JSON
	Autor  string `json:"autor"`  // A tag json é usada para definir o nome da chave no JSON
}

// Lista de livros para simular um banco de dados
var livros = []Livro{
	{ID: 1, Titulo: "O Guarani", Autor: "José de Alencar"},
	{ID: 2, Titulo: "Iracema", Autor: "José de Alencar"},
	{ID: 3, Titulo: "Dom Casmurro", Autor: "Machado de Assis"},
	{ID: 4, Titulo: "A Hora da Estrela", Autor: "Clarice Lispector"},
	{ID: 5, Titulo: "Grande Sertão: Veredas", Autor: "Guimarães Rosa"},
}

// Função que será chamada quando o usuário acessar a rota principal
func rotaPrincipal(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bem-vindo!"))
}

// Função que será chamada quando o usuário acessar a rota /livros
func listarLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o header da resposta para indicar que estamos retornando um JSON
	w.WriteHeader(http.StatusCreated)                  // Define o status da resposta para 201 (Created)
	json.NewEncoder(w).Encode(livros)                  // Codifica a lista de livros em JSON e envia para o cliente
}

// Função que será chamada quando o usuário enviar um POST para a rota /livros
func cadastrarLivro(w http.ResponseWriter, r *http.Request) {
	var livro Livro
	err := json.NewDecoder(r.Body).Decode(&livro) // Decodifica o JSON enviado pelo cliente e armazena na variável livro
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Retorna um erro 400 (Bad Request) se houver algum erro na decodificação do JSON
		return
	}
	livro.ID = len(livros) + 1                         // Atribui um ID para o livro (neste caso, o próximo número inteiro disponível)
	livros = append(livros, livro)                     // Adiciona o livro à lista de livros
	w.Header().Set("Content-Type", "application/json") // Define o header da resposta para indicar que estamos retornando um JSON
	w.WriteHeader(http.StatusCreated)                  // Define o status da resposta para 201 (Created)
	json.NewEncoder(w).Encode(livro)                   // Codifica o livro em JSON e envia para o cliente
}

// Função que será chamada quando o usuário enviar um DELETE para a rota /livros/{id}
func excluirLivro(w http.ResponseWriter, r *http.Request) {
	// Extrai o id da rota e converte de string para int
	partes := strings.Split(r.URL.Path, "/")
	if len(partes) < 3 {
		w.WriteHeader(http.StatusBadRequest) // Retorna um erro 400 (Bad Request) se o ID
		return
	}
	id, err := strconv.Atoi(partes[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Busca o livro com o id fornecido
	var indiceDoLivro = -1
	for indice, livro := range livros {
		if livro.ID == id {
			indiceDoLivro = indice
			break
		}
	}
	// Retorna um erro 404 se o livro não foi encontrado
	if indiceDoLivro < 0 || indiceDoLivro >= len(livros) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove o livro da lista
	livros = append(livros[:indiceDoLivro], livros[indiceDoLivro+1:]...)

	w.WriteHeader(http.StatusOK)
}

func rotearLivros(w http.ResponseWriter, r *http.Request) {

	//divide o caminho da URL recebida em partes
	partes := strings.Split(r.URL.Path, "/")
	// Se a URL possui 2 partes ou se possui 3 partes, sendo que a terceira parte é vazia, a função entra no switch
	if len(partes) == 2 || len(partes) == 3 && partes[2] == "" {
		switch r.Method {
		case http.MethodGet:
			listarLivros(w, r)
		case http.MethodPost:
			cadastrarLivro(w, r)
		case http.MethodPut:
			modificarLivro(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}

	} else if len(partes) == 3 || len(partes) == 4 && partes[3] == "" {
		if r.Method == http.MethodGet {
			buscarLivros(w, r)
		} else if r.Method == http.MethodDelete {
			excluirLivro(w, r)
		} else if r.Method == http.MethodPut {
			modificarLivro(w, r)
		}

	} else {
		//Se nenhuma das duas condições anteriores for verdadeira, significa que a rota não existe e retorna:
		w.WriteHeader(http.StatusNotFound)
	}
}
func buscarLivros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.URL.Path)

	//"/livros/123" ---> ["", "livros", "123"]
	partes := strings.Split(r.URL.Path, "/")

	//extrai o id e converte o int para string
	id, err := strconv.Atoi(partes[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//encontrar o livro na lista
	for _, livro := range livros {
		if livro.ID == id {
			json.NewEncoder(w).Encode(livro)
			return
		}
	}
	http.NotFound(w, r)
}
func modificarLivro(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(partes[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, erroBody := ioutil.ReadAll(r.Body)
	if erroBody != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var livroModificado Livro
	erroJson := json.Unmarshal(body, &livroModificado)

	if erroJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indiceDoLivro := -1
	for indice, livro := range livros {
		if livro.ID == id {
			indiceDoLivro = indice
			break
		}
	}
	if indiceDoLivro < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	livros[indiceDoLivro] = livroModificado

	json.NewEncoder(w).Encode(livroModificado)
}
func main() {
	http.HandleFunc("/", rotaPrincipal)
	http.HandleFunc("/livros", rotearLivros)
	http.HandleFunc("/livros/", rotearLivros)
	log.Fatal(http.ListenAndServe(":5558", nil))
}
