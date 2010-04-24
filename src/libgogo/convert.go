// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo conversion functions
//

package libgogo

//
// Converts a byte value to an integer
// See asm_linux_amd64.s for details
//
func ToIntFromByte(b byte) uint64;

//
// Converts an unsigned 64bit integer to a byte
// See asm_linux_amd64.s for details
//
func ToByteFromInt(i uint64) byte;

func ToUint64FromUint64Ptr(value *uint64) uint64;

//
// Returns the address of the byte as uint64.
// See asm_linux_amd64.s for details
//
func ToUint64FromBytePtr(char *byte) uint64;

//
// Converter returning the integer (uint64) representation of a given string.
//
func StringToInt(str string) uint64 {
    var n uint64 = StringLength(str);
    var i uint64;
    var val uint64 = 0;
    for i = 0; i < n ; i = i +1 {
        val = val * 10;
        val = val + ToIntFromByte(str[i]) - 48;
    }
    return val;
}

//
// Converter returning a string representation (heap) of a given number.
//
func IntToString(num uint64) string {
    var str string = "";
    var i uint64;
    var buf [255]byte;
    for i = 0; num != 0; i = i +1 {
        buf[i] = ToByteFromInt( num - (num/10) * 10 + 48 );
        num = num / 10;
    }
    if i == 0 { //Special case: 0
        buf[0] = 48;
        i = 0;
    } else {
        i = i -1;
    }
    for ; i != 0; i = i -1 {
        CharAppend(&str,buf[i]);
    }
    CharAppend(&str,buf[0]);
    return str;
}

