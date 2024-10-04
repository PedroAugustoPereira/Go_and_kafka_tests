package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"productSystem/internal/infra/akafka"
	"productSystem/internal/infra/repository"
	"productSystem/internal/infra/web"
	"productSystem/internal/usecase"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/products") //Assim nos conectamos ao banco de dados

	if err != nil {
		print("erro, não conectou com bd")
		panic(err)
	}
	defer db.Close() //fecha a conexão com o banco de dados

	//Agora vamos criar o canal para receber mensagens do kafka
	//Ele nos ajuda a trabalhar em threads diferentes
	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"product"}, "host.docker.internal:9094", msgChan) //Assima estamos lendo nesse canal as mensagens do kafka
	//Criamos uam thrad para rodar de forma concorrente

	//Agora vamos criar nosso objeto do repositório
	repository := repository.NewProductRepositoryMysql(db)

	//Agora vamos criar nosso useCase
	createProductUsecase := usecase.NewCreateProductUseCase(repository)
	listProductsUseCase := usecase.NewListProductsUseCase(repository)

	//Criar productHanlders
	productHandlers := web.NewProductHandlers(createProductUsecase, listProductsUseCase)

	//Criamos um roteamento
	//E estabelecemos nossas rotas
	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductsHandler)

	//Vamos aqui criar um servidor http em uma nova thread, para contniuar a execução na main
	go http.ListenAndServe(":8000", r)
	print("server on")

	//Agora vamos fazer um for aqui, para ficar lendo as mensagens que ficamos recebendo do kafka
	//Lembrando aqui, esse for fica bloqueado enquanto não tiver nenhuma menasgem no canal
	for msg := range msgChan {
		//Criamos um dto de input
		//Para hidradatar esse dto com dados em json, preciso colocar umas notações la na struct
		//`json:"name"`
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto) //converte a mensagem de json para a struct dto que é passada por referencia

		if err != nil {
			//logar o erro aqui
		}

		//Executamos a criação do produto
		_, err = createProductUsecase.Execute(dto)
	}
}
