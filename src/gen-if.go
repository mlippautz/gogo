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
    var stacksize uint64;

    labelString = GenerateSubLabel(ed,1,"END");

    if ed.Not == 0 {
        jmp = GetJump(item.C, 0);
    } else {
        labelString = GenerateSubLabel(ed,1,"END");
        jmp = GetJump(item.C, 1);
    }
    PrintJump(jmp, labelString);

    // Important: Since last jump is a positive one, we have to start with the
    // negative path
    for stacksize = libgogo.GetStackItemCount(&ed.FS); stacksize > 0 ; stacksize = libgogo.GetStackItemCount(&ed.FS) {
        PrintLabelWrapped(ed, 1 /*local*/, 0 /*negative*/, "END");
        libgogo.Pop(&ed.FS);
        libgogo.Pop(&ed.FDepthS);
    }
    PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "END");

    // Positive branch starts after this label, thus insert last remaining 
    // positive label (if available) here
    for stacksize = libgogo.GetStackItemCount(&ed.TS); stacksize > 0 ; stacksize = libgogo.GetStackItemCount(&ed.TS) {
        PrintLabelWrapped(ed, 1 /*local*/, 1 /*positive*/, "END");
        libgogo.Pop(&ed.TS);
        libgogo.Pop(&ed.TDepthS);
    }

    item.C = 0;
    FreeRegisterIfRequired(item);
}

func GenerateIfEnd(ed *ExpressionDescriptor) {
    PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "END");
    PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "ELSE_END");
}

func GenerateElseStart(ed *ExpressionDescriptor) {
    PrintJumpWrapped("JMP", ed, 0 /*global*/, 0 /*unused*/, "ELSE_END");
    PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "END");
}

func GenerateElseEnd(ed *ExpressionDescriptor) {
    PrintLabelWrapped(ed, 0 /*global*/, 0 /*unused*/, "ELSE_END");
}

