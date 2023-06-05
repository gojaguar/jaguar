package repository

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type firestoreTestSuite struct {
	suite.Suite
}

func TestFirestore(t *testing.T) {
	suite.Run(t, new(firestoreTestSuite))
}
