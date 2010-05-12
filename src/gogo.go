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
// Set to 100 to enable all symbol tables
//
var DEBUG_LEVEL uint64 = 0;

//
// Entry point of the compiler
//
func main() {
    var errno uint64;
    var i uint64;
    var temptype *libgogo.TypeDesc;
    var tempobject *libgogo.ObjectDesc;

    libgogo.GetArgv();

    if libgogo.Argc <= 1 {
        libgogo.ExitError("Usage: gogo file1.go [file2.go ...]",1);
    }

    if libgogo.Argc > 11 {
        libgogo.ExitError("Cannot compile more than 10 files at once",1);
    }

    //Default data types
    temptype = libgogo.NewType("uint64", "", 8, nil);
    GlobalTypes = libgogo.AppendType(temptype, GlobalTypes);
    temptype = libgogo.NewType("byte", "", /*1*/ 8, nil); //Use size 8 as everything has to be 64-bit aligned anyways
    GlobalTypes = libgogo.AppendType(temptype, GlobalTypes);
    temptype = libgogo.NewType("string", "", 16, nil);
    GlobalTypes = libgogo.AppendType(temptype, GlobalTypes);

    //Default objects
    tempobject = libgogo.NewObject("nil", libgogo.CLASS_VAR);
    libgogo.SetObjType(tempobject, nil);
    libgogo.FlagObjectTypeAsPointer(tempobject); //nil is a pointer to no specified type (universal)
    GlobalObjects = libgogo.AppendObject(tempobject, GlobalObjects);

    InitFreeRegisters(); //Init registers for code generation

    for i=1; i < libgogo.Argc ; i= i+1 {
        curFileIndex = i-1;
        fileInfo[curFileIndex].filename = libgogo.Argv[i];
        fileInfo[curFileIndex].lineCounter = 1;
        fileInfo[curFileIndex].charCounter = 1;
        
        fileInfo[curFileIndex].fd = libgogo.FileOpen(libgogo.Argv[i], 0);
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

    if CheckDebugLevel(100) == 1 { //Global symbol table
        libgogo.PrintString("\nGlobal symbol table:\n");
        libgogo.PrintString("--------------------\n");
        libgogo.PrintTypes(GlobalTypes);
        libgogo.PrintObjects(GlobalObjects);
    }
}
