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

func CopyMem(source uint64, dest uint64, size uint64);

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
