// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo conversion functions
//

package libgogo

//
// Converts a byte to an unsigned 64-bit integer
// Implemented in assembler (see corresponding .s file)
//
func ToIntFromByte(b byte) uint64;

//
// Converts an unsigned 64-bit integer to a byte
// Implemented in assembler (see corresponding .s file)
//
func ToByteFromInt(i uint64) byte;

//
// Converts a pointer to an unsigned 64-bit integer
// Implemented in assembler (see corresponding .s file)
//
func ToUint64FromBytePtr(char *byte) uint64;

//
// Converts a pointer to an unsigned 64-bit integer
// Implemented in assembler (see corresponding .s file)
//
func ToUint64FromUint64Ptr(value *uint64) uint64;

//
// Converts an unsigned 64-bit integer to a pointer
// Implemented in assembler (see corresponding .s file)
//
func ToUint64PtrFromUint64(value uint64) *uint64;

//
// Converts a string pointer to an unsigned 64-bit integer
// Implemented in assembler (see corresponding .s file)
//
func ToUint64FromStringPtr(value *string) uint64;

//
// Interprets the given address as a string pointer and returns it
// Implemented in assembler (see corresponding .s file)
//
func GetStringFromAddress(addr uint64) *string;

//
// Converts a string to a string to an unsigned 64-bit integer by intepreting it as valid a decimal number in form of ASCII characters
//
func StringToInt(str string) uint64 {
    var n uint64;
    var i uint64;
    var temp uint64;
    var val uint64 = 0;
    n = StringLength(str);
    for i = 0; i < n ; i = i + 1 { //Process digit by digit
        val = val * 10; //Next digit => Move old value one digit to the left
        temp = ToIntFromByte(str[i]);
        val = val + temp - 48; //Add the new digit as the last (rightmost) one (ASCII 48 = '0', 49 = '1' etc.)
    }
    return val;
}

//
// Converts an unsigned 64-bit integer to a string denoting its decimal digits
//
func IntToString(num uint64) string {
    var str string = "";
    var i uint64;
    var buf [20]byte; //The decimal representation of the longest unsigned 64-bit integer possible is 20 digits long
    for i = 0; num != 0; i = i + 1 { //Process digits one by one in reverse order
        buf[i] = ToByteFromInt(num - (num / 10) * 10 + 48); //Get first (leftmost) digit (ASCII '0' = 48, '1' = 49 etc.)
        num = num / 10; //Next digit => Move value one digit to the right
    }
    if i == 0 { //Special case: 0 (no digits processed)
        buf[0] = 48; //Put ASCII '0' into buffer
        i = 0;
    } else {
        i = i - 1; //Decrement i due to the last loop increment
    }
    for ; i != 0; i = i - 1 { //Reverse the digit order (for all digits but the last one)
        CharAppend(&str, buf[i]);
    }
    CharAppend(&str,buf[0]); //Get the last digit (has to be appended if the value is 0)
    return str;
}

