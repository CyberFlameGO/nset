package nlookup

import (
	"fmt"
	"strings"
)

var _ fmt.Stringer = &NLookup[uint]{}

type StorageType uint64

const StorageTypeBits = 64

type IntsIf interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type NLookup[T IntsIf] struct {
	Data []StorageType
}

func (n *NLookup[T]) Add(x T) {

	unitIndex := n.GetStorageUnitIndex(x)
	if unitIndex >= n.Size() {
		storageUnitsToAdd := unitIndex - n.Size() + 1
		n.Data = append(n.Data, make([]StorageType, storageUnitsToAdd)...)
	}

	n.Data[unitIndex] |= 1 << (x % StorageTypeBits)
}

func (n *NLookup[T]) Remove(x T) {

	unitIndex := n.GetStorageUnitIndex(x)
	if unitIndex >= n.Size() {
		return
	}

	n.Data[unitIndex] ^= 1 << (x % StorageTypeBits)
}

func (n *NLookup[T]) Contains(x T) bool {
	return n.isSet(x)
}

func (n *NLookup[T]) ContainsAny(values ...T) bool {

	for _, x := range values {
		if n.isSet(x) {
			return true
		}
	}

	return false
}

func (n *NLookup[T]) ContainsAll(values ...T) bool {

	for _, x := range values {
		if !n.isSet(x) {
			return false
		}
	}

	return true
}

func (n *NLookup[T]) isSet(x T) bool {
	unitIndex := n.GetStorageUnitIndex(x)
	return unitIndex < n.Size() && n.Data[unitIndex]&(1<<(x%StorageTypeBits)) != 0
}

func (n *NLookup[T]) GetStorageUnitIndex(x T) uint64 {
	return uint64(x) / StorageTypeBits
}

func (n *NLookup[T]) GetStorageUnit(x T) StorageType {
	return n.Data[x/StorageTypeBits]
}

//Size returns len(n.Data)
func (n *NLookup[T]) Size() uint64 {
	return uint64(len(n.Data))
}

func (n *NLookup[T]) ElementCap() uint64 {
	return uint64(len(n.Data) * StorageTypeBits)
}

//String returns a string of the storage as bytes separated by spaces. A comma is between each storage unit
func (n *NLookup[T]) String() string {

	b := strings.Builder{}
	b.Grow(len(n.Data)*StorageTypeBits + len(n.Data)*2)

	for i := 0; i < len(n.Data); i++ {

		x := n.Data[i]
		shiftAmount := StorageTypeBits - 8
		for shiftAmount >= 0 {

			byteToShow := uint8(x >> shiftAmount)
			if shiftAmount > 0 {
				b.WriteString(fmt.Sprintf("%08b ", byteToShow))
			} else {
				b.WriteString(fmt.Sprintf("%08b", byteToShow))
			}

			shiftAmount -= 8
		}
		b.WriteString(", ")
	}

	return b.String()
}

func NewNLookup[T IntsIf]() NLookup[T] {

	return NLookup[T]{
		Data: make([]StorageType, 1),
	}
}

//NewNLookupWithMax creates a nlookup that already has capacity to hold till at least largestNum without resizing.
//Note that this is NOT the count of elements you want to store, instead you input the largest value you want to store. You can store larger values as well.
func NewNLookupWithMax[T IntsIf](largestNum T) NLookup[T] {
	return NLookup[T]{
		Data: make([]StorageType, largestNum/StorageTypeBits+1),
	}
}