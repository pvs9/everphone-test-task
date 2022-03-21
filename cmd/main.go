package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/pvs9/everphone-test-task/pkg/handler"
	"github.com/pvs9/everphone-test-task/pkg/migrations"
	"github.com/pvs9/everphone-test-task/pkg/queue"
	"github.com/pvs9/everphone-test-task/pkg/repository"
	"github.com/pvs9/everphone-test-task/pkg/repository/pgx"
	"github.com/pvs9/everphone-test-task/pkg/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := pgx.NewPGXDB(pgx.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize db: %s", err.Error())
	}

	initModels(db)

	q, err := queue.NewSQSQueue(queue.Config{
		Options: session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("AWS_ACCESS_KEY_ID"),
					os.Getenv("AWS_SECRET_ACCESS_KEY"),
					"",
				),
				Endpoint: aws.String(viper.GetString("queue.host")),
				Region:   aws.String(viper.GetString("queue.region")),
			},
		},
		QueueName: viper.GetString("queue.name"),
		ConsumerConfig: queue.ConsumerConfig{
			MaxNumberOfMessages:      viper.GetInt64("queue.consumer.max_messages"),
			MessageVisibilityTimeout: viper.GetInt64("queue.consumer.visibility_timeout"),
			PollDelayInMilliseconds:  viper.GetInt("queue.consumer.poll_delay"),
			Receivers:                viper.GetInt("queue.consumer.receivers"),
		},
	})

	if err != nil {
		log.Fatalf("failed to initialize queue: %s", err.Error())
	}

	repositories := repository.NewRepository(db)
	queues := queue.NewQueue(q)
	services := service.NewService(queues, repositories)
	handlers := handler.NewHandler(services)

	app := &cli.App{
		Name: "app",
		Commands: []*cli.Command{
			initAppCommand(db, handlers),
			initConsumerCommand(queues, services),
			newDBCommand(db, migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[FATAL] Failed to initialize app: %s", err.Error())
	}
}

func initAppCommand(db *bun.DB, handlers *handler.Handler) *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "start API server",
		Action: func(c *cli.Context) error {
			srv := new(christmas.Server)

			go func() {
				if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
					log.Fatalf("error occured while running http server: %s", err.Error())
				}
			}()

			log.Print("Application started and running")

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
			<-quit

			log.Print("Application shutting down")

			if err := srv.Shutdown(context.Background()); err != nil {
				log.Errorf("error occured on server shutting down: %s", err.Error())
			}

			if err := db.Close(); err != nil {
				log.Errorf("error occured on db connection close: %s", err.Error())
			}

			return nil
		},
	}
}

func initConsumerCommand(queues *queue.Queue, services *service.Service) *cli.Command {
	return &cli.Command{
		Name:  "consume",
		Usage: "start consuming messages from queue",
		Action: func(c *cli.Context) error {
			queues.Consumer.Consume(services.DatasetMessageHandler)

			log.Print("Application consumer started and running")
			return nil
		},
	}
}

func newDBCommand(db *bun.DB, migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					migrator := migrate.NewMigrator(db, migrations)

					return migrator.Init(context.Background())
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					migrator := migrate.NewMigrator(db, migrations)

					group, err := migrator.Migrate(context.Background())
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("There are no new migrations to run\n")
						return nil
					}

					fmt.Printf("Successfully migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					migrator := migrate.NewMigrator(db, migrations)

					group, err := migrator.Rollback(context.Background())
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("There are no groups to roll back\n")
						return nil
					}

					fmt.Printf("Successfully rolled back %s\n", group)
					return nil
				},
			},
		},
	}
}

func initConfig() error {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func initModels(db *bun.DB) {
	db.RegisterModel((*christmas.EmployeeTag)(nil))
	db.RegisterModel((*christmas.GiftTag)(nil))
}
