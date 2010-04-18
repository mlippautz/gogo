// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

//
// Struct holding the information about a file that is compiled
//
type FileInfo struct {
    filename string;
    lineCounter uint64;
    charCounter uint64;
    fd uint64;
}; 

//
// Fileinformation for all files that are compiled in this run
// Is limited by 10 to reduce memory etc.
//
var fileInfo [10]FileInfo;
var fileInfoLen uint64 = 0;
var curFileIndex uint64 = 0;

//
// A very basic debug flag
// Set to 1000 to enable all parsing strings
//
var DEBUG_LEVEL uint64 = 0;

//
// Entry point of the compiler
//
func main() {
    var errno uint64;
    var i uint64;

    libgogo.GetArgv()

    if libgogo.Argc <= 1 {
        libgogo.ExitError("Usage: gogo file1.go [file2.go ...]",1);
    }

    if libgogo.Argc > 11 {
        libgogo.ExitError("Cannot compile more than 10 files at once",1);
    }

    for i=1; i < libgogo.Argc ; i= i+1 {
        curFileIndex = i-1;
        fileInfo[curFileIndex].filename = libgogo.Argv[i];
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
