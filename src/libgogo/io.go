// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo I/O functions
//

package libgogo

//
// Writes the specified string with the maximium given length to the file descriptor given
// Write returns the number of characters actually written or 0 in case of an error
// Implemented in assembler (see corresponding .s file)
//
func Write(fd uint64, text string, length uint64) uint64;

//
// Writes the specified string to the (console) standard output
//
func PrintString(msg string) {
    var strLen uint64;
    strLen = StringLength(msg);
    Write(1, msg, strLen); //1 = standard output
}

//
// Prints a single character to the (console) standard output
// Implemented in assembler (see corresponding .s file)
//
func PrintChar(char byte);

//
// Prints the specified number to the (console) standard output
// This function simplified string handling as the conversion is wrapped and has not to be done every time outside the library
//
func PrintNumber(num uint64) {
    var str string;
    str = IntToString(num);
    PrintString(str);
}

//
// Reads a maximum of buffer_size bytes from the file descriptor given to the specified string
// Note that the string has to be at least of size buffer_size. Read returns 0 if an error occured or the end of the file has been reached
// Implemented in assembler (see corresponding .s file)
//
func Read(fd uint64, buffer string, buffer_size uint64) uint64;

//
// Reads a single character from the file descriptor specified
// Note that Read returns 0 if there was an error (p.e. end of file)
// Implemented in assembler (see corresponding .s file)
//
func GetChar(fd uint64) byte;

//
// Opens the file with the name and flags specified and returns a file descriptor or 0 if there was an error
// Implemented in assembler (see corresponding .s file)
//
func FileOpen(filename string, flags uint64) uint64;

//
// Closes the specified file descriptor and returns errno
// Implemented in assembler (see corresponding .s file)
//
func FileClose(fd uint64) uint64;
