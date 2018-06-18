package signingKey

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {

	dappKey, err := New("My New Key")
	require.Nil(t, err)
	require.Equal(t, "My New Key", dappKey.Name)

	exportedKey, err := dappKey.Export([]byte("pw"))
	require.Nil(t, err)

	recoveredDappKey, err := Decrypt(exportedKey, []byte("pw"))
	require.Nil(t, err)

	require.Equal(t, dappKey.Name, recoveredDappKey.Name)
	require.Equal(t, dappKey.privateKey, recoveredDappKey.privateKey)
	require.Equal(t, dappKey.Version, recoveredDappKey.Version)
	require.Equal(t, dappKey.PublicKey, recoveredDappKey.PublicKey)

}
