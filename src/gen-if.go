// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// Code that is used to generate the 'if' and 'else' constructs.
//

package main

import "./libgogo/_obj/libgogo"

func GenerateIfStart(item *libgogo.Item, ed *ExpressionDescriptor) {
    var labelString string;
    var jmp string;

    labelString = GenerateSubLabel(ed,1,"END");
    if ed.Not == 0 {
        jmp = GetJump(item.C, 0);
    } else {
        labelString = GenerateSubLabel(ed,0,"END");
        jmp = GetJump(item.C, 1);
        SwapExpressionBranches(ed);
    }
    PrintJump(jmp, labelString);

    // Important: Since last jump is a positive one, we have to start with the
    // negative path
    if ed.F != 0 {
        PrintLabelWrapped(ed, 1 /*local*/, 0 /*positive*/, "END");
        //labelString = GetSubLabel(ed,0,"END");
        //PrintLabel(labelString);
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

