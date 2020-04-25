package uuid

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"os"
	"strings"
	"time"
)

var defaultNamespace = []byte("uuid")
var defaultHostCache []byte
var baseTimeCache *time.Time

var hs = hex.EncodeToString

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
		defaultHostCache = hash[:2]
	}

	noise := make([]byte, 3)
	_, _ = rand.Read(noise)

	return &UUID{
		Namespace: defaultNamespace,
		Time:      time.Now(),
		Machine:   defaultHostCache,
		Noise:     noise,
	}
}

// Parse a uuid
func Parse(data []byte) *UUID {
	timeBin := append([]byte{0}, data[4:11]...)
	micro := binary.BigEndian.Uint64(timeBin)
	t := (*baseTimeCache).Add(time.Duration(micro) * time.Microsecond)

	return &UUID{data[:4], t, data[11:13], data[13:]}
}

func (id *UUID) String() string {
	return strings.Join([]string{
		string(id.Namespace),
		id.Time.Format("2006_01_02T15:04:05"),
		hs(id.Machine),
		hs(id.Noise),
	}, "-")
}

// Bin of a new uuid. If namespace is nil it will be set to "uuid".
func (id *UUID) Bin() []byte {
	if len(id.Namespace) != 4 {
		panic("[ysmood/uuid] length of namespace must be 4")
	}
	if len(id.Machine) != 2 {
		panic("[ysmood/uuid] length of machine must be 2")
	}
	if len(id.Noise) != 3 {
		panic("[ysmood/uuid] length of noise must be 2")
	}

	if baseTimeCache == nil {
		base, _ := time.Parse("2006", "2020")
		baseTimeCache = &base
	}

	bin := []byte{}

	timeBin := make([]byte, 8)
	micro := id.Time.Sub(*baseTimeCache).Microseconds()
	binary.BigEndian.PutUint64(timeBin, uint64(micro))

	bin = append(bin, id.Namespace...)
	bin = append(bin, timeBin[1:]...)
	bin = append(bin, id.Machine...)
	bin = append(bin, id.Noise...)

	return bin
}

// Hex of a new uuid in hex format. If namespace is empty it will be set to "uuid".
func (id *UUID) Hex() string {
	return hs(id.Bin())
}
