// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.


package main

import "./libgogo/_obj/libgogo"

func GenerateIfStart(item *libgogo.Item, ed *ExpressionDescriptor) {
    var labelString string;
    var jmp string;
    var tmp uint64;

    labelString = GenerateSubLabel(ed,1,"END");
    if ed.Not == 0 {
        jmp = GetJump(item.C, 0);
    } else {
        labelString = GenerateSubLabel(ed,0,"END");
        jmp = GetJump(item.C, 1);
        tmp = ed.T;
        ed.T = ed.F;
        ed.F = tmp;
    }
    PrintJump(jmp, labelString);

    // Important: Since last jump is a positive one, we have to start with the
    // negative path
    if ed.F != 0 {
        labelString = GetSubLabel(ed,0,"END");
        PrintLabel(labelString);
    }
    labelString = GetGlobLabel(ed,"END");
    PrintJump("JMP", labelString);

    // Positive branch starts after this label, thus insert last remaining 
    // positive label (if available) here
    if ed.T != 0 {
        labelString = GetSubLabel(ed,1,"END");
        PrintLabel(labelString);
    }
}

func GenerateIfEnd(ed *ExpressionDescriptor) {
    var labelString string;
    labelString = GetGlobLabel(ed, "END");
    PrintLabel(labelString);
}

func GenerateElseStart(ed *ExpressionDescriptor) {
    var labelString string;
    labelString = GetGlobLabel(ed, "ELSE_END");
    PrintJump("JMP", labelString);
    labelString = GetGlobLabel(ed, "END");
    PrintLabel(labelString);
}

func GenerateElseEnd(ed *ExpressionDescriptor) {
    var labelString string;
    labelString = GetGlobLabel(ed, "ELSE_END");
    PrintLabel(labelString);
}

func GenerateSubLabel(ed *ExpressionDescriptor, i uint64, label string) string {
    var str string;
    var tmpStr string;
    libgogo.StringAppend(&str, ed.CurFile);
    libgogo.StringAppend(&str, "_");
    tmpStr = libgogo.IntToString(ed.CurLine);
    libgogo.StringAppend(&str, tmpStr);
    libgogo.StringAppend(&str, "_");

    if i == 0 {
        if ed.F == 0 {
            tmpStr = libgogo.IntToString(ed.IncCnt);
            ed.F = ed.IncCnt;
            ed.FDepth = ed.ExpressionDepth+1;
            ed.IncCnt = ed.IncCnt + 1;
        } else {
            tmpStr = libgogo.IntToString(ed.F);
        }
    } else {
        if ed.T == 0 {
            tmpStr = libgogo.IntToString(ed.IncCnt);
            ed.T = ed.IncCnt;
            ed.TDepth = ed.ExpressionDepth+1;
            ed.IncCnt = ed.IncCnt + 1;
        } else {
            tmpStr = libgogo.IntToString(ed.T);
        }
    }
    libgogo.StringAppend(&str, tmpStr);

    libgogo.StringAppend(&str, "_");
    libgogo.StringAppend(&str, label);
    return str;
}

func GetSubLabel(ed *ExpressionDescriptor, i uint64, label string) string {
    var str string;
    var tmpStr string;
    libgogo.StringAppend(&str, ed.CurFile);
    libgogo.StringAppend(&str, "_");
    tmpStr = libgogo.IntToString(ed.CurLine);
    libgogo.StringAppend(&str, tmpStr);
    libgogo.StringAppend(&str, "_");

    if i == 0 {
        tmpStr = libgogo.IntToString(ed.F);
    } else {
        tmpStr = libgogo.IntToString(ed.T);
    }
    libgogo.StringAppend(&str, tmpStr);

    libgogo.StringAppend(&str, "_");
    libgogo.StringAppend(&str, label);
    return str;
}

func GetGlobLabel(ed *ExpressionDescriptor, label string) string {
    var str string;
    var tmpStr string;
    libgogo.StringAppend(&str, ed.CurFile);
    libgogo.StringAppend(&str, "_");
    tmpStr = libgogo.IntToString(ed.CurLine);
    libgogo.StringAppend(&str, tmpStr);
    libgogo.StringAppend(&str, "_");
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
