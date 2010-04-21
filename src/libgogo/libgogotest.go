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
  fmt.Printf("Library test (should write Hello world): ");
  test_ret = libgogo.Write(1, test_text, test_len);
  fmt.Printf("Library test (should write %d and %d): %d and %d\n", test_len, test_len, test_ret, libgogo.StringLength(test_text));

  //Library test II
  var fd uint64;
  var ten_char_string string = "##########";
  var read_ret uint64;
  fd = libgogo.FileOpen(os.Args[1], 0);
  fmt.Printf("Library test for '%s' (should neither be 0 nor -1): %d\n", os.Args[1], fd);
  read_ret = libgogo.Read(fd, ten_char_string, 10);
  var next_byte byte;
  next_byte = libgogo.GetChar(fd);
  fmt.Printf("Library test (should write the first 11 characters of '%s' and 10): '%s%c' and %d\n", os.Args[1], ten_char_string, next_byte, read_ret);
  libgogo.FileClose(fd);

  //Library test III
  var testString string = "Hell";
  testString += "o";
  fmt.Printf("Library test: string length of '%s' (len vs. myLen): %d vs. %d\n", testString, len(testString), libgogo.StringLength(testString));
  var testStringPtr *string = &testString;
  fmt.Printf("Library test 2: string length 2 of '%s' (len vs. myLen): %d vs. %d\n", *testStringPtr, len(*testStringPtr), libgogo.StringLength2(testStringPtr));
  fmt.Printf("Library test: '%s' and '", testString);
  libgogo.PrintString(testString);
  fmt.Printf("' are identical\n");

  //Library test IV
  var chr byte = '\n';
  for fd = libgogo.FileOpen(os.Args[1], 0); chr != 0; chr = libgogo.GetChar(fd) {
    //libgogo.PrintChar(chr);
  }
  libgogo.FileClose(fd);

  //Library test V
  var oldsize uint64 = libgogo.GetBrk();
  if oldsize != 0 {
    fmt.Printf("Brk returned: %d\n", oldsize);
    var newsize uint64 = oldsize + 100 * 1024 * 1024;
    var errno uint64 = libgogo.Brk(newsize);
    if errno == 0 {
      fmt.Printf("Brk successfully allocated 100 MB\n");
      var newreadsize uint64 = libgogo.GetBrk();
      if newreadsize != 0 {
        fmt.Printf("Brk returned: %d, so %d MB of space have been allocated\n", newreadsize, (newreadsize - oldsize) >> 20);
        newreadsize = libgogo.TestMem(oldsize + 1);
        if newreadsize == 0 {
          fmt.Printf("Successfully wrote to address %d\n", oldsize + 1);
        } else {
          fmt.Printf("Failed to write to address %d\n", oldsize + 1);
        }
      }
    } else {
      fmt.Printf("Brk failed to allocate 100 MB: errno %d\n", errno);
    }
  } else {
    fmt.Printf("Brk failed\n");
  }

  //Library test VI
  libgogo.Exit(0);
  fmt.Printf("If you can read this, something is wrong");
}
