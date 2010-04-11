// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "os"
import "./libgogo/_obj/libgogo"

var filename string;

func main() {
    var fd uint64;
    var errno uint64;
    var doRest uint64 = 1;

    if len(os.Args) != 2 {
        libgogo.PrintString("Usage: gogo file.go\n");
        doRest = 0;
    }
    
    if doRest != 0 {

        filename = os.Args[1];

        fd = libgogo.FileOpen(filename, 0);
        if fd != 0 {
            //ScannerTest(fd);
            Parse(fd);
            errno = libgogo.FileClose(fd);
            if errno != 0 {
                libgogo.ExitError("Error closing file", errno);
            }
        } else {
            libgogo.PrintString("Error opening file ");
            libgogo.PrintString(filename);
            libgogo.PrintString(".\n");
        }
    }
}
