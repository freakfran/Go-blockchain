package crypto

import (
	"MyChain/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// Sign 使用私钥对给定数据进行数字签名。
//
// 参数:
// data []byte - 需要签名的数据。
//
// 返回值:
// *Signature - 生成的签名对象。
// error - 如果签名过程中发生错误，则返回错误对象。
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	// 使用ECDSA算法对数据进行签名
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err // 如果签名过程中出现错误，返回错误
	}
	return &Signature{r, s}, nil // 成功签名后，返回签名对象
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{key}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{key: &k.key.PublicKey}
}

type PublicKey struct {
	key *ecdsa.PublicKey
}

// ToSlice 将公钥转换为字节切片。
//
// 参数:
//
//	k PublicKey - 公钥的结构体实例，包含未压缩的密钥信息。
//
// 返回值:
//
//	[]byte - 经过压缩的公钥字节切片。
func (k PublicKey) ToSlice() []byte {
	// 使用椭圆曲线算法将公钥压缩并转换为字节切片
	return elliptic.MarshalCompressed(k.key, k.key.X, k.key.Y)
}

// Address 方法根据公钥计算并返回对应的地址。
// 该方法没有参数。
// 返回值:
//
//	types.Address: 从公钥计算得到的地址。
func (k PublicKey) Address() types.Address {
	// 使用SHA256算法计算公钥的哈希值
	bytes := sha256.Sum256(k.ToSlice())
	// 从哈希值的后20字节创建并返回地址
	return types.AddressFromBytes(bytes[len(bytes)-20:])
}

type Signature struct {
	r, s *big.Int
}

// Verify 使用给定的公钥验证签名是否有效。
//
// 参数:
//
//	pubKey - 公钥，用于验证签名。
//	data - 被签名的数据。
//
// 返回值:
//
//	返回一个布尔值，表示签名是否有效。
func (s Signature) Verify(pubKey PublicKey, data []byte) bool {
	// 使用ECDSA算法验证签名
	return ecdsa.Verify(pubKey.key, data, s.r, s.s)
}
