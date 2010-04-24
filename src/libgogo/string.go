// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo string functions
//

package libgogo

//
// Function returning the length of an ASCII (!) string.
// Parameter the Go string.
// See asm_linux_amd64.s for details
//
func StringLength(str string) uint64;

//
// Function returning the length of an ASCII (!) string.
// Takes the pointer of a Go string.
// See asm_linux_amd64.s for details
//
func StringLength2(str *string) uint64;

func GetStringAddress(str *string) uint64;

func GetStringFromAddress(addr uint64) *string;

//
// Simple string compare function.
// Returns 0 if strings are equal, 1 otherwise.
//
func StringCompare(str1 string, str2 string) uint64 {
    var i uint64;
    var equal uint64 = 0;
    var strlen1 uint64 = StringLength(str1);
    var strlen2 uint64 = StringLength(str2);
    if strlen1 != strlen2 {
       equal = 1;
    } else {
        for i = 0; i < strlen1; i = i +1 {
            if str1[i] != str2[i] {
                equal = 1;
            }
        }
    }
    return equal;
}

func SetStringAddressAndLength(str *string, new_addr uint64, new_length uint64);

//
// Function appending a single character to a string.
// Basically moving a new copy of the string with the additional character 
// appended to a new place in the heap.
//
func CharAppend(str *string, char byte) {
    var nullByte byte = 0;
    var strlen uint64 = StringLength2(str);
    var new_length uint64 = strlen + 1;
    var new_addr uint64 = Alloc(new_length+1);
    var old_addr uint64 = GetStringAddress(str);
    CopyMem(old_addr, new_addr, strlen);
    CopyMem(ToUint64FromBytePtr(&char), new_addr + strlen, 1);
    CopyMem(ToUint64FromBytePtr(&nullByte), new_addr+strlen +1, 1);
    SetStringAddressAndLength(str, new_addr, new_length);
}

//
// Function appending a whole string to a given string.
// Moving both strings to a new allocated place in the heap.
//
func StringAppend(str *string, append_str string) {
    var nullByte byte = 0;
    var strlen uint64 = StringLength2(str);
    var strappendlen uint64 = StringLength(append_str);
    var new_length uint64 = strlen + strappendlen;
    var new_addr uint64 = Alloc(new_length+1);
    var old_addr uint64 = GetStringAddress(str);
    var append_addr uint64 = GetStringAddress(&append_str);
    CopyMem(old_addr, new_addr, strlen);
    CopyMem(append_addr, new_addr + strlen, strappendlen);
    CopyMem(ToUint64FromBytePtr(&nullByte), new_addr + strlen + strappendlen + 1, 1);
    SetStringAddressAndLength(str, new_addr, new_length);
}
