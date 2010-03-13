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
  fmt.Printf("Library test (should write %d and 18): %d and %d\n", test_len, test_ret, libgogo.StringLength("18 characters long\000"));

  //Library test II
  var one_char_string string = "_";
  var read_ret uint64;
  fmt.Printf("Library test: please type one character and hit return\n");
  read_ret = libgogo.Read(0, one_char_string, 1);
  fmt.Printf("Library test (should write the typed character and 1): %s and %d\n", one_char_string, read_ret);

  //Library test III
  libgogo.Exit(0);
  fmt.Printf("If you can read this, something is wrong");
}
