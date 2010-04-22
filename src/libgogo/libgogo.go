// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo library functions
//

package libgogo

//
// Argc/Argv variables.
// Are available AFTER they have been set by libgogo.GetArgv()
//
var Argv [255]string;
var Argc uint64 = 0;

//
// Function setting the globally available variables libgogo.Argc and libgogo.Argv.
// This function uses the Linux proc fs to determine its command line!
//
func GetArgv() {
    var fd uint64;
    var errno uint64;    
    var char string = "#";

    fd = FileOpen("/proc/self/cmdline", 0);
    if fd == 0 {
        ExitError("Error opening /proc/self/cmdline. Currently GoGo is only supported on systems with /proc enabled.", 1);
    }

    for errno = Read(fd, char, 1) ; errno != 0 ; errno = Read(fd, char, 1) {
        if char[0] == 0 {
            Argc = Argc + 1;            
        } else {
            CharAppend(&Argv[Argc], char[0]);
        }
    }

    errno = FileClose(fd);
    if errno != 0 {
        ExitError("Error closing file /proc/self/cmdline",1);
    }
}

//
// Simple minimum function
//
func Min(a uint64, b uint64) uint64 {
    var result uint64 = b;
    if a < b {
        result = a;
    }
    return result;
}

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

func GetStringAddress(str *string) uint64;

func GetStringFromAddress(addr uint64) *string;

func CopyMem(source uint64, dest uint64, size uint64);

//
// Returns the address of the byte as uint64.
// See asm_linux_amd64.s for details
//
func ToUint64FromBytePtr(char *byte) uint64;

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
    CopyMem(ToUint64FromBytePtr(&nullByte), new_addr+strlen +1, 1);
    SetStringAddressAndLength(str, new_addr, new_length);
}

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

//
// Exit the current program. (syscall)
// Takes the error number as parameter.
// See asm_linux_amd64.s for details
//
func Exit(code uint64);

//
// Wrapper printing a given message and exiting the program with an error number
//
func ExitError(msg string, code uint64) {
    PrintString(msg);
    PrintChar(10);
    Exit(code);
}

func Write(fd uint64, text string, length uint64) uint64;

func PrintString(msg string) {
    Write(1, msg, StringLength(msg));
}

func PrintChar(char byte);

func PrintNumber(num uint64) {
    PrintString(IntToString(num));
}

func Read(fd uint64, buffer string, buffer_size uint64) uint64;

func GetChar(fd uint64) byte;

func FileOpen(filename string, flags uint64) uint64;

func FileClose(fd uint64) uint64;
