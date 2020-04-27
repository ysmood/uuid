package uuid

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

const lenNamespace = 3
const lenMachine = 2
const lenTime = 7
const lenNoise = 4

var defaultNamespace = bytes.Repeat([]byte("_"), 3)
var defaultHostCache []byte
var baseTimeCache *time.Time

// UUID holds the time prefix and random suffix
type UUID struct {
	Namespace []byte
	Time      time.Time
	Machine   []byte
	Noise     []byte
}

// New uuid struct with default values
func New() *UUID {
	if defaultHostCache == nil {
		host, _ := os.Hostname()

		hash := sha256.Sum256([]byte(host))
		defaultHostCache = hash[:lenMachine]
	}

	noise := make([]byte, lenNoise)
	_, _ = rand.Read(noise)

	return &UUID{
		Namespace: defaultNamespace,
		Time:      time.Now(),
		Machine:   defaultHostCache,
		Noise:     noise,
	}
}

// N set the namespace
func (id *UUID) N(ns []byte) *UUID {
	id.Namespace = ns
	return id
}

// M set the machine
func (id *UUID) M(m []byte) *UUID {
	id.Machine = m
	return id
}

// Parse a uuid binary
func Parse(data []byte) *UUID {
	machineStart := lenNamespace + lenTime
	machineEnd := machineStart + lenMachine
	timeBin := append(make([]byte, 8-lenTime), data[lenNamespace:machineStart]...)
	micro := binary.BigEndian.Uint64(timeBin)
	t := (*baseTimeCache).Add(time.Duration(micro) * time.Microsecond)

	return &UUID{data[:lenNamespace], t, data[machineStart:machineEnd], data[machineEnd:]}
}

// ParseHex a uuid hex string
func ParseHex(data string) *UUID {
	b, _ := hex.DecodeString(data)
	return Parse(b)
}

func (id *UUID) String() string {
	return strings.Join([]string{
		string(id.Namespace),
		id.Time.Format("2006_01_02T15:04:05"),
		string(id.Machine),
		hex.EncodeToString(id.Noise),
	}, "-")
}

// Bin of a new uuid.
func (id *UUID) Bin() []byte {
	if baseTimeCache == nil {
		base, _ := time.Parse("2006", "2020")
		baseTimeCache = &base
	}

	if len(id.Namespace) != lenNamespace {
		panic(fmt.Sprintf("length of namespace must be %d", lenNamespace))
	}
	if id.Time.Sub(*baseTimeCache) < 0 {
		panic(fmt.Sprintf("time must be greater than %v", *baseTimeCache))
	}
	if len(id.Machine) != lenMachine {
		panic(fmt.Sprintf("length of machine must be %d", lenMachine))
	}
	if len(id.Noise) != lenNoise {
		panic(fmt.Sprintf("length of noise must be %d", lenNoise))
	}

	bin := []byte{}

	timeBin := make([]byte, 8)
	micro := id.Time.Sub(*baseTimeCache).Microseconds()
	binary.BigEndian.PutUint64(timeBin, uint64(micro))

	bin = append(bin, id.Namespace...)
	bin = append(bin, timeBin[8-lenTime:]...)
	bin = append(bin, id.Machine...)
	bin = append(bin, id.Noise...)

	return bin
}

// Hex of a new uuid in hex format. If namespace is empty it will be set to "uuid".
func (id *UUID) Hex() string {
	return hex.EncodeToString(id.Bin())
}
