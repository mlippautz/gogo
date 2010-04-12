// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "os"
import "./libgogo/_obj/libgogo"

type FileInfo struct {
    filename string;
    lineCounter uint64;
    charCounter uint64;
    fd uint64;
}; 

var fileInfo [255]FileInfo;
var fileInfoLen uint64 = 0;
var curFileIndex uint64 = 0;

var DEBUG_LEVEL uint64 = 0;

func main() {
    var errno uint64;
    var i uint64;

    if len(os.Args) == 1 {
        libgogo.ExitError("Usage: gogo file1.go [file2.go ...]",1);
    }

    for i=1; i < uint64(len(os.Args)) ; i= i+1 { // (SC) TODO: remove cast and os.Args
        curFileIndex = i-1;
        fileInfo[curFileIndex].filename = os.Args[i];
        fileInfo[curFileIndex].lineCounter = 1;
        fileInfo[curFileIndex].charCounter = 1;

        fileInfo[curFileIndex].fd = libgogo.FileOpen(fileInfo[curFileIndex].filename, 0);
        if (fileInfo[curFileIndex].fd == 0) {
            GlobalError("Cannot open file.");
        }
    }
    fileInfoLen = i-1;

    for curFileIndex=0;curFileIndex<fileInfoLen;curFileIndex=curFileIndex+1 {
        Parse();
    }

    for curFileIndex=0;curFileIndex<fileInfoLen;curFileIndex=curFileIndex+1 {
        errno = libgogo.FileClose(fileInfo[curFileIndex].fd);
        if errno != 0 {
            GlobalError("Cannot close file.");
        }
    }
}
