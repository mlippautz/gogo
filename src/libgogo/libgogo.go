// Copyright 2009 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo library functions
//

package libgogo

import "fmt" //TODO (SC): Don't rely on Go libraries

func Min(a uint64, b uint64) uint64 {
    var result uint64 = b;
    if a < b {
        result = a;
    }
    return result;
}

func ByteBufLength(buf [255]byte) uint64 {
    var i uint64;
    for i = 0; buf[i] != 0; i = i +1 {
        }
    return i;
}

func StringLength(str string) uint64;

func StringByteBufCmp(str string, buf [255]byte) uint64 {
    var i uint64;
    var equal uint64 = 1;
    var strlen uint64 = StringLength(str);
    var bufsize uint64 = ByteBufLength(buf);
    var size uint64 = Min(strlen, bufsize);
    if strlen != bufsize {
        equal = 0;
    } else {
        for i = 0; i < size; i = i +1 {
            if str[i] != buf[i] {
                equal = 0;
            }
        }
    }
    return equal;
}

func StringToByteBuf(from string) [255]byte {
  var i uint64;
  var to [255]byte;
  for i = 0; i < StringLength(from) ; i = i+1 {
    to[i] = from[i];
  }
  to[i] = 0;
  return to;
}

func ToIntFromByte(b byte) uint64;

func ToByteFromInt(i uint64) byte;

func ByteBufToInt(byteBuf [255]byte, bufLen uint64) uint64 {
    var i uint64;    
    var val uint64 = 0;
    for i = 0; i < bufLen ; i = i +1 {
        val = val * 10;
        val = val + ToIntFromByte(byteBuf[i]) - 48;
    }
    return val;
}

func PrintByteBuf(buf [255]byte) {
    var i uint64;
    for i = 0; buf[i] != 0; i = i +1 {
        PrintChar(buf[i]);
    }
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

func PrintChar(char byte) { //TODO (SC): Don't rely on Printf
    fmt.Printf("%c", char);
}

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
    }
    for ; i != 0; i = i -1 {
        PrintChar(buf[i]);
    }
    PrintChar(buf[0]);
}

//func GetChar(fd uint64) byte;

//--- Cleanup necessary from here onwards (most functions don't work properly!)

func Read(fd uint64, buffer string, buffer_size uint64) uint64; //Reads the specified number of characters from the file with the given file descriptor to the given buffer (string)
func FileOpen(filename string, flags uint64) uint64; //Opens a file with the specified flags and returns the corresponding file descriptor
func FileClose(fd uint64); //Closes the given file descriptor

func GetChar(fd uint64) byte {
  var one_char_buf string = "\000";
  if Read(fd, one_char_buf, 0) == 0 {
    ; //TODO: Error handling?
    return 0;
  }
  return one_char_buf[0];
}
