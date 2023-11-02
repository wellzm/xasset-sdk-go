package auth

import (
	"strings"

	"github.com/xuperchain/crypto/client/service/xchain"
	"github.com/xuperchain/crypto/core/hdwallet/wordlist"
)

// 助记词强度：弱、中、强
type MnemStrgth int

const (
	_ MnemStrgth = iota
	// 1:弱(12个助记词)
	MnemStrgthWeak
	// 2:中(18个助记词)
	MnemStrgthMedium
	// 3:强(24个助记词)
	MnemStrgthStrong
)

// 助记词语言：英文、中文
type MnemLang int

const (
	_ MnemLang = iota
	// 1:中文
	MnemLangCN
	// 2:英文
	MnemLangEN
)

type Account struct {
	// 钱包地址
	Address string `json:"address,omitempy"`
	// 私钥
	PrivateKey string `json:"private_key,omitempy"`
	// 公钥
	PublicKey string `json:"public_key,omitempy"`
	// 助记词
	Mnemonic string `json:"mnemonic,omitempy"`
}

// 新创建xuperchain ecdsa账户
func NewXchainEcdsaAccount(strg MnemStrgth, lang MnemLang) (*Account, error) {
	cryptoCli := &xchain.XchainCryptoClient{}

	ecdsaAcc, err := cryptoCli.CreateNewAccountWithMnemonic(int(lang), uint8(strg))
	if err != nil {
		return nil, err
	}

	acc := &Account{
		Address:    ecdsaAcc.Address,
		PrivateKey: ecdsaAcc.JsonPrivateKey,
		PublicKey:  ecdsaAcc.JsonPublicKey,
		Mnemonic:   ecdsaAcc.Mnemonic,
	}

	return acc, nil
}

// 根据助记词生成历史账户
func RetrieveAccountByMnemonic(mnemonic string, language int) (*Account, error) {
	cryptoCli := &xchain.XchainCryptoClient{}

	acount, err := cryptoCli.RetrieveAccountByMnemonic(mnemonic, language)
	if err != nil {
		return nil, err
	}

	acc := &Account{
		Address:    acount.Address,
		PrivateKey: acount.JsonPrivateKey,
		PublicKey:  acount.JsonPublicKey,
		Mnemonic:   acount.Mnemonic,
	}

	return acc, nil
}

// InferLanguage 根据助记词推测语言
//  如果推测失败，返回 0
//  这里不做合法性检查，仅推测可能值
func InferLanguage(mnemonic string) int {
	ss := strings.Split(mnemonic, " ")
	if len(ss) == 0 {
		return 0
	}

	if _, exist := wordlist.ReversedEnglishWordMap[ss[0]]; exist {
		return int(MnemLangEN)
	} else if _, exist := wordlist.ReversedSimplifiedChineseWordMap[ss[0]]; exist {
		return int(MnemLangCN)
	}

	return 0
}
