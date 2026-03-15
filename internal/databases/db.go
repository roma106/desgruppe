package databases

import (
	"desgruppe/internal/config"
	"desgruppe/internal/logger"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDB(cfg *config.Config) *sqlx.DB {
	logger.Info(fmt.Sprintf("Connecting to database %s...", cfg.PostgresDB))
	query := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	db, err := sqlx.Open("postgres", query)
	if err != nil {
		logger.Error("Failed to connect to database: " + err.Error())
		return nil
	}
	err = db.Ping()
	if err != nil {
		logger.Error("Failed to connect to database: " + err.Error())
		return nil
	}
	logger.Info("Database " + cfg.PostgresDB + " connected!")
	return db
}

func InitProducts(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.products
	(
		"ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 ),
		"Name" character varying COLLATE pg_catalog."default" NOT NULL,
		"Type" character varying COLLATE pg_catalog."default" NOT NULL,
		"SectionID" integer,
		"ProducerID" integer,
		"DesignerID" integer,
		"Size" character varying COLLATE pg_catalog."default",
		"Available" boolean NOT NULL,
		"Price" decimal NOT NULL,
		"OnSale" boolean NOT NULL,
		"Sale" decimal,
		"Photo" bytea,
		"Description" text COLLATE pg_catalog."default",
		"Colors" integer[],
		"Position" integer,
		"Slug" character varying UNIQUE COLLATE pg_catalog."default",
		"Free Form" character varying,
		CONSTRAINT products_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Products table ready!")
	return nil
}

func InitColors(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.colors
	(
		"ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
		"Name" character varying COLLATE pg_catalog."default",
		"Code" character varying COLLATE pg_catalog."default",
		"Picture" bytea,
		"Position" bigint,
		CONSTRAINT colors_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Colors table ready!")
	return nil
}

func InitDesigners(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.designers
	(
		"ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
		"Name" character varying COLLATE pg_catalog."default",
		"Description" text COLLATE pg_catalog."default",
		"Picture" bytea,
		"Position" bigint,
		"Slug" character varying UNIQUE COLLATE pg_catalog."default",
		CONSTRAINT designers_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Designers table ready!")
	return nil
}

func InitProducers(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.producers
	(
		"ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
		"Name" character varying COLLATE pg_catalog."default",
		"Description" text COLLATE pg_catalog."default",
		"Picture" bytea,
		"Position" bigint,
		"Slug" character varying UNIQUE COLLATE pg_catalog."default",
		CONSTRAINT producers_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Producers table ready!")
	return nil
}

func InitSections(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.sections
	(
		"ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
		"Name" character varying COLLATE pg_catalog."default",
		"Type" character varying COLLATE pg_catalog."default",
		CONSTRAINT sections_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Sections table ready!")
	return nil
}

func InitProductsGallery(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.products_gallery
	(
		"ID" serial NOT NULL,
		"ProductID" integer NOT NULL,
		"Photo" bytea,
		"Position" integer,
		CONSTRAINT products_gallery_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Gallery table ready!")
	return nil
}

func InitSettings(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.settings
	(
		"ID" serial NOT NULL,
		"ExchangeRate" decimal,
		"Email" character varying,
		CONSTRAINT settings_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Settings table ready!")
	return nil
}
func InitCarts(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.carts
	(
		"ID" serial NOT NULL,
		"ProductIDs" integer[],
		"ColorIDs" integer[],
		"QuantitiesIDs" integer[],
		CONSTRAINT carts_pkey PRIMARY KEY ("ID")
	);
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Carts table ready!")
	return nil
}

func InitOrders(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.orders
	(
		"ID" serial NOT NULL,
		"Date" date,
		"Name" character varying,
		"Phone" character varying NOT NULL,
		"Email" character varying,
		"Comment" character varying,
		"CartID" character varying NOT NULL,
		"Seen" boolean,
		CONSTRAINT orders_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Orders table ready!")
	return nil
}

func InitQuestions(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS public.questions
	(
		"ID" serial NOT NULL,
		"Date" date,
		"Name" character varying,
		"Phone" character varying NOT NULL,
		"Email" character varying,
		"Message" character varying,
		"Seen" boolean,
		CONSTRAINT questions_pkey PRIMARY KEY ("ID")
	)
	`)
	if err != nil {
		logger.Error("Failed to create table: " + err.Error())
		return err
	}
	logger.Info("Questions table ready!")
	return nil
}
