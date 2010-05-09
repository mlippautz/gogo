// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

func CheckDebugLevel(debugLevel uint64) uint64 {
    var retVal uint64 = 0;
    if debugLevel <= DEBUG_LEVEL {
        retVal = 1;
    }
    return retVal;
}

func PrintDebugString(msg string, debugLevel uint64) {
    if CheckDebugLevel(debugLevel) == 1 {
        PrintHead();   
        libgogo.PrintString(": DEBUG: ");
        libgogo.PrintString(msg);
        libgogo.PrintString("\n");
    }
}

func PrintDebugChar(char byte, debugLevel uint64) {
    if CheckDebugLevel(debugLevel) == 1 {
        PrintHead();   
        libgogo.PrintString(": DEBUG: ");
        libgogo.PrintChar(char);
        libgogo.PrintString("\n");
    }
}

func PrintHead() {
    libgogo.PrintString(fileInfo[curFileIndex].filename);
    libgogo.PrintString(":");
    libgogo.PrintNumber(fileInfo[curFileIndex].lineCounter);
    libgogo.PrintString(":");
    libgogo.PrintNumber(fileInfo[curFileIndex].charCounter);   
}

func GlobalError(msg string) {
    PrintHead();
    libgogo.PrintString(": error: ");
    libgogo.ExitError(msg,1);
}

func ScanErrorString(msg string) {
    PrintHead();
    libgogo.PrintString(": syntax error: ");
    libgogo.ExitError(msg,2);
}

func ScanErrorChar(char byte) {
    PrintHead();
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
func ParseError(ue uint64, e [2]uint64, eLen uint64) {
    var i uint64;
    var str string;

    PrintHead();
    libgogo.PrintString(": syntax error: unexpected token '");
    str = TokenToString(ue);
    libgogo.PrintString(str);
    libgogo.PrintString("'");

    if tok.id == TOKEN_INTEGER {
        libgogo.PrintString(" (value: ");
        libgogo.PrintNumber(tok.intValue);
        libgogo.PrintString(")");
    }

    if tok.id == TOKEN_STRING {
        libgogo.PrintString(" (value: ");
        libgogo.PrintString(tok.strValue);
        libgogo.PrintString(")");
    }

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

func SymbolTableError(msg string, position string, msg2 string, identifier string) {
    PrintHead();
    libgogo.PrintString(": symbol table error: ");
    libgogo.PrintString(msg);
    libgogo.PrintString(" ");
    libgogo.PrintString(position);
    if libgogo.StringLength(position) != 0 {
        libgogo.PrintString(" ");
    }
    libgogo.PrintString(msg2);
    libgogo.PrintString(" '");
    libgogo.PrintString(identifier);
    libgogo.PrintString("'\n");
    libgogo.Exit(4);
}

