// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.


package main

import "./libgogo/_obj/libgogo"

func GenerateIfStart(item *libgogo.Item, ed ExpressionDescriptor) {
    var labelString string;
    var jmp string;
    labelString = GenerateIfLabel(ed.CurFile,ed.CurLine,0,"END");
    jmp = GetJump(item.C, 1);
    PrintJump(jmp, labelString);
    labelString = GenerateIfLabel(ed.CurFile, ed.CurLine, ed.ExpressionDepth-1, "END");
    PrintLabel(labelString);
    labelString = GenerateIfLabel(ed.CurFile, ed.CurLine, 0, "END");
    PrintJump("JMP",labelString);
    labelString = GenerateIfLabel(ed.CurFile, ed.CurLine, 0, "OK");
    PrintLabel(labelString);
}

func GenerateIfEnd(item *libgogo.Item, ed ExpressionDescriptor) {
    var labelString string;
    labelString = GenerateIfLabel(ed.CurFile, ed.CurLine, 0, "END");
    PrintLabel(labelString);
}

func GenerateIfLabel(filename string, line uint64, local uint64, label string) string {
    var str string;
    var tmpStr string;
    libgogo.StringAppend(&str, filename);
    libgogo.StringAppend(&str, "_");
    tmpStr = libgogo.IntToString(line);
    libgogo.StringAppend(&str, tmpStr);
    libgogo.StringAppend(&str, "_");
    if local != 0 {
        tmpStr = libgogo.IntToString(local);
        libgogo.StringAppend(&str, tmpStr);
        libgogo.StringAppend(&str, "_");
    }
    libgogo.StringAppend(&str, label);
    return str;
}

//
// Returns the jump expression that represents a given conditional
//
func GetJump(condition uint64, invert uint64) string {
    var jmp string;

    if invert != 0 { // invert logic if necessary (T|F jump)
        if (condition == 0) || (condition == 2) || (condition == 4) {
            condition = condition + 1;
        } else {
            condition = condition - 1;
        }
    }

    // Information taken from go/src/cmd/6a/<file>
    if condition == libgogo.REL_EQ {
        jmp = "JE"; // lex.c:466
    }
    if condition == libgogo.REL_NEQ {
        jmp = "JNE"; // lex.c:468
    }
    if condition == libgogo.REL_GT {
        jmp = "JG"; // lex.c:494
    }
    if condition == libgogo.REL_GTE {
        jmp = "JGE"; // lex.c:489
    }
    if condition == libgogo.REL_LT {
        jmp = "JL"; // lex.c:487
    }
    if condition == libgogo.REL_LTE {
        jmp = "JLE"; // lex.c:491
    }

    return jmp;
}
