// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// Code that is used to generate the required conditionals (labels, jumps) that
// are used by expressions.
//

package main

import "./libgogo/_obj/libgogo"

//
// Function used to print labels.
// Internally using the asm version.
//
func PrintLabelWrapped(ed *ExpressionDescriptor, global uint64, localBranch uint64, suffix string) {
    var labelString string;
    if global == 0 {
        labelString = GetGlobLabel(ed, suffix);
    } else {
        labelString = GetSubLabel(ed, localBranch, suffix);
    }
    PrintLabel(labelString);
}

//
// Generates a new Sublabel that can be used in a jump and later on be fetched
// by GetSubLabel().
//
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

//
// Returns a Sublabel depending on an ExpressionDescriptors internal state.
// (=Global Label + True|False Branch)
//
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

//
// Returns the global label that is represented by an ExpressionDescriptor.
//
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
