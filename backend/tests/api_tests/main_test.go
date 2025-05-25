package api_tests

import (
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/testutils"
)

func TestMain(m *testing.M) {
	// Initialize the test suite with the database
	testutils.InitializeTestSuite(m)
}
