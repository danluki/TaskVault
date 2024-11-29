package taskvault

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestLogger() *logrus.Entry {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	entry := logrus.NewEntry(log)
	return entry
}

func TestStore(t *testing.T) {
	log := getTestLogger()
	s, err := NewStore(log)
	require.NoError(t, err)
	defer s.Shutdown() // nolint: errcheck

	err = s.SetValue("foo", "bar")
	assert.NoError(t, err)
	value, err := s.GetValue("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", value)
}
