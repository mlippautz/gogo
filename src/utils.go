// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

var maxErrors uint64 = 10;
var errors uint64 = 0;

func CheckDebugLevel(debugLevel uint64) uint64 {
    var retVal uint64 = 0;
    if debugLevel <= DEBUG_LEVEL {
        retVal = 1;
    }
    return retVal;
}

func PrintDebugString(msg string, debugLevel uint64) {
    var flag uint64;
    flag = CheckDebugLevel(debugLevel);
    if flag  == 1 {
        PrintHead();   
        libgogo.PrintString(": DEBUG: ");
        libgogo.PrintString(msg);
        libgogo.PrintString("\n");
    }
}

func PrintDebugChar(char byte, debugLevel uint64) {
    var flag uint64;
    flag = CheckDebugLevel(debugLevel);
    if flag  == 1 {
        PrintHead();   
        libgogo.PrintString(": DEBUG: ");
        libgogo.PrintChar(char);
        libgogo.PrintString("\n");
    }
}

func BuildHead() string {
    var text string = "";
    var temp string;
    libgogo.StringAppend(&text, fileInfo[curFileIndex].filename);
    libgogo.StringAppend(&text, ":");
    temp = libgogo.IntToString(fileInfo[curFileIndex].lineCounter);
    libgogo.StringAppend(&text, temp);
    libgogo.StringAppend(&text, ":");
    temp = libgogo.IntToString(fileInfo[curFileIndex].charCounter);
    libgogo.StringAppend(&text, temp);
    return text;
}

func PrintHead() {
    var head string;
    head = BuildHead();
    libgogo.PrintString(head);
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

func ParseErrorWeak(ue uint64, expectedToken1 uint64, expectedToken2 uint64, expectedLen uint64) {
    var str string;
    
    errors = errors+1;
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

    if expectedLen > 0 {
        libgogo.PrintString(", expecting one of: ");
        str = TokenToString(expectedToken1);
        libgogo.PrintString(str);
        if expectedLen >= 2 {
            libgogo.PrintString(", ");
            str = TokenToString(expectedToken2);
            libgogo.PrintString(str);
        }
    } 
    libgogo.PrintString("\n");
    if errors == maxErrors {
        libgogo.PrintString("Maximum number of errors (");
        libgogo.PrintNumber(maxErrors);
        libgogo.PrintString(") reached. Exiting.\n");
        libgogo.Exit(3);
    }
}

//
// Function printing a parse error using only libgogo.
// ue ... unexpected token
// expectedToken1,2 .... expected tokens
// expectedLen . actual number of expected tokens (others will be ignored)
//
func ParseErrorFatal(ue uint64, expectedToken1 uint64, expectedToken2 uint64, expectedLen uint64) {
    ParseErrorWeak(ue, expectedToken1, expectedToken2, expectedLen);
    libgogo.Exit(3); 
}

//
// Print an error, stop compilation and parse rest.
//
func GenErrorWeak(msg string) {
    PrintHead();
    libgogo.PrintString(": generation error: ");
    libgogo.PrintString(msg);
    libgogo.PrintString("\n");
    Compile = 0;
    ParserSync();
}

//
// Linker error. Always Fatal.
//
func LinkError(msg1 string, msg2 string, msg3 string, msg4 string, msg5 string) {
    PrintHead();
    libgogo.PrintString(": linker error: ");
    libgogo.PrintString(msg1);
    libgogo.PrintString(msg2);
    libgogo.PrintString(msg3);
    libgogo.PrintString(msg4);
    libgogo.PrintString(msg5);
    libgogo.PrintString("\n");
    libgogo.Exit(5);
}

//
// Linker warning
//
func LinkWarn(msg1 string, msg2 string, msg3 string, msg4 string, msg5 string) {
    PrintHead();
    libgogo.PrintString(": linker warning: ");
    libgogo.PrintString(msg1);
    libgogo.PrintString(msg2);
    libgogo.PrintString(msg3);
    libgogo.PrintString(msg4);
    libgogo.PrintString(msg5);
    libgogo.PrintString("\n");
}

func SymbolTableError(msg string, position string, msg2 string, identifier string) {
    var strLen uint64;
    PrintHead();
    libgogo.PrintString(": symbol table error: ");
    libgogo.PrintString(msg);
    libgogo.PrintString(" ");
    libgogo.PrintString(position);
    strLen = libgogo.StringLength(position);
    if strLen != 0 {
        libgogo.PrintString(" ");
    }
    libgogo.PrintString(msg2);
    libgogo.PrintString(" '");
    libgogo.PrintString(identifier);
    libgogo.PrintString("'\n");
    PrintGlobalSymbolTable(); //May be useful for debugging purposes
    PrintLocalSymbolTable();    
    libgogo.Exit(4);
}

