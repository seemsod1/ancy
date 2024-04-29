package config

import (
	"github.com/alexedwards/scs/v2"
	"gorm.io/gorm"
	"html/template"
	"log"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	DB            *gorm.DB
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Env           *EnvVariables
	Session       *scs.SessionManager
}

type EnvVariables struct {
	PostgresHost   string
	PostgresUser   string
	PostgresPass   string
	PostgresDBName string
}
