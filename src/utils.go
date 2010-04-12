// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

func PrintErrorHead() {
    libgogo.PrintString(fileInfo[curFileIndex].filename);
    libgogo.PrintString(":");
    libgogo.PrintNumber(fileInfo[curFileIndex].lineCounter);
    libgogo.PrintString(":");
    libgogo.PrintNumber(fileInfo[curFileIndex].charCounter);   
}

func GlobalError(msg string) {
    PrintErrorHead();
    libgogo.PrintString(": error: ");
    libgogo.ExitError(msg,1);
}

func ScanErrorString(msg string) {
    PrintErrorHead();
    libgogo.PrintString(": syntax error: ");
    libgogo.ExitError(msg,2);
}

func ScanErrorChar(char byte) {
    PrintErrorHead();
    libgogo.PrintString(": syntax error: Unknown char '");
    libgogo.PrintChar(char);
    libgogo.PrintString("'.\n");
    libgogo.Exit(2);
}

//
// Function printing a parse error using only libgogo.
// ue ... unexpected token
// e .... array of expected tokens
// eLen . actual length (items) of array
//
func ParseError(ue uint64, e [255]uint64, eLen uint64) {
    var i uint64;
    var str string;

    PrintErrorHead();
    libgogo.PrintString(": syntax error: unexpected token '");
    str = TokenToString(ue);
    libgogo.PrintString(str);
    libgogo.PrintString("'");

    if eLen > 0 {
        libgogo.PrintString(", expecting one of: ");
        str = TokenToString(e[i]);
        libgogo.PrintString(str);
        for i = 1; i < eLen; i = i+1 {
            str = TokenToString(e[i]);
            libgogo.PrintString(str);
            libgogo.PrintString(", ");        
        }
    } 
    libgogo.PrintString("\n");
    libgogo.Exit(3); 
}

