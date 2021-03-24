package consumer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConsumer_CreateConsumer_ShouldReturnErrorOnSubscribeTopicError(t *testing.T) {
	c, err := CreateConsumer("test", "test", "")
	require.Nil(t, c)
	require.NotNil(t, err)
}

func TestConsumer_CreateConsumer_ShouldReturnNilErrorOnSuccess(t *testing.T) {
	c, err := CreateConsumer("test", "test", "test")
	require.Nil(t, err)
	require.NotNil(t, c)
}
