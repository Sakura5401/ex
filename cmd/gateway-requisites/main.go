package main

import (
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"gitlab.safecrow.ru/safecrow/gateway-requisites/v2/internal/handlers"
	"gitlab.safecrow.ru/safecrow/gateway-requisites/v2/internal/models"
	"gitlab.safecrow.ru/safecrow/gateway-requisites/v2/utils"
	rabbitmq "gitlab.safecrow.ru/safecrow/gateway-rmq"
)

var settings = utils.GetEnvs() //nolint
var rmq *amqp.Connection       //nolint

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Str("handler", "Requisites").Msg("Init service")

	_ = godotenv.Load()

	// подключиние к базе данных (PostgreSQL)
	db, err := models.Connection(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal().Err(err)
	}

	err = models.ModifyColumnIfTableExist(db)
	if err != nil {
		log.Fatal().Err(err)
	}

	err = models.CreateTableIfNotExists(db)
	if err != nil {
		log.Fatal().Err(err)
	}

	// подключиние к брокеру сообщений (RabbitMQ)
	if rmq, err = rabbitmq.Connect(os.Getenv("AMQP_URL")); err != nil {
		log.Fatal().Err(err).Str("handler", "Requisites").Msg("Connect to rmq")
	}

	defer func() {
		_ = rmq.Close()
	}()

	for i := 1; i <= settings.WorkerCount; i++ {
		log.Info().Msgf("Run worker: %d", i)

		go worker(db, rmq, log.Logger)
	}

	log.Info().Msg("[*] Waiting for messages. To exit press CTRL+C")

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal)
	msgSignal := <-quitSignal
	log.Info().Msgf("Got signal: %s", msgSignal.String())
}

// Данный метод создаёт очередь, принимает сообщения из очереди, на основе routing key вызывает хендлер
// и отправляет сообщение(ответ) обратно в очередь
func worker(db *gorm.DB, rmq *amqp.Connection, log zerolog.Logger) {
	// принмает соединение и структуру на подписку, создает новую структуру брокера
	currentQueue := rabbitmq.QueueNew(rmq, rabbitmq.Subscribe{
		ExchangeName:     "gateway",
		ExchangeType:     "topic",
		RoutingKey:       "gateway-requisites.#",
		QueueName:        "gateway-requisites",
		ConsumerTag:      "",
		QosPrefetchCount: 1, //nolint
		QosPrefetchSize:  0,
		QosGlobal:        false,
	})

	//открывает канал для получения сообщений AMQP.
	msgChannel, _, err := currentQueue.QueueUp()
	if err != nil {
		log.Error().Err(err).Str("handler", "Requisites").Msg("subscribe to channel")
		return
	}

	// принимает сообщение (msg) из очереди и на основе routing key вызывает хендлер
	for msg := range msgChannel {
		topics := strings.Split(msg.RoutingKey, ".")

		var got []byte

		switch topics[1] {
		case "create-bank-detail":
			got = handlers.HandlerSetBankDetail(msg.Body, db)
		case "get-bank-details":
			got = handlers.HandlerGetBankDetails(msg.Body, db)
		case "get-bank-detail":
			got = handlers.HandlerGetBankDetail(msg.Body, db)
		case "update-payment-method":
			got = handlers.HandlerUpdateDefaultPayments(msg.Body, db)
		case "create-business":
			got = handlers.HandlerSetBusinessInfo(msg.Body, db)
		case "get-business":
			got = handlers.HandlerGetBusinessInfo(msg.Body, db)
		case "update-business":
			got = handlers.HandlerUpdateBusinessInfo(msg.Body, db)
		case "create-business-contacts":
			got = handlers.HandlerSetBusinessContactInfo(msg.Body, db)
		case "get-business-contacts":
			got = handlers.HandlerGetBusinessContactInfo(msg.Body, db)
		case "update-business-contacts":
			got = handlers.HandlerUpdateBusinessContactInfo(msg.Body, db)
		case "create-bank-card":
			got = handlers.HandlerSetBankCard(msg.Body, db)
		case "get-bank-cards":
			got = handlers.HandlerGetBankCards(msg.Body, db)
		case "get-bank-card":
			got = handlers.HandlerGetBankCard(msg.Body, db)
		case "create-customer":
			got = handlers.HandlerSetCustomerInfo(msg.Body, db)
		case "get-customer":
			got = handlers.HandlerGetCustomerInfo(msg.Body, db)
		case "update-customer":
			got = handlers.HandlerUpdateCustomerInfo(msg.Body, db)
		case "create-customer-contacts":
			got = handlers.HandlerSetCustomerContactInfo(msg.Body, db)
		case "get-customer-contacts":
			got = handlers.HandlerGetCustomerContactInfo(msg.Body, db)
		case "update-customer-contacts":
			got = handlers.HandlerUpdateCustomerContactInfo(msg.Body, db)
		default:
			log.Info().Msgf("Invalid router_key: %s", msg.RoutingKey)

			if err = msg.Ack(false); err != nil {
				panic(err)
			}

			continue
		}

		if msg.ReplyTo == "" {
			msg.ReplyTo = "gateway-bff"
		}

		// отправляет сообщение(ответ) в канал через Exchange по настройкам из Publish
		err = currentQueue.Publish(got, rabbitmq.Publish{
			ExchangeName:  "gateway",
			ExchangeType:  "topic",
			RoutingKey:    msg.ReplyTo + ".",
			QueueName:     msg.ReplyTo,
			CorrelationID: msg.CorrelationId,
			ContentType:   msg.ContentType,
			ReplyTo:       "",
		})

		if err != nil {
			log.Fatal().Err(err)
		}

		if err = msg.Ack(false); err != nil {
			panic(err)
		}
	}
}
