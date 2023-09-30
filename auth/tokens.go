package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/fatih/structs"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kataras/golog"
)

const (
	tokenTimeout         = 4 * time.Hour
	TokenRememberTimeout = 30 * 24 * time.Hour
)

var (
	UserTokenRedisConn *redis.Client
	UserTokenRedisCtx  context.Context
)

type UserToken struct {
	TokenID   string    `json:"token_id" structs:"-"`
	ID        int64     `json:"id" structs:"id"`
	Email     string    `json:"email" structs:"email"`
	CreatedAt time.Time `json:"created_at" structs:"-"`
}

func (token UserToken) IsNil() bool {
	return token.TokenID == ""
}

func (token UserToken) ToUserData() map[string]interface{} {
	return structs.Map(token)
}

// Store :
// tokenID == "" si une erreur s'est produite
// sinon retourne le tokenID de l'objet 'UserToken' créé.
func (userToken *UserToken) Store(remember bool) (tokenID string) {

	tokenID = uuid.New().String()
	userToken.TokenID = tokenID

	encryptedData, err := marshallAndEncryptUser(*userToken, secretKey)
	if err != nil {
		golog.Error(err)
		return ""
	}
	t := tokenTimeout
	if remember {
		t = TokenRememberTimeout
	}
	err = UserTokenRedisConn.Set(UserTokenRedisCtx, tokenID, encryptedData, t).Err()
	if err != nil {
		golog.Error(err)
		return ""
	}

	return
}

// GetUserToken :
// err != nil si une erreur est présente
// sinon 'userToken' est bien présent et correct
func GetUserToken(tokenID string) (userToken UserToken, err error) {

	encryptedData, err := UserTokenRedisConn.Get(UserTokenRedisCtx, tokenID).Result()
	if err == redis.Nil || encryptedData == "" {
		return
	}

	userToken, err = unmarshallAndDecryptUser(encryptedData, secretKey)
	if err != nil {
		UserTokenRedisConn.Del(UserTokenRedisCtx, tokenID)
	}

	return
}

func RevokeUserToken(tokenID string) (success bool) {
	err := UserTokenRedisConn.Del(UserTokenRedisCtx, tokenID).Err()
	if err == redis.Nil {
		return true
	} else if err != nil {
		golog.Error(err)
		return
	}

	return true
}

// NewEncryptSecretKey:
// Genère une nouvelle clé de chiffrement à chaque démarrage du backend
func NewEncryptSecretKey() string {

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		golog.Fatalf("error occured when creating encrypt key : %v", err)
	}

	return hex.EncodeToString(bytes)
}

// unmarshallAndDecryptUser :
// dechiffre et deserialise une string recupérée de redis
func unmarshallAndDecryptUser(marshalledEncryptedUser, Key string) (unmarshalledUser UserToken, err error) {
	data := Decrypt(marshalledEncryptedUser, Key)
	err = json.Unmarshal([]byte(data), &unmarshalledUser)
	return
}

// marshallAndEncryptUser:
// serialise et chiffre un objet de type 'UserToken' à l'aide d'une 'Key' de chiffrement
func marshallAndEncryptUser(unmarshalledUser UserToken, Key string) (marshalledEncryptedUser string, err error) {
	data, err := json.Marshal(unmarshalledUser)
	if err != nil {
		return
	}

	return Encrypt(string(data), Key), err
}

// Encrypt :
// fonction basique de chiffrement
func Encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		golog.Errorf("failed encrypt with cipher, error : %v", err)
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		golog.Errorf("cipher newGCM with error : %v", err)
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		golog.Errorf("reading aes error : %v", err)
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

// decrypt :
// fonction basique de déchiffrement
func Decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		golog.Errorf("decrypt new cipher error: %v", err)
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		golog.Errorf("decrypt new gcm error : %v", err)
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		decryptedString = ""
	} else {
		decryptedString = string(plaintext)
	}

	return
}
