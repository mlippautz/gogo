// Copyright 2009 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo library functions
//

package libgogo

var Argv [255]string;
var Argc uint64 = 0;

func GetArgv() {
    var fd uint64;
    var errno uint64;    
    var char string = "#";
    var inArgv uint64 = 0;

    fd = FileOpen("/proc/self/cmdline", 0);
    if fd == 0 {
        ExitError("Error opening /proc/self/cmdline. Currently GoGo is only supported on systems with /proc enabled.", 1);
    }

    for errno = Read(fd, char, 1) ; errno != 0 ; errno = Read(fd, char, 1) {
        if char[0] == 0 {
            inArgv = 1;
            Argc = Argc + 1;
        } else {
            if inArgv == 1 {
                Argv[Argc] += string(char[0]); // (SC) TODO: Remove cast, str append
            }
        }
    }

    errno = FileClose(fd);
    if errno != 0 {
        ExitError("Error closing file /proc/self/cmdline",1);
    }
}

func StringAppend(str *string, char byte) { //TODO (SC): Get rid of magic string concatenation (needs memory management!) and cast operator
  *str += string(char);
}

func Min(a uint64, b uint64) uint64 {
    var result uint64 = b;
    if a < b {
        result = a;
    }
    return result;
}

func StringLength(str string) uint64;

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

func ToIntFromByte(b byte) uint64;

func ToByteFromInt(i uint64) byte;

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

func Exit(code uint64);

func ExitError(msg string, code uint64) {
    PrintString(msg);
    PrintChar('\n');
    Exit(code);
}

func Write(fd uint64, text string, length uint64) uint64;

func PrintString(msg string) {
    Write(1, msg, StringLength(msg));
}

func PrintChar(char byte);

func PrintNumber(num uint64) {
    var i uint64;
    var buf [255]byte;
    for i = 0; num != 0; i = i +1 {
        buf[i] = ToByteFromInt(num - (num / 10) * 10 + 48);
        num = num / 10;
    }
    if i == 0 { //Special case: 0
        buf[0] = 48;
        i = 1;
    } else {
        i = i -1;
    }
    for ; i != 0; i = i -1 {
        PrintChar(buf[i]);
    }
    PrintChar(buf[0]);
}

func Read(fd uint64, buffer string, buffer_size uint64) uint64;

func GetChar(fd uint64) byte;

func FileOpen(filename string, flags uint64) uint64;

func FileClose(fd uint64) uint64;
