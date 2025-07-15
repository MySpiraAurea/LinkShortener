package utils

import (
    "crypto/rand"
    "math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortID() string {
    const length = 6
    id := make([]byte, length)
    for i := range id {
        num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        id[i] = charset[num.Int64()]
    }
    return string(id)
}