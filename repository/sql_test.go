package repository

import (
	"context"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

type Test struct {
	gorm.Model
	FirstName string
	LastName  string
}

func TestSQL(t *testing.T) {
	suite.Run(t, new(SQLTestSuite))
}

type SQLTestSuite struct {
	suite.Suite

	db *gorm.DB

	tx *gorm.DB

	repository Repository[Test, uint]
}

func (s *SQLTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open("sql_test.db"))
	s.Require().NoError(err)
	s.Require().NotNil(db)
	s.db = db
}

func (s *SQLTestSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.tx = s.tx.Debug()
	s.repository = &SQL[Test, uint]{
		db: s.tx,
	}
	s.Require().NoError(s.tx.Migrator().AutoMigrate(&Test{}))
}

func (s *SQLTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *SQLTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	s.Require().NoError(sqlDB.Close())
	s.Require().NoError(os.Remove("sql_test.db"))
}

func (s *SQLTestSuite) createMockData() {
	s.Require().NoError(s.tx.Create(&Test{
		FirstName: "Marcos",
		LastName:  "Huck",
	}).Error)

	s.Require().NoError(s.tx.Create(&Test{
		FirstName: "Andres",
		LastName:  "Huck",
	}).Error)

	s.Require().NoError(s.tx.Create(&Test{
		FirstName: "Andrew",
		LastName:  "Baker",
	}).Error)
}

func (s *SQLTestSuite) TestFind() {
	s.createMockData()
	list, err := s.repository.Find(context.Background(), []uint{1, 2, 3})
	s.Assert().NoError(err)
	s.Assert().Len(list, 3)
}

func (s *SQLTestSuite) TestGet() {
	s.createMockData()
	result, err := s.repository.Get(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().NotZero(result)
	s.Assert().Equal("Marcos", result.FirstName)
}

func (s *SQLTestSuite) TestGet_NotFound() {
	s.createMockData()
	result, err := s.repository.Get(context.Background(), 5)
	s.Assert().Error(err)
	s.Assert().Zero(result)
}

func (s *SQLTestSuite) TestUpdate() {
	s.createMockData()
	result, err := s.repository.Update(context.Background(), 1, Test{FirstName: "Changed"})
	s.Assert().NoError(err)
	s.Assert().NotZero(result)

	result, err = s.repository.Get(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().Equal("Changed", result.FirstName)
	s.Assert().Equal("Huck", result.LastName)

	result, err = s.repository.Get(context.Background(), 2)
	s.Assert().NoError(err)
	s.Assert().Equal("Andres", result.FirstName)
	s.Assert().Equal("Huck", result.LastName)
}

func (s *SQLTestSuite) TestRemove() {
	s.createMockData()
	result, err := s.repository.Get(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().NotZero(result)

	removed, err := s.repository.Remove(context.Background(), 1)
	s.Assert().NotZero(removed)
	s.Assert().NoError(err)
	s.Assert().NotZero(removed.DeletedAt)
	s.Assert().Equal(result.ID, removed.ID)

	result, err = s.repository.Get(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Zero(result)
}

func (s *SQLTestSuite) TestCreateBulk() {
	result, err := s.repository.Get(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Zero(result)

	created, err := s.repository.CreateBulk(context.Background(), []Test{
		{
			FirstName: "Marcos",
			LastName:  "Huck",
		},
		{
			FirstName: "Marcos",
			LastName:  "Huck",
		},
		{
			FirstName: "Marcos",
			LastName:  "Huck",
		},
	})
	s.Assert().NoError(err)
	s.Assert().NotEmpty(created)
	s.Assert().Len(created, 3)

	res, err := s.repository.Find(context.Background(), []uint{1, 2, 3})
	s.Assert().NoError(err)
	s.Assert().Len(res, 3)
}
