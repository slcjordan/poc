package assert

import "testing"

type ErrorChecker interface {
	CheckError(*testing.T, string, error)
}

type BoolChecker interface {
	CheckBool(*testing.T, string, bool)
}

type ByteChecker interface {
	CheckByte(*testing.T, string, byte)
}

type Complex64Checker interface {
	CheckComplex64(*testing.T, string, complex64)
}

type Complex128Checker interface {
	CheckComplex128(*testing.T, string, complex128)
}

type Float32Checker interface {
	CheckFloat32(*testing.T, string, float32)
}

type Float64Checker interface {
	CheckFloat64(*testing.T, string, float64)
}

type IntChecker interface {
	CheckInt(*testing.T, string, int)
}

type Int16Checker interface {
	CheckInt16(*testing.T, string, int16)
}

type Int32Checker interface {
	CheckInt32(*testing.T, string, int32)
}

type Int64Checker interface {
	CheckInt64(*testing.T, string, int64)
}

type Int8Checker interface {
	CheckInt8(*testing.T, string, int8)
}

type RuneChecker interface {
	CheckRune(*testing.T, string, rune)
}

type StringChecker interface {
	CheckString(*testing.T, string, string)
}

type UintChecker interface {
	CheckUint(*testing.T, string, uint)
}

type Uint16Checker interface {
	CheckUint16(*testing.T, string, uint16)
}

type Uint32Checker interface {
	CheckUint32(*testing.T, string, uint32)
}

type Uint64Checker interface {
	CheckUint64(*testing.T, string, uint64)
}

type Uint8Checker interface {
	CheckUint8(*testing.T, string, uint8)
}

type UintptrChecker interface {
	CheckUintptr(*testing.T, string, uintptr)
}
