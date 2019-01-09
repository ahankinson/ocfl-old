package libocfl

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"sync"
)

type ByteSize int64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

type HashStruct struct {
	Path		string
	Hash		string
	HashType	string
}

/*
	Takes a file path and an algorithm and provides the
	base16-encoded string of the hash of the file. First
	checks the size of the file and, if it's under 1MB,
	passes it off to a single-threaded hashing function;
	otherwise, passes it off to a hasher that is capable
	of concurrent processing of the file.
*/
func HashFile(path string, algorithm string) *HashStruct {
	file, oerr := os.Open(path)

	if oerr != nil {
		panic(oerr.Error())
	}

	info, _ := file.Stat()
	filesize := info.Size()
	cerr := file.Close()

	if cerr != nil {
		panic(cerr.Error())
	}

	var hsh *HashStruct

	if filesize <= int64(MB) {
		hsh = singleHashing(path, algorithm)
	} else {
		hsh = concurrentHashing(path, algorithm)
	}

	return hsh
}


func chooseHashAlg(algorithm string) hash.Hash {
	switch algorithm {
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha256":
		return sha256.New()
	case "sha512":
		return sha512.New()
	case "blake2b":
		return crypto.BLAKE2b_512.New()
	default:
		return sha512.New()
	}
}

type chunk struct {
	bufsize int64
	offset  int64
}

func concurrentHashing(path string, algorithm string) *HashStruct {
	res := &HashStruct{
		Path: path,
		HashType: algorithm,
	}

	file, _ := os.Open(path)
	info, _ := file.Stat()
	filesize := info.Size()

	concurrency := NumWorkers
	bufsize := int64(8 * KB)

	chunksizes := make([]chunk, concurrency)

	for i := 0; i < NumWorkers; i++ {
		chunksizes[i].bufsize = bufsize
		chunksizes[i].offset = bufsize * int64(i)
	}

	if remainder := filesize % bufsize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency) * bufsize}
		concurrency++
		chunksizes = append(chunksizes, c)
	}

	var wg sync.WaitGroup
	byteChan := make(chan []byte)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func(chunksizes []chunk, i int) {
			defer wg.Done()

			chunk := chunksizes[i]
			buffer := make([]byte, chunk.bufsize)
			_, err := file.ReadAt(buffer, chunk.offset)

			if err != nil {
				fmt.Println(err)
				return
			}

			byteChan <- buffer


		}(chunksizes, i)
	}

	hsh := chooseHashAlg(algorithm)

	for chk := range byteChan {
		hsh.Write(chk)
	}

	h := hex.EncodeToString(hsh.Sum(nil))
	close(byteChan)

	_ := file.Close()

	res.Hash = h
	return res
}

func singleHashing(path string, algorithm string) *HashStruct {
	res := &HashStruct{
		Path: path,
		HashType: algorithm,
	}

	hsh := chooseHashAlg(algorithm)
	file, _ := os.Open(path)
	defer file.Close()

	buf := make([]byte, 1024)

	for {
		bytesread, err := file.Read(buf)

		if err != nil {
			hsh.Write(buf[:bytesread])
			break
		}

		hsh.Write(buf[:bytesread])
	}

	h := hex.EncodeToString(hsh.Sum(nil))
	res.Hash = h

	return res
}
