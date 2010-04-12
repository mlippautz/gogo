// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

func ScanErrorString(msg string) {
    libgogo.PrintString(filename);
    libgogo.PrintString(":");
    libgogo.PrintNumber(lineCounter);
    libgogo.PrintString(": syntax error: ");
    libgogo.ExitError(msg,1);
}

func ScanErrorChar(char byte) {
    libgogo.PrintString(filename);
    libgogo.PrintString(":");
    libgogo.PrintNumber(lineCounter);
    libgogo.PrintString(": syntax error: Unknown char '");
    libgogo.PrintChar(char);
    libgogo.PrintString("'.\n");
    libgogo.Exit(1);
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

    libgogo.PrintString(filename);
    libgogo.PrintString(":");
    libgogo.PrintNumber(lineCounter);
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
    libgogo.Exit(2); 
}

