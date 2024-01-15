package beta

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestExpiredPromotion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up mock expectations for the database query
	rows := sqlmock.NewRows([]string{"expired"}).AddRow("2024-01-15")
	mock.ExpectQuery("SELECT expired FROM tbl_beta_promotion LIMIT 1").WillReturnRows(rows)

	// Create a real GORM instance using the mocked database
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	// Call the function to be tested
	data := ExpiredPromotion(gormDB)

	// Assert the results
	if data.Expired != "2024-01-15" {
		t.Errorf("Expected expired value to be '2024-01-15', got '%s'", data.Expired)
	}

	// Ensure all mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
