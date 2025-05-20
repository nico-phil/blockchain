package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/Nico2220/blockchain/utils"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	// 1. Creating ECDSA private key(32 bytes) public key(64 bytes)
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &privateKey.PublicKey

	//2. Perform SHA-256
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	//3
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	//4
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	//5
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	//6
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	//7
	chsum := digest6[:4]

	//8
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum)

	// 9
	address := base58.Encode(dc8)
	w.blockchainAddress = address
	return w
}



func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func(w *Wallet) MarshalJSON()([]byte, error){
	return json.Marshal(struct{
		PrivateKey string `json:"private_key"`
		PublicKey string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey: w.PrivateKeyStr(),
		PublicKey: w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

type Transaction struct {
	senderPrivateKey      *ecdsa.PrivateKey
	senderPublickKey      *ecdsa.PublicKey
	sendBlockchainAddress string
	recepientBlockAddress string
	value                 float32
}

func NewTransaction(
	privateKey *ecdsa.PrivateKey,
	publicKey *ecdsa.PublicKey,
	sender string,
	recipient string,
	value float32,
) *Transaction {
	return &Transaction{
		privateKey,
		publicKey,
		sender,
		recipient,
		value,
	}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshaJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.sendBlockchainAddress,
		Recipient: t.recepientBlockAddress,
		Value:     t.value,
	})
}


type TransactionRequest struct {
	SenderPrivateKey  *string   `json:"sender_private_key"`   
	SenderPublickKey     *string `json:"sender_public_key"`
	SendBlockchainAddress *string `json:"sender_blockchain_address"`
	RecepientBlockAddress *string `json:"recipient_block_address"`
	Value   *string  `json:"value"`
}