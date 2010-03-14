// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "os"
import "fmt"
import "./_obj/libgogo"

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: libgogotest file\n");
		return;	
	}

  //Library test I
  const test_text string = "Hello world\n";
  const test_len uint64 = 12;
  var test_ret uint64;
  test_ret = libgogo.Write(1, test_text, test_len);
  fmt.Printf("Library test (should write %d and %d): %d and %d\n", test_len, test_len, test_ret, libgogo.StringLength(test_text));

  //Library test II
  var fd uint64;
  var ten_char_string string = "##########";
  var read_ret uint64;

  //Does not work yet as FileOpenByteBuf does not work yet
  /*var filename [255]byte;
  filename = libgogo.StringToByteBuf("libgogotest.go"); //Works
  //filename = libgogo.StringToByteBuf(os.Args[1]); //Does not work!?
  fmt.Printf("Filename: '");
  libgogo.WriteByteBuf(1, filename, libgogo.ByteBufLength(filename));
  fmt.Printf("'\n");
  fd = libgogo.FileOpenByteBuf(filename, 0);*/

  fd = libgogo.FileOpen("libgogotest.go", 0); //Works
  //fd = libgogo.FileOpen(os.Args[1], 0); //Does not work!?
  fmt.Printf("Library test for '%s' (should neither be 0 nor -1): %d\n", os.Args[1], fd);
  read_ret = libgogo.Read(fd, ten_char_string, 10);
  var next_byte byte;
  next_byte = libgogo.GetChar(fd);
  fmt.Printf("Library test (should write the first 11 characters of '%s' and 10): '%s%c' and %d\n", os.Args[1], ten_char_string, next_byte, read_ret);
  libgogo.FileClose(fd);

  //Library test III
  var test [255]byte;
  test = libgogo.StringToByteBuf(test_text);
  libgogo.WriteByteBuf(1, test, libgogo.ByteBufLength(test));
  var test_string string;
  test_string = libgogo.ByteBufToString(test);
  fmt.Printf("%s", test_string);

  //Library test IV
  libgogo.Exit(0);
  fmt.Printf("If you can read this, something is wrong");
}
