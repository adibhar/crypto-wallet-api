package models;

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"strconv"
)

func ProofOfWork(index int, transactions, prevHash, timestamp string, difficulty int) (string, int) {
	nonce := 0;
	var hash string;

	record := strconv.Itoa(index) + timestamp + transactions + prevHash + strconv.Itoa(nonce);
	rawHash := sha256.Sum256([]byte(record));
	hash = hex.EncodeToString(rawHash[:]);

	for !strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
		nonce++;
		record = strconv.Itoa(index) + timestamp + transactions + prevHash + strconv.Itoa(nonce);
		rawHash = sha256.Sum256([]byte(record));
		hash = hex.EncodeToString(rawHash[:]);
	}

	return hash, nonce;
}

