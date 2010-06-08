// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

type LineDesc struct {
    Line string;
    Offset uint64;
};

//
// Function fetches the next line from the assembly file.
//
func GetLine(ld *LineDesc) {
    var line string;
    var singleChar byte;
    for singleChar = GetCharWrapped(); (singleChar != 0) && (singleChar != 10); singleChar = GetCharWrapped() {
        libgogo.CharAppend(&line, singleChar);
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
    return 16;
}

//
// Function processing a line and fixing offsets if necessary
//
func FixOffsetIfNecessary(ld *LineDesc) string {
    var packageName string;
    var functionName string;
    var i uint64;
    var j uint64;
    var strLen uint64;
    var position uint64 = 0;
    var size uint64;
    var oldsize uint64;
    var numstr string;
    var newLine string;

    strLen = libgogo.StringLength(ld.Line);
    for i=0 ; i < strLen ; i=i+1 {
        if position == 3 { // Create new offset using the parametersize of packageName·functionName
            size = GetParameterSize(packageName, functionName);
            // Dismiss '-', which is mandatory. I.e.: -8(SP)
            libgogo.CharAppend(&newLine, ld.Line[i]); 
            i = i+1; 
            // Read the number in -<number>(SP)
            for ;ld.Line[i]!='(';i=i+1 {
                libgogo.CharAppend(&numstr, ld.Line[i]);
            }
            oldsize = libgogo.StringToInt(numstr); // Convert to uint64
            size = size + oldsize; // Caluclate new size
            numstr = libgogo.IntToString(size); // Convert back
            libgogo.StringAppend(&newLine, numstr); // Append
            libgogo.CharAppend(&newLine,'('); // Finally append '(' again
            position = 0; // Reset flag to starting indication
            continue;
        }
        if position == 2 { // Fetch functionName until "##"
            if ld.Line[i] == '#' {
                j = i-1;
                if ld.Line[j] == '#' {
                    position = position +1;
                }
            } else {
                libgogo.CharAppend(&functionName, ld.Line[i]);
            }
        }
        if position == 1 { // Fetch packageName until '·'
            if ld.Line[i] != 194 {
                if ld.Line[i] == 183 {
                    position = position + 1;
                } else {
                    libgogo.CharAppend(&packageName, ld.Line[i]);
                }
            }
        }
        if position == 0 { // Append characters until separator "##" is reached
            if ld.Line[i] == '#' {
                if i==0 {
                    LinkError("Found '#' at linestart. Not a valid object file.");
                } else {
                    j = i-1;
                    if ld.Line[j] == '#' {
                        position = position +1;
                    }
                }                 
            }  else {
                if (i>0) {
                    j = i-1;
                    if ld.Line[j] == '#' {                    
                        LinkError("Single '#' not allowed. Not a valid object file.");
                    }
                } 
                libgogo.CharAppend(&newLine, ld.Line[i]);
            }
        }
    }
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
            newLine = FixOffsetIfNecessary(&ld);
            libgogo.PrintString(newLine);
            libgogo.PrintString("\n");
        }
        GetLine(&ld);
    }
}
