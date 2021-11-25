package modules

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/math"
	"strings"
)

type OnBoarding struct {
	Host                      string
	EthSigner                 EthSigner
	NetworkId                 int
	EthAddress                string
	StarkPublicKey            string
	StarkPublicKeyYCoordinate string
	Singer                    *SignOnboardingAction
}

type ApiKeyCredentials struct {
	Key        string
	Secret     string
	Passphrase string
}

func (board OnBoarding) RecoverDefaultApiCredentials(ethereumAddress string) ApiKeyCredentials {
	signature := board.Singer.Sign(ethereumAddress, map[string]interface{}{"action": OffChainOnboardingAction})
	rHex := signature[2:66]
	rInt, _ := math.MaxBig256.SetString(rHex, 16)

	rData := [][]interface{}{{"uint256"}, {rInt.String()}}

	keccak := SolidityKeccak(rData)
	hashedRBytes := []byte(keccak)
	secretBytes := hashedRBytes[:30]
	sHex := signature[66:130]
	sInt, _ := math.MaxBig256.SetString(sHex, 16)
	sData := [][]interface{}{{"uint256"}, {sInt.String()}}
	hashedSBytes := []byte(SolidityKeccak(sData))
	keyBytes := hashedSBytes[:16]
	passphraseBytes := hashedSBytes[16:31]

	keyHex := hex.EncodeToString(keyBytes)
	keyUuid := strings.Join([]string{
		keyHex[:8],
		keyHex[8:12],
		keyHex[12:16],
		keyHex[16:20],
		keyHex[20:],
	}, "-")

	return ApiKeyCredentials{
		Key:        base64.URLEncoding.EncodeToString(secretBytes),
		Secret:     keyUuid,
		Passphrase: base64.URLEncoding.EncodeToString(passphraseBytes),
	}
}

func (board OnBoarding) sign(signerAddress, action string) string {
	return board.Singer.Sign(signerAddress, map[string]interface{}{"action": action})
}
