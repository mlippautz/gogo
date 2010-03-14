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

func StringLength(str string) uint64; //Determines the length of a string

func StringToByteBuf(from string) [255]byte {
  var i uint64;
  var to [255]byte;
  for i = 0; i < StringLength(from) ; i = i+1 {
    to[i] = from[i];
  }
  to[i] = 0;
  return to;
}

//Alternative version of ByteBufToString in pure Go, although with string allocation
/*func ByteBufToString(from [255]byte) string {
  var i uint64;
  var to string = "";
  for i = 0; from[i] != 0 ; i = i+1 {
    to += string(from[i]);
  }
  to += "\000";
  return to;
}*/

func InternalByteBufToString(from *byte, to string);

func ByteBufToString(from [255]byte) string {
  var to string = "\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000";
  InternalByteBufToString(&from[0], to);
  return to;
}

func InternalByteBufLength(buf *byte) uint64;

func ByteBufLength(buf [255]byte) uint64 {
  return InternalByteBufLength(&buf[0]);
}

func WriteByteBuf(fd uint64, buf [255]byte, buflen uint64) {
  Write(fd, ByteBufToString(buf), buflen);
}

func GetChar(fd uint64) byte {
  var one_char_buf string = "\000";
  if Read(fd, one_char_buf, 0) == 0 {
    ; //TODO: Error handling?
    return 0;
  }
  return one_char_buf[0];
}

//Does not work!?
func FileOpenByteBuf(filename [255]byte, flags uint64) uint64 {
  return FileOpen(ByteBufToString(filename), flags);
}
