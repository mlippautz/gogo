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
// Sets the program's arguments as globally available variables libgogo.Argc and libgogo.Argv
// This function uses the Linux proc fs to determine its command line arguments. It is not possible to get these arguments from the stack as the Go runtime also uses the latter. In order to be compatible with the Go runtime, the arguments are therefore not read from the stack
//
func GetArgv() {
    var fd uint64;
    var errno uint64;    
    var singleChar byte;
    var lastChar byte = 1; // needs to be != 0 at start

    fd = FileOpen("/proc/self/cmdline", 0); //Open file that contains the program's arguments
    if fd == 0 { //Error check (the system may have been compiled with proc fs disabled)
        ExitError("Error opening /proc/self/cmdline. Currently GoGo is only supported on systems with /proc enabled.", 1);
    }    

    for singleChar = GetChar(fd);(singleChar != 0) || (lastChar != 0); singleChar = GetChar(fd) {
        if (singleChar == 0) {
            Argc = Argc +1;
        } else {
            CharAppend(&Argv[Argc], singleChar);
        }
        lastChar = singleChar;
    }

    errno = FileClose(fd);
    if errno != 0 {
        ExitError("Error closing file /proc/self/cmdline",1);
    }
}

//
// Copies size bytes of memory from one location to another
// Implemented in assembler (see corresponding .s file)
//
func CopyMem(source uint64, dest uint64, size uint64);

//
// Exits the current program with an error number (return value) as parameter
// Implemented in assembler (see corresponding .s file)
//
func Exit(code uint64);

//
// Wrapper printing a given message and exiting the program with an error number
//
func ExitError(msg string, code uint64) {
    PrintString(msg);
    PrintChar(10); //Print new line ('\n' = 10)
    Exit(code);
}
