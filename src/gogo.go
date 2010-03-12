// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "os"
import "unsafe"
import "syscall"
import "./libgogo/_obj/libgogo"

func get_next_byte(fd int) byte {
	var b [1]byte;
	_, _, _ = syscall.Syscall(syscall.SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(&b)), 1);
	return b[0];
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: gogo file.go\n");
		return;	
	}

    // just a simple library test calling assembler code
    libgogo.Test();
    
	var r0 uintptr;
	var e1 uintptr;
	r0, _, e1 = syscall.Syscall(syscall.SYS_OPEN, uintptr(unsafe.Pointer(syscall.StringBytePtr(os.Args[1]))), 0, 0);
	var e int = int(e1);
	var fd int = int(r0);
	if e == 0 {
		fmt.Printf("%c\n", get_next_byte(fd) );
		syscall.Syscall(syscall.SYS_CLOSE, uintptr(fd), 0, 0);
	} else {
		fmt.Printf("Error opening file %s.\n", os.Args[1]);
	}
}
