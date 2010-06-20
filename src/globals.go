// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// File holding global variables that are needed by two or more modules
//

package main

//
// Struct holding the information about a file that is compiled
//
type FileInfo struct {
    filename string;
    lineCounter uint64;
    charCounter uint64;
    fd uint64;
};

type ExpressionDescriptor struct {
    //
    // Labeling information
    //
    ExpressionDepth uint64; // The current expression depth.
    IncCnt uint64; // Some incremental counter to guarantee uniqueness
    CurFile string; // Current file begining with a specified prefix. 
    CurLine uint64; // Current line in parser. Used for label generation.

    //
    // True/False branches (merge) information
    //
    T uint64; // True branch
    F uint64; // False branch
    TDepth uint64; /* Depth when true branch has been started. Used for merge 
      and printing. */
    FDepth uint64; // Same as true depth.
    Not uint64; // Flag indicating not branch

    //
    // Break continue information
    //
    ForEd *ExpressionDescriptor;
    ForPost uint64;
};

//
// Fileinformation for all files that are compiled in this run
// Is limited by 40 to reduce memory consumption, but to allow
// self compilation via "./gogo libgogo/*.go *.go"
//
var fileInfo [40]FileInfo;
var fileInfoLen uint64 = 0;
var curFileIndex uint64 = 0;

//
// Compiler flag indicating in which mode the compiler is
// 0 ... parsing only
// 1 ... compile (code generation)
// 2 ... link
//
var Compile uint64 = 0;

//
// A very basic debug flag
// Set to 1000 to enable all parsing strings
// Set to 100 to enable all symbol tables
// Set to 10 to enable asm debugging
//
var DEBUG_LEVEL uint64 = 10;

//
// Package name of currently processed file
//
var CurrentPackage string = "<no package>";

var InsideFunction uint64 = 0;
var InsideStructDecl uint64 = 0;
var InsideFunctionVarDecl uint64 = 0;


