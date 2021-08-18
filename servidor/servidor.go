package servidor

import (
	"crud/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:id`
	Nome  string `json:nome`
	Email string `json:email`
}

func CriarUsuario(rw http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		rw.Write([]byte("Bad Request"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		rw.Write([]byte("Erro ao converter usuario - Bad Request"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		rw.Write([]byte("Erro ao conectar ao banco de dados - Erro Interno"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if erro != nil {
		rw.Write([]byte("Erro ao criar o statement - Erro Interno"))
		return
	}

	defer statement.Close()

	insercao, erro := statement.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		rw.Write([]byte("Erro ao executar o statement"))
		return
	}
	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		rw.Write([]byte("Erro ao obter o último id inserido"))
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso! Id: %d", idInserido)))

}

func BuscarUsuarios(rw http.ResponseWriter, r *http.Request) {
	db, erro := db.Conectar()
	if erro != nil {
		rw.Write([]byte("Erro ao conectar ao banco de dados - Erro Interno"))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from usuarios")
	if erro != nil {
		rw.Write([]byte("Erro ao buscar usuários"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario

	for linhas.Next() {
		var usuario usuario

		if err := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			rw.Write([]byte("Erro ao escanear o usuário!"))
			return
		}

		usuarios = append(usuarios, usuario)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(rw).Encode(usuarios); erro != nil {
		rw.Write([]byte("Erro ao converter usuários em JSON"))
		return
	}
}

func BuscarUsuario(rw http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
		rw.Write([]byte("Parâmetro invalido"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		rw.Write([]byte("Erro ao conectar ao banco de dados - Erro Interno"))
		return
	}
	defer db.Close()

	linha, erro := db.Query("select * from usuarios where id = ?", ID)
	if erro != nil {
		rw.Write([]byte("Erro ao buscar usuários"))
		return
	}
	defer linha.Close()

	var usuario usuario

	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			rw.Write([]byte("Erro ao buscar usuário"))
			return
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(rw).Encode(usuario); erro != nil {
		rw.Write([]byte("Erro ao converter usuário em JSON"))
		return
	}
}

func AtualizarUsuario(rw http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
		rw.Write([]byte("Parâmetro invalido"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		rw.Write([]byte("Bad Request"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		rw.Write([]byte("Erro ao converter usuario - Bad Request"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		rw.Write([]byte("Erro ao conectar ao banco de dados - Erro Interno"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if erro != nil {
		rw.Write([]byte("Erro ao criar o statement - Erro Interno"))
		return
	}

	defer statement.Close()

	if _, erro := statement.Exec(usuario.Nome, usuario.Email, ID); erro != nil {
		rw.Write([]byte("Erro ao criar o statement - Erro Interno"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func DeletarUsuario(rw http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
		rw.Write([]byte("Parâmetro invalido"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		rw.Write([]byte("Erro ao conectar ao banco de dados - Erro Interno"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("delete from usuarios  where id = ?")
	if erro != nil {
		rw.Write([]byte("Erro ao criar o statement - Erro Interno"))
		return
	}

	if _, erro := statement.Exec(ID); erro != nil {
		rw.Write([]byte("Erro ao deletar o usuário - Erro Interno"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
