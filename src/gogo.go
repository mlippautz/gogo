// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "os"
import "unsafe"
import "syscall"
import "./libgogo/_obj/libgogo"

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: gogo file.go\n");
		return;	
	}
    
	var r0 uintptr;
	var e1 uintptr;
    var e int;
    var fd uint64;
  var errno uint64;
	r0, _, e1 = syscall.Syscall(syscall.SYS_OPEN, uintptr(unsafe.Pointer(syscall.StringBytePtr(os.Args[1]))), 0, 0);
	e = int(e1);
	fd = uint64(r0);
    
    if e == 0 {
        //ScannerTest(fd);
        Parse(fd);
        errno = libgogo.FileClose(fd);
        if errno != 0 {
            libgogo.ExitError("Error closing file", errno);
        }
	} else {
		fmt.Printf("Error opening file %s.\n", os.Args[1]);
	}
}
