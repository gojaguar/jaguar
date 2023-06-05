package repository

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type gormTestSuite struct {
	suite.Suite
}

func TestGorm(t *testing.T) {
	suite.Run(t, new(gormTestSuite))
}
