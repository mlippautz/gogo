// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

var InSymTable uint64 = 0;

func GetLine() string {
    var line string;
    var singleChar byte;
    for singleChar = GetCharWrapped(); (singleChar != 0) && (singleChar != 10); singleChar = GetCharWrapped() {
        libgogo.CharAppend(&line, singleChar);
    }
    if singleChar == 0 {
        tok.id = TOKEN_EOS;
    }
    return line;
}

func GetNextSymToken(lineRest string, offset *uint64) string {
    var symtoken string;
    var i uint64;
    var len uint64;
    len = libgogo.StringLength(lineRest);
    for i = *offset; (i < len) && (lineRest[i] != ','); i=i+1 {
        if (lineRest[i] != '/') && (lineRest[i] != ' ') {
            libgogo.CharAppend(&symtoken, lineRest[i]);
        }
    }
    *offset = i+1;
    return symtoken;
}

func ParseLine(line string) {
    // Something like 
    // Type, Ndx, Name, Ret, Params [,...]
    // FUNC ,UND    ,test·test      ,           ,uint64
    var offset uint64 = 0;
    var symtoken string;
    symtoken = GetNextSymToken(line,&offset);
    libgogo.PrintString("Type: ");
    libgogo.PrintString(symtoken);
    libgogo.PrintString("\n");

    symtoken = GetNextSymToken(line,&offset);
    libgogo.PrintString("Defined: ");
    libgogo.PrintString(symtoken);
    libgogo.PrintString("\n");

    symtoken = GetNextSymToken(line,&offset);
    libgogo.PrintString("Packagename+Var-/Functionname: ");
    libgogo.PrintString(symtoken);
    libgogo.PrintString("\n");

    symtoken = GetNextSymToken(line,&offset);
    libgogo.PrintString("Return, vartype: ");
    libgogo.PrintString(symtoken);
    libgogo.PrintString("\n");

}

func Link() {
    var line string;
    var strCmp uint64;

    tok.id = 0;
    tok.nextChar = 0;
    tok.nextToken = 0;   
    tok.llCnt = 0; 
    
    for line = GetLine(); tok.id != TOKEN_EOS ;line = GetLine() {
        strCmp = libgogo.StringCompare("// ##START_SYM_TABLE", line)
        if (strCmp == 0) {
            libgogo.PrintString("Symboltable start\n");
            InSymTable = 1;
            line = GetLine();
        } 
        strCmp = libgogo.StringCompare("// ##END_SYM_TABLE", line)
        if (strCmp == 0) {
            libgogo.PrintString("Symboltable end\n");
            InSymTable = 0;
        }
        if InSymTable != 0 {
            ParseLine(line);
            libgogo.PrintString("\n");
        }        
    }
}
