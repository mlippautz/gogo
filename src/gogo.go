// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

//
// Entry point of the compiler
//
func main() {
    var errno uint64;
    var i uint64;
    var j uint64;
    var k uint64;
    var singleChar byte;

    libgogo.GetArgv();

    ParseOption();

    if libgogo.Argc > 40 {
        libgogo.ExitError("Cannot compile more than 38 files at once",1);
    }

    InitSymbolTable(); //Initialize symbol table
    InitFreeRegisters(); //Init registers for code generation

    ResetCode();

    for i=2; i < libgogo.Argc ; i= i+1 {
        curFileIndex = i-2;
        fileInfo[curFileIndex].filename = libgogo.Argv[i];
        fileInfo[curFileIndex].lineCounter = 1;
        fileInfo[curFileIndex].charCounter = 1;
        
        fileInfo[curFileIndex].fd = libgogo.FileOpen(libgogo.Argv[i], 0);
        if (fileInfo[curFileIndex].fd == 0) {
            GlobalError("Cannot open file.");
        }
    }
    fileInfoLen = i-2;

    if (Compile == 0) || (Compile == 1) { // Compile == 0|1 => Parse/Compile
        for curFileIndex=0;curFileIndex<fileInfoLen;curFileIndex=curFileIndex+1 {
            i = libgogo.StringLength(fileInfo[curFileIndex].filename);
            if i > 2 { //Check for assembly files
                j = i - 2;
                k = i - 1;
                if (Compile == 1) && (fileInfo[curFileIndex].filename[j] == '.') && (fileInfo[curFileIndex].filename[k] == 's') { 
                    // Assembly file: compile if in compiler mode
                    if curFileIndex == 0 {
                        GlobalError("The first file in the list cannot be an assembly file");
                    }
                    for singleChar = libgogo.GetChar(fileInfo[curFileIndex].fd); singleChar != 0; singleChar = libgogo.GetChar(fileInfo[curFileIndex].fd) { //Copy file to output character by character
                        /* middot '·' is U+00C7 => UTF-8: C2B7*/
                        if singleChar != 194 /*C2*/ {
                            if singleChar == 183 /*B7*/ { //Prepend package name
                                PrintCodeOutput(CurrentPackage); 
                                PrintCodeOutputChar(194);
                                PrintCodeOutputChar(183);
                            } else {
                                PrintCodeOutputChar(singleChar);
                            }
                        }
                    }
                } else {
                    // Go file: parse and/or compile
                    Parse();
                }
            } else { //Go file: parse and/or compile
                Parse();
            }
        }
    } else { // Compile==2 => Link
        Link();
    }

    for curFileIndex=0;curFileIndex<fileInfoLen;curFileIndex=curFileIndex+1 {
        errno = libgogo.FileClose(fileInfo[curFileIndex].fd);
        if errno != 0 {
            GlobalError("Cannot close file.");
        }
    }

    PrintGlobalSymbolTable();
    
    if Compile == 1 {
        i = libgogo.GetAlignedObjectListSize(GlobalObjects); //Get required data segment size
        SetDataSegmentSize(i); //Set data segment size
        PrintFile(GlobalFunctions, GlobalTypes); //Print compiled output to file
    }
}

func ParseOption() {
    var strIndicator uint64;
    var done uint64 = 0;

    // handle -h and --help    
    strIndicator = libgogo.StringCompare("--help", libgogo.Argv[1]);
    if strIndicator != 0 {
        strIndicator = libgogo.StringCompare("-h", libgogo.Argv[1]);
    }

    if strIndicator == 0 {
        libgogo.PrintString("Usage: gogo option file1.go [file2.go ...]\n\n");
        libgogo.PrintString("GoGo - A go compiler\n\n");
        libgogo.PrintString("Options:\n");
        libgogo.PrintString("-h, --help     show this help message and exit\n");
        libgogo.PrintString("-p,            parser mode\n");
        libgogo.PrintString("-c             compiler mode\n");
        libgogo.PrintString("-l             linker mode\n");
        libgogo.Exit(1);
    }

    strIndicator = libgogo.StringCompare("-c", libgogo.Argv[1]);
    if (done == 0) && (strIndicator == 0) {
        Compile = 1;
        done = 1;
    }

    strIndicator = libgogo.StringCompare("-p", libgogo.Argv[1]);
    if (done == 0) && (strIndicator == 0) {
        Compile = 0;
        done = 1;
    }

    strIndicator = libgogo.StringCompare("-l", libgogo.Argv[1]);
    if (done == 0) && (strIndicator == 0) {
        Compile = 2;
        done = 1;
    }
    
    if done == 0 {
        libgogo.ExitError("Usage: gogo option file1.go [file2.go ...]",1);
    }

    if libgogo.Argc <= 2 {
        libgogo.ExitError("Usage: gogo option file1.go [file2.go ...]",1);
    }
}
