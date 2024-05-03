package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/seemsod1/ancy/internal/config"
	"github.com/seemsod1/ancy/internal/handlers"
	"github.com/seemsod1/ancy/internal/models"
	"github.com/seemsod1/ancy/internal/render"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func setup(app *config.AppConfig) error {
	env, err := loadEnv()
	if err != nil {
		return err
	}

	app.Env = env

	db, err := connectDB(env)
	if err != nil {
		return err
	}

	app.DB = db

	if err = runSchemasMigration(db); err != nil {
		return err
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app.Session = session

	repo := handlers.NewRepo(app)
	handlers.NewHandlers(repo)
	render.NewRenderer(app)

	return nil

}

func connectDB(env *config.EnvVariables) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		env.PostgresHost, env.PostgresUser, env.PostgresDBName, env.PostgresPass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	return db, nil
}

func loadEnv() (*config.EnvVariables, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASS")
	postgresDBName := os.Getenv("POSTGRES_DBNAME")

	return &config.EnvVariables{
		PostgresHost:   postgresHost,
		PostgresUser:   postgresUser,
		PostgresPass:   postgresPass,
		PostgresDBName: postgresDBName,
	}, nil
}

func runSchemasMigration(db *gorm.DB) error {

	if err := db.AutoMigrate(&models.ExhibitType{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.ExhibitStatus{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.UserRole{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Exhibit{}); err != nil {
		return err
	}

	if err := addDefaultRoles(db); err != nil {
		return err
	}

	if err := addAdminUser(db); err != nil {
		return err
	}
	if err := addExhibitTypes(db); err != nil {
		return err

	}
	if err := addExhibitStatuses(db); err != nil {
		return err
	}
	return nil

}

func addDefaultRoles(db *gorm.DB) error {
	if count := db.Find(&models.UserRole{}).RowsAffected; count > 0 {
		return nil
	}

	roles := []models.UserRole{
		{
			Name: "User",
		},
		{
			Name: "Admin",
		},
	}

	for _, role := range roles {
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
func addAdminUser(db *gorm.DB) error {

	if count := db.Find(&models.User{Username: "admin"}).RowsAffected; count > 0 {
		return nil
	}

	pass, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		Username: "admin",
		Email:    "vadim@mail.com",
		Password: string(pass),
		RoleID:   2,
	}

	if err = db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
func addExhibitTypes(db *gorm.DB) error {
	if count := db.Find(&models.ExhibitType{}).RowsAffected; count > 0 {
		return nil
	}

	types := []models.ExhibitType{
		{
			Name: "Photo",
		},
		{
			Name: "Video",
		},
		{
			Name: "Audio",
		},
		{
			Name: "Text",
		},
	}

	for _, t := range types {
		if err := db.Create(&t).Error; err != nil {
			return err
		}
	}

	return nil
}
func addExhibitStatuses(db *gorm.DB) error {
	if count := db.Find(&models.ExhibitStatus{}).RowsAffected; count > 0 {
		return nil
	}

	statuses := []models.ExhibitStatus{
		{
			Name: "Pending",
		},
		{
			Name: "Approved",
		},
		{
			Name: "Rejected",
		},
	}

	for _, s := range statuses {
		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}

	return nil
}
