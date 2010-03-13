// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "./_obj/libgogo"

func main() {
  //Library test I
  const test_text string = "Hello world\n";
  const test_len uint64 = 12;
  var test_ret uint64;
  test_ret = libgogo.Write(1, test_text, test_len);
  fmt.Printf("Library test (should write %d, %d and 18): %d, %d and %d\n", test_len, test_len, test_ret, libgogo.StringLength(test_text), libgogo.GoGoStringLength("18 characters long\000"));

  //Library test II
  var fd uint64;
  var ten_char_string string = "##########";
  var read_ret uint64;
  fd = libgogo.FileOpen("libgogo.go", 0);
  read_ret = libgogo.Read(fd, ten_char_string, 10);
  fmt.Printf("Library test (should neither be 0 nor -1): %d\n", fd);
  fmt.Printf("Library test (should write the first 10 characters of libgogo.go and 10): '%s' and %d\n", ten_char_string, read_ret);
  libgogo.FileClose(fd);

  //Library test III
  libgogo.Exit(0);
  fmt.Printf("If you can read this, something is wrong");
}
