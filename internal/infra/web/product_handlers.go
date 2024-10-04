package web

import (
	"encoding/json"
	"net/http"
	"productSystem/internal/usecase"
)

// Struct do objeto
type ProductHandlers struct {
	CreateProductUseCase *usecase.CreateProductUseCase
	ListProductsUseCase  *usecase.ListProductsUseCase
}

// Construtor, recebemos os parametros por ponteiros de valor e retornamos o ponteiro para uma referencia do ProductHandlers
func NewProductHandlers(createProductUseCase *usecase.CreateProductUseCase, listProductsUseCase *usecase.ListProductsUseCase) *ProductHandlers {
	return &ProductHandlers{
		CreateProductUseCase: createProductUseCase,
		ListProductsUseCase:  listProductsUseCase,
	}
}

// Agora criamos um controller entre aspas de fato, para receber requisições http
// Esse cara não sabe o que ta acontecendo na aplicação, ele apenas utiliza ela
// é um extremo de uma arquiterua  exdagonal
func (p *ProductHandlers) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	print("chegou aqui")
	var input usecase.CreateProductInputDto       //Criamos a nossa referencia a struct de input
	err := json.NewDecoder(r.Body).Decode(&input) //Criamos um newDecoder json e fazemos o decode de r.Body e passamos a referencia de input para ser atualizado
	if err != nil {
		print("Erro no decoder do json\n")
		w.WriteHeader(http.StatusBadRequest) //Escxrevemos na requisição que tivemos um bad Request
		return
	}
	output, err := p.CreateProductUseCase.Execute(input) //Executamos o CreateProductUsecase
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //Se tiver um erro, damos um erro de internal server com w response
		return
	}

	w.Header().Set("Content-Type", "application/json") //setamos o retorno como aplication json
	w.WriteHeader(http.StatusCreated)                  //escrevemos um estado como created
	json.NewEncoder(w).Encode(output)                  //fazemops um encoder do w, e fazemos o encode do outpt com o produto criado, com output sendo do tipo dto de outp
}

// Agora vamos fazer um endpoint para listar todos os produtos, da mesma forma usando requisições http
func (p *ProductHandlers) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	output, err := p.ListProductsUseCase.Execute()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
