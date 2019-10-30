package auth

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/keyerror"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keybase struct {
	name  string
	addr  sdk.AccAddress
	armor string
}

func NewHotKeybase(name string, passphrase string, pk crypto.PrivKey) *Keybase {
	armor := mintkey.EncryptArmorPrivKey(pk, passphrase)
	addr := sdk.AccAddress(pk.PubKey().Address())

	return &Keybase{
		name:  name,
		addr:  addr,
		armor: armor,
	}
}

func (k *Keybase) GetAddr() sdk.AccAddress {
	return k.addr
}

func (k *Keybase) GetName() string {
	return k.name
}

func (*Keybase) List() ([]keys.Info, error) {
	panic("not implemented")
}

func (*Keybase) Get(name string) (keys.Info, error) {
	panic("not implemented")
}

func (*Keybase) GetByAddress(address types.AccAddress) (keys.Info, error) {
	panic("not implemented")
}

func (*Keybase) Delete(name, passphrase string, skipPass bool) error {
	panic("not implemented")
}

func (k *Keybase) Sign(name string, passphrase string, msg []byte) ([]byte, crypto.PubKey, error) {
	if k.name != name {
		return nil, nil, keyerror.NewErrKeyNotFound(name)
	}

	priv, err := mintkey.UnarmorDecryptPrivKey(k.armor, passphrase)
	if err != nil {
		return nil, nil, err
	}

	sig, err := priv.Sign(msg)
	if err != nil {
		return nil, nil, err
	}

	pub := priv.PubKey()
	return sig, pub, nil
}

func (*Keybase) CreateMnemonic(name string, language keys.Language, passwd string, algo keys.SigningAlgo) (info keys.Info, seed string, err error) {
	panic("not implemented")
}

func (*Keybase) CreateAccount(name, mnemonic, bip39Passwd, encryptPasswd string, account uint32, index uint32) (keys.Info, error) {
	panic("not implemented")
}

func (*Keybase) Derive(name, mnemonic, bip39Passwd, encryptPasswd string, params hd.BIP44Params) (keys.Info, error) {
	panic("not implemented")
}

func (*Keybase) CreateOffline(name string, pubkey crypto.PubKey) (info keys.Info, err error) {
	panic("not implemented")
}

func (*Keybase) CreateMulti(name string, pubkey crypto.PubKey) (info keys.Info, err error) {
	panic("not implemented")
}

func (*Keybase) Update(name, oldpass string, getNewpass func() (string, error)) error {
	panic("not implemented")
}

func (*Keybase) Import(name string, armor string) (err error) {
	panic("not implemented")
}

func (*Keybase) ImportPubKey(name string, armor string) (err error) {
	panic("not implemented")
}

func (*Keybase) Export(name string) (armor string, err error) {
	panic("not implemented")
}

func (*Keybase) ExportPubKey(name string) (armor string, err error) {
	panic("not implemented")
}

func (*Keybase) ExportPrivateKeyObject(name string, passphrase string) (crypto.PrivKey, error) {
	panic("not implemented")
}

func (*Keybase) CloseDB() {
	panic("not implemented")
}

func (k *Keybase) CreateLedger(name string, algo keys.SigningAlgo, hrp string, account, index uint32) (info keys.Info, err error) {
	panic("not implemented")
}

func (k *Keybase) ImportPrivKey(name, armor, passphrase string) error {
	panic("not implemented")
}

func (k *Keybase) ExportPrivKey(name, decryptPassphrase, encryptPassphrase string) (armor string, err error) {
	panic("not implemented")
}
