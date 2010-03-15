// Copyright 2009 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo library functions
//

package libgogo

import "os"
import "fmt"

func ToIntFromByte(b byte) uint64 {
    return uint64(b);
}

func ByteBufToInt(byteBuf [255]byte, bufLen uint64) uint64 {
    var i uint64;    
    var val uint64;
    
    val = 0;

    for i = 0; i < bufLen ; i = i +1 {
        val = val * 10;
        val = val + uint64(byteBuf[i]) - 48;
    }

    return val;
}

func ExitError(msg string, code uint64) {
    fmt.Printf("%s\n",msg);
    os.Exit(int(code));
}

func PrintString(msg string) {
    fmt.Printf(msg);
}

func PrintChar(char byte) {
    fmt.Printf("%c",char);
}

func PrintNumber(num uint64) {
    fmt.Printf("%d",num);
}

func PrintByteBuf(buf [255]byte) {
    var i uint64;
    for i = 0; buf[i] != 0; i = i+1 {
        fmt.Printf("%c",buf[i]);
    }
}

//--- Cleanup necessary from here onwards (most functions don't work properly!)

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
