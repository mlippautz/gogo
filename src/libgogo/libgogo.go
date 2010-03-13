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
