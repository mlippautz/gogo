// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.


package main

import "./libgogo/_obj/libgogo"

func GenerateIfStart(item *libgogo.Item, ed ExpressionDescriptor) {
    var labelString string;
    labelString = GenerateIfLabel(ed.Prefix, ed.GlobalCounter, ed.ExpressionDepth-1, "END");
    PrintLabel(labelString);
    labelString = GenerateIfLabel(ed.Prefix, ed.GlobalCounter, 0, "END");
    PrintJump("JMP",labelString);
    labelString = GenerateIfLabel(ed.Prefix, ed.GlobalCounter, 0, "OK");
    PrintLabel(labelString);
}

func GenerateIfEnd(item *libgogo.Item, ed ExpressionDescriptor) {
    var labelString string;
    labelString = GenerateIfLabel(ed.Prefix, ed.GlobalCounter, 0, "END");
    PrintLabel(labelString);
}

func GenerateIfLabel(prefix string, global uint64, local uint64, label string) string {
    var str string;
    var tmpStr string;
    libgogo.StringAppend(&str, prefix);
    libgogo.StringAppend(&str, "_");
    tmpStr = libgogo.IntToString(global);
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
