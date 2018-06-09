package block

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"strconv"

	"github.com/sirupsen/logrus"
)

// -- Notes
// doing the work is hard, but verifying the proof is easy
// work is to find a hash for a block, hash serves as a proof, so finding a proof is the actual work
// smaller target will mean a smaller upper-bound creating a smaller pool of elgible numbers making it harder to find a valid hash

const targetBits = 24
const maxNonce = math.MaxInt64

// ProofOfWork is the unit of work needed to add a block
type ProofOfWork struct {
	block  *Block   // block we're trying to find a hash for
	target *big.Int // upper-bound on the range of hashes we accept, ie hash we're looking MUST be less than this target value
}

// NewProofOfWork returns a ProofOfWork task for a given block
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{b, target}
}

// prepareData concats all of the data fields in a block AND the nonce (pseudo-random counter to increase) into 1 byte array
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	intToHexHelp := func(n interface{}) []byte {
		var normN int64
		switch v := n.(type) {
		case int:
			normN = int64(v)
		case int64:
			normN = v
		}

		return []byte(strconv.FormatInt(normN, 16))
	}

	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			intToHexHelp(pow.block.Timestamp),
			intToHexHelp(targetBits),
			intToHexHelp(nonce),
		},
		[]byte{},
	)

	return data
}

// run finds a hash for the block that satisfies the upper-bound target
func (pow *ProofOfWork) run() (int, []byte) {
	var hash [32]byte // 32 bytes = 256{32*8} bits
	var hashInt big.Int

	logrus.Printf("Mining the block containing \"%s\"", pow.block.Data)
	var nonce int
	for nonce = 0; nonce < maxNonce; nonce++ {
		// serialize block and nonce data
		data := pow.prepareData(nonce)

		// get sha256 hash
		hash = sha256.Sum256(data)

		// copy the hash value into a BigInt type
		hashInt.SetBytes(hash[:])

		// confirm the hashInt is less than the target
		if hashInt.Cmp(pow.target) == -1 {
			break
		}

		//fmt.Printf("\r%v: %x", nonce, hash)
	}

	return nonce, hash[:]
}

// Validate hashses the nonce and the block-data to verify it meets the target upper-bound and matches the stored hash
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	var storedHashInt big.Int

	// make sure nounce/block hash is under target
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1

	// make sure nounce/block hash matches the one stored int he block
	storedHashInt.SetBytes(pow.block.Hash[:])
	isMatch := hashInt.Cmp(&storedHashInt) == 0

	return isValid && isMatch
}
