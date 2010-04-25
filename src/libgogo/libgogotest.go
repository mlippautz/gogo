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
  var a uint64 = 0;
  var someString string = libgogo.IntToString(a);
  fmt.Printf("Library test: %s (should be '0')\n", someString);

  //Library test VII
  var test_str string = "Hello ";
  var test_appendstr string = "world";
  libgogo.StringAppend(&test_str, test_appendstr);
  libgogo.PrintString(test_str);
  libgogo.PrintString("\n");

  //Library test VIII
  var stack libgogo.Stack;
  libgogo.InitializeStack(&stack);
  fmt.Printf("Stack initialized: item count: %d, capacity: %d\n", libgogo.GetStackItemCount(&stack), libgogo.GetStackCapacity(&stack));
  var val uint64;
  fmt.Printf("Pushing the 20 values onto the stack\n");
  for val = 1; val <= 20; val++ {
    libgogo.Push(&stack, val);
    //fmt.Printf("DEBUG %d: item count: %d, capacity: %d, peek item: %d\n", val, libgogo.GetStackItemCount(&stack), libgogo.GetStackCapacity(&stack), libgogo.Peek(&stack));
  }
  fmt.Printf("Taking 19 values from the stack\n");
  var j uint64;
  for j = 20; j > 1; j-- {
    val = libgogo.Pop(&stack);
    //fmt.Printf("DEBUG %d: item count: %d, capacity: %d, peek item: %d\n", j, libgogo.GetStackItemCount(&stack), libgogo.GetStackCapacity(&stack), libgogo.Peek(&stack));
    if (j != val) {
      fmt.Printf("ERROR: Unexpected value from stack: %d (should be %d)\n", val, j);
    }
  }
  fmt.Printf("Library test: Number of values left on stack and peek value (should both be 1): %d and %d\n", libgogo.GetStackItemCount(&stack), libgogo.Peek(&stack));

  //Library test IX
  var list libgogo.List;
  libgogo.InitializeList(&list);
  fmt.Printf("List initialized: item count: %d, capacity: %d\n", libgogo.GetListItemCount(&list), libgogo.GetListCapacity(&list));
  fmt.Printf("Putting the 20 values into the list\n");
  for val = 1; val <= 20; val++ {
    libgogo.AddItem(&list, val);
    //fmt.Printf("DEBUG %d: item count: %d, capacity: %d, first item: %d, last item: %d\n", val, libgogo.GetListItemCount(&list), libgogo.GetListCapacity(&list), libgogo.GetItemAt(&list, 0), libgogo.GetItemAt(&list, libgogo.GetListItemCount(&list) - 1));
  }
  fmt.Printf("Removing 19 values from the list\n");
  for j = 1; j < 20; j++ {
    val = libgogo.RemoveItemAt(&list, 0);
    //fmt.Printf("DEBUG %d: item count: %d, capacity: %d, first item: %d, last item: %d\n", j, libgogo.GetListItemCount(&list), libgogo.GetListCapacity(&list), libgogo.GetItemAt(&list, 0), libgogo.GetItemAt(&list, libgogo.GetListItemCount(&list) - 1));
    if (j != val) {
      fmt.Printf("ERROR: Unexpected value from list: %d (should be %d)\n", val, j);
    }
  }
  fmt.Printf("Library test: Number of values left in list and value (should be 1 and 20): %d and %d\n", libgogo.GetListItemCount(&list), libgogo.GetItemAt(&list, 0));

  //Library test X
  libgogo.Exit(0);
  fmt.Printf("If you can read this, something is wrong");
}
