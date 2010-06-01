// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// Code that is used to generate the 'for' constructs.
//

package main

import "./libgogo/_obj/libgogo"

func GenerateForStart(item *libgogo.Item, ed *ExpressionDescriptor) {
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
        PrintLabelWrapped(ed, 1 /*local*/, 0 /*negative*/, "END");
    }
    PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "END");

    // Positive branch starts after this label, thus insert last remaining 
    // positive label (if available) here
    if ed.T != 0 {
        PrintLabelWrapped(ed, 1 /*local*/, 1 /*positive*/, "END");
    }
}

func GenerateForEnd(ed *ExpressionDescriptor) {
    if ed.ForPost != 0 {
        PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "EXTENDED_BODY");
    } else {
        PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "EXPR_START");
    }
    PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "END");
}

func GenerateForBodyExtended(ed *ExpressionDescriptor) {
    PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "EXTENDED_BODY");
}

func GenerateForBody(ed *ExpressionDescriptor) {
    if ed.ForPost != 0 {
        if ed.ForExpr != 0 {
            PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "EXPR_START");
        }
        PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "BODY");
    }
}

func GenerateForBodyJump(ed *ExpressionDescriptor) {
    PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "BODY");
}

func GenerateExpressionStart(ed *ExpressionDescriptor) {
    PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "EXPR_START");
}

func GenerateContinue(ed *ExpressionDescriptor) {
    if ed.ForPost != 0 {
        PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "EXTENDED_BODY");
    } else {
        PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "EXPR_START");
    }
}

