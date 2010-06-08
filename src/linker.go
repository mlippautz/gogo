// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

type LineDesc struct {
    Line string;
    Offset uint64;

    NeedsFix uint64;
    PackageName string;
    FunctionName string;
};

//
// Function fetches the next line from the assembly file.
//
func GetLine(ld *LineDesc) {
    var line string;
    var singleChar byte;

    singleChar = GetCharWrapped();
    if singleChar == '#' {
        GetCharWrapped(); // Abolish '#'
        singleChar = GetCharWrapped(); // should be number
        ld.NeedsFix = libgogo.ToIntFromByte(singleChar);
        ld.NeedsFix = ld.NeedsFix - 48;
        GetCharWrapped(); // Abolish '#'
        GetCharWrapped(); // Abolish '#'
        for singleChar = GetCharWrapped(); singleChar != 194 ; singleChar = GetCharWrapped() {
            libgogo.CharAppend(&ld.PackageName, singleChar);
        }
        GetCharWrapped(); // Abolish second part of dot (183)
        for singleChar = GetCharWrapped(); (singleChar != '#'); singleChar = GetCharWrapped() {
            libgogo.CharAppend(&ld.FunctionName, singleChar);
        }
        for singleChar = GetCharWrapped(); (singleChar != 0) && (singleChar != 10); singleChar = GetCharWrapped() { // Dismiss rest of line
        }
    } else {
        libgogo.CharAppend(&line, singleChar);
        for singleChar = GetCharWrapped(); (singleChar != 0) && (singleChar != 10); singleChar = GetCharWrapped() {
            libgogo.CharAppend(&line, singleChar);
        }
    }

    if singleChar == 0 {
        tok.id = TOKEN_EOS;
    }

    ld.Line = line;
    ld.Offset = 0;
}

func GetNextSymToken(ld *LineDesc) string {
    var symtoken string;
    var i uint64;
    var len uint64;
    len = libgogo.StringLength(ld.Line);
    for i = ld.Offset; (i < len) && (ld.Line[i] != ','); i=i+1 {
        if (ld.Line[i] != '/') && (ld.Line[i] != ' ') {
            libgogo.CharAppend(&symtoken, ld.Line[i]);
        }
    }
    ld.Offset = i+1;
    return symtoken;
}


func ParseLine(ld *LineDesc) {
    // Something like 
    // Type, Ndx, Name, Ret, Params [,...]
    // FUNC ,UND    ,test·test      ,           ,uint64

    var symtoken string;
    symtoken = GetNextSymToken(ld);
    symtoken = GetNextSymToken(ld);
    symtoken = symtoken;
}

func GetParameterSize(packageName string, functionName string) uint64 {
    return 48;
}

//
// Function processing a line and fixing offsets if necessary
//
func FixOffset(ld *LineDesc) string {
    var i uint64;
    var strLen uint64;
    var size uint64;
    var oldsize uint64;
    var numstr string;
    var newLine string;

    if ld.NeedsFix == 1 { // Type 1 fix of offsets
        strLen = libgogo.StringLength(ld.Line);
        size = GetParameterSize(ld.PackageName, ld.FunctionName);
        for i = 0; ld.Line[i] != '-' ; i = i +1 {
            libgogo.CharAppend(&newLine, ld.Line[i]);
        }
        libgogo.CharAppend(&newLine, ld.Line[i]);
        for i = i+1 /*Dismiss '-'*/ ;ld.Line[i]!='(';i=i+1 {
            libgogo.CharAppend(&numstr, ld.Line[i]);
        }
        oldsize = libgogo.StringToInt(numstr);
        size = size + oldsize;
        size = size - 100000;
        numstr = libgogo.IntToString(size);
        libgogo.StringAppend(&newLine, numstr);
        for ; i < strLen; i = i +1 {
            libgogo.CharAppend(&newLine, ld.Line[i]);
        }
    }
    if ld.NeedsFix == 2 { // Type 2 fix of offsets

    }

    ld.NeedsFix = 0;
    ld.Line = newLine;
    return newLine;
}

//
// Main linking method
//
func Link() {
    var newLine string;
    var strCmp uint64;
    var symtable uint64 = 0;
    var ld LineDesc;

    ld.NeedsFix = 0;
    ResetToken();
    GetLine(&ld);
    for ;tok.id != TOKEN_EOS;{
        strCmp = libgogo.StringCompare("// ##START_SYM_TABLE", ld.Line);
        if (strCmp == 0) {
            symtable = 1;
            GetLine(&ld); // Proceed to next line
        } 
        strCmp = libgogo.StringCompare("// ##END_SYM_TABLE", ld.Line);
        if (strCmp == 0) {
            symtable = 0;
            GetLine(&ld); // Proceed to next line
        }
        if symtable != 0 { // Parse symtable entries
            ParseLine(&ld);
        } else { // Parse normal lines and fix everything
            if ld.NeedsFix != 0 {
                GetLine(&ld);
                newLine = FixOffset(&ld);
                libgogo.PrintString(newLine);
                libgogo.PrintString("\n");
            } else {
                //libgogo.PrintString(ld.Line);
                //libgogo.PrintString("\n");
            }

        }
        GetLine(&ld);
    }
}
