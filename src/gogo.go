// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "fmt"
import "os"
//import "unsafe"
//import "syscall"
import "./libgogo/_obj/libgogo"

func gogo_strcpy (s string) [255]byte {
    var b [255]byte;
    var i int;
    for i=0;i<len(s);i++ {
        b[i] = s[i]
    }
    b[i] = 0;
    return b;
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: gogo file.go\n");
		return;	
	}
    
//	var r0 uintptr;
//	var e1 uintptr;
//	r0, _, e1 = syscall.Syscall(syscall.SYS_OPEN, uintptr(unsafe.Pointer(syscall.StringBytePtr(os.Args[1]))), 0, 0);
//	var e int = int(e1);
//	var fd int = int(r0);
    var fd uint64;
    var test [255]byte;
    libgogo.StringToByteBuf(os.Args[1],test);
    fmt.Printf("%s\n",test);
    fd = libgogo.FileOpen("gogo.go",0);
    fmt.Printf("%d\n",fd);    
	//if e == 0 {
    //scanner_test(fd);		
    libgogo.FileClose(fd);
      //  syscall.Syscall(syscall.SYS_CLOSE, uintptr(fd), 0, 0);
	//} else {
	//	fmt.Printf("Error opening file %s.\n", os.Args[1]);
	//}
}
