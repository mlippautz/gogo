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
  libgogo.PrintByteBuf(test);
  fmt.Printf("Equality test: %d\n", libgogo.StringByteBufCmp(test_text, test));

  //Library test IV
  var testString string = "Hell";
  testString += "o";
  fmt.Printf("Test: string length of '%s' (len vs. myLen): %d vs. %d\n", testString, len(testString), libgogo.StringLength(testString));
  fmt.Printf("Write verification (currently failing): '%s' and '", testString);
  libgogo.PrintString(testString);
  fmt.Printf("' are identical\n");
  var buf1 [255]byte;
  buf1 = libgogo.StringToByteBuf(os.Args[1]);
  fmt.Printf("Length of '%s': %d\n", os.Args[1], libgogo.ByteBufLength(buf1));

  //Library test V
  libgogo.Exit(0);
  fmt.Printf("If you can read this, something is wrong");
}
