package tests

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

var testConfig TestConfig

type TestConfig struct {
	adapter Client
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		TestSuiteInitializer: InitTestSuite,
		ScenarioInitializer:  InitScenarios,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(
		`^a user creates a full deck that is not shuffled$`, aUserCreatesAFullDeckThatIsNotShuffled,
	)
	ctx.Step(`^a user creates a full deck that is shuffled$`, aUserCreatesAFullDeckThatIsShuffled)
	ctx.Step(
		`^a user creates a partial deck that is not shuffled with the following cards:$`,
		aUserCreatesAPartialDeckThatIsNotShuffledWithTheFollowingCards,
	)
	ctx.Step(
		`^the user should receive a deck ID and the following results:$`,
		theUserShouldReceiveADeckIDAndTheFollowingResults,
	)
	ctx.Step(`^the user requests to open the created deck$`, theUserRequestsToOpenTheCreatedDeck)
	ctx.Step(
		`^the user should receive a deck with the following results:$`,
		theUserShouldReceiveADeckWithTheFollowingResults,
	)
}

func InitTestSuite(ctx *godog.TestSuiteContext) {
	clientInst := NewClient("http://localhost:8080/api/v1", time.Second*10)
	testConfig = TestConfig{adapter: clientInst}
	dsn := fmt.Sprintf("host=localhost user=toggl password=toggl dbname=cards port=5432 sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	runMigrations(db)
}

func runMigrations(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgresMigrate.WithInstance(sqlDB, &postgresMigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Println("Database wasn't dropped", "err", err.Error())
	}
	if err := m.Up(); err != nil {
		log.Println("Migrations weren't updated", "err", err.Error())
	}
}