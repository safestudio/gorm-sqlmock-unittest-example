package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Rosaniline/gorm-ut/pkg/model"
	"github.com/go-test/deep"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository Repository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(
		mysql.New(mysql.Config{Conn: db,
			SkipInitializeWithVersion: true,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	require.NoError(s.T(), err)

	s.repository = CreateRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_repository_Get() {
	var (
		id   = uuid.New()
		name = "test-name"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `person` WHERE id = ?")).
		WithArgs(id.String()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id.String(), name))
	res, err := s.repository.Get(id)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&model.Person{ID: id, Name: name}, res))
}

func (s *Suite) Test_repository_Create() {
	var (
		id   = uuid.New()
		name = "test-name"
	)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `person` (`id`,`name`) VALUES (?,?)")).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.Create(id, name)

	require.NoError(s.T(), err)
}
