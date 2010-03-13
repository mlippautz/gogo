// Copyright 2009 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo library functions
//

package libgogo

func Exit(return_code uint64); //Exits the program
func Write(fd uint64, text string, length uint64) uint64; //Writes a defined number of characters of a given string to the file with the given file descriptor
func Read(fd uint64, buffer string, buffer_size uint64) uint64; //Reads the specified number of characters from the file with the given file descriptor to the given buffer (string)
func FileOpen(filename string, flags uint64) uint64; //Opens a file with the specified flags and returns the corresponding file descriptor
func FileClose(fd uint64); //Closes the given file descriptor

func StringCopy(source string, destination string, count uint64); //Copies a maximum of count characters of the source string into the destination string
func StringLength(str string) uint64; //Determines the length of a string

func GoGoStringLength(str string) uint64 { //Determines the length of a GoGo string (not to be used for Go strings!)
  var len uint64;
  for len = 0; str[len] != 0; len = len + 1 {
    }
  return len;
}

func StringToByteBuf(from string, to [255]byte) {
    var i int;
    for i = 0; i < StringLength(from) ; i = i+1 {
        to[i] = from[i];
    }
    to[i] = 0;
    return to;
}
