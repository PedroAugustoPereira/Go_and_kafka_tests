package akafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Consume(topics []string, servers string, msgChan chan *kafka.Message) {

	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,                 //servidor de comunicação que vamos falar com o kafka
		"group.id":          "imersao12-go-esquenta", //grupo de consumidores caso tenha várias intâncias
		"auto.offset.reset": "earliest",              //pegar mensagens sempre od início
	})

	if err != nil {
		print("erro, não conectou com kafka")
		panic(err)
	}

	kafkaConsumer.SubscribeTopics(topics, nil) //esses são os tópicos que vamos consumir, onde vamos nos inscrever de fato

	//Agora precisamos rodar um looping infinito para ficar lendo.
	for {
		msg, err := kafkaConsumer.ReadMessage(-1) //vamos ler a mensagem, -1 é um timeout infinito
		if err == nil {
			msgChan <- msg //Toda vez que lermos uma mensagem nesse canal, a outra thread vai saber que temos algo para processar
		}
	}
}

// Quando rodar essa função, vamos começar a ler e receber dados do apache kafka
