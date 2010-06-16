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

type ParamDesc struct {
    Name string;
    TypePackage string;
    TypeName string;
    Ptr uint64;
};

//
// Function fetches the next line from the assembly file.
//
func GetLine(ld *LineDesc) {
    var line string;
    var i uint64 =0;
    var singleChar byte;
    var cmp uint64;

    for singleChar = GetCharWrapped(); (i<10) && (singleChar!=10) && (singleChar!=0); singleChar = GetCharWrapped() {
        libgogo.CharAppend(&line, singleChar);
        i = i +1;
    }
    cmp = libgogo.StringCompare(line, "  //--- ##");
    if cmp == 0 {
        singleChar = GetCharWrapped(); // should be number; the case of fixing
        ld.NeedsFix = libgogo.ToIntFromByte(singleChar);
        ld.NeedsFix = ld.NeedsFix - 48;
        ld.PackageName = "";
        ld.FunctionName = "";

        GetCharWrapped(); // Abolish '#'
        GetCharWrapped(); // Abolish '#'

        for singleChar = GetCharWrapped(); singleChar != 194 ; singleChar = GetCharWrapped() {
            libgogo.CharAppend(&ld.PackageName, singleChar);
        }
        GetCharWrapped(); // Abolish second part of dot (183)
        for singleChar = GetCharWrapped(); (singleChar != '#'); singleChar = GetCharWrapped() {
            libgogo.CharAppend(&ld.FunctionName, singleChar);
        }

        if (singleChar != 10) && (singleChar != 0) {
            for ; (singleChar != 0) && (singleChar != 10); singleChar = GetCharWrapped() {
                libgogo.CharAppend(&line, singleChar);
            }
        }
    } else {
        for ; (singleChar != 0) && (singleChar != 10); singleChar = GetCharWrapped() {
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
    var symtoken string = "";
    var i uint64;
    var len uint64;
    len = libgogo.StringLength(ld.Line);
    for i = ld.Offset; (i < len) && (ld.Line[i] != ','); i=i+1 {
        if (ld.Line[i] != '/') {
            libgogo.CharAppend(&symtoken, ld.Line[i]);
        }
    }
    ld.Offset = i+1;
    return symtoken;
}


func ParseSymTblLine(ld *LineDesc) {
    var symtype string;
    var strCmp uint64;

    symtype = GetNextSymToken(ld);
    strCmp = libgogo.StringCompare(symtype, "TYPE");
    if strCmp == 0 {
        ParseSymTblType(ld);
    }
    strCmp = libgogo.StringCompare(symtype, "FUNC");
    if strCmp == 0 {
        ParseSymTblFunc(ld);
    }
}

func IsDefaultType(pkgType string) uint64 {
    var strCmp uint64;
    var retValue uint64 = 0;
    strCmp = libgogo.StringCompare(pkgType, "·uint64");
    if strCmp == 0 {
        retValue = 1;
    }
    strCmp = libgogo.StringCompare(pkgType, "·byte");
    if strCmp == 0 {
        retValue = 1;
    }
    strCmp = libgogo.StringCompare(pkgType, "·string");
    if strCmp == 0 {
        retValue = 1;
    }
    strCmp = libgogo.StringCompare(pkgType, "·bool");
    if strCmp == 0 {
        retValue = 1;
    }
    return retValue;
}

func GetPackageName(pkgType string) string {
    var retStr string;
    var i uint64;
    for i = 0; pkgType[i] != 194 ; i=i+1 {
        libgogo.CharAppend(&retStr, pkgType[i]);
    }
    return retStr;
}

func GetFuncName(pkgType string) string {
    var retStr string;
    var i uint64;
    var strLen uint64;
    strLen = libgogo.StringLength(pkgType);
    for i = 0; pkgType[i] != 183; i = i +1 {
    }
    i = i+1;
    for ; i < strLen; i = i +1 {
        libgogo.CharAppend(&retStr, pkgType[i]);
    }
    return retStr;
}

func ParseSymbolParam(pd *ParamDesc, paramType string) {
    var i uint64 = 0;
    var strLen uint64;

    pd.Name = "";
    pd.TypePackage = "";
    pd.TypeName = "";
    pd.Ptr = 0;

    strLen = libgogo.StringLength(paramType);
    if strLen > 0 {
        for i = 0; paramType[i] != ':' ; i=i+1 {
            libgogo.CharAppend(&pd.Name, paramType[i]);
        }
        i=i+1;
        if paramType[i] == '*' {
            pd.Ptr = 1;
            i=i+1;
        }
        for ;paramType[i] != 194; i=i+1 {
            libgogo.CharAppend(&pd.TypePackage, paramType[i]);
        }
        for i=i+2;i < strLen; i=i+1 {
            libgogo.CharAppend(&pd.TypeName, paramType[i]);
        }
    }
}

//
// Function parsing the main part of a function in a symbol table and adding it 
// to the global symbol table.
//
func ParseSymTblFunc(ld *LineDesc) {
    var pd ParamDesc;
    var func_t *libgogo.TypeDesc = nil;
    var some_t *libgogo.TypeDesc = nil;
    var fwdStr string;
    var fwdNum uint64;
    var pkgFunc string;
    var funcName string;
    var pkgName string;
    var ind uint64;
    var paramType string;
    var tmpParam *libgogo.ObjectDesc;

    fwdStr = GetNextSymToken(ld);
    fwdNum = libgogo.StringToInt(fwdStr);
    pkgFunc = GetNextSymToken(ld);
    ind = libgogo.StringCompare(pkgFunc, "main·init");
    if ind != 0 {
        pkgName = GetPackageName(pkgFunc);
        funcName = GetFuncName(pkgFunc);
        func_t = NewFunction(funcName, pkgName, fwdNum);
        paramType = GetNextSymToken(ld);
        ind = libgogo.StringLength(paramType);
        for ; ind != 0 ; {
            ParseSymbolParam(&pd, paramType);
            tmpParam = libgogo.NewObject(pd.Name, "", libgogo.CLASS_PARAMETER);
            some_t = libgogo.GetType(pd.TypeName, pd.TypePackage, GlobalTypes, 1);
            if some_t != nil {
                tmpParam.ObjType = some_t;
                tmpParam.PtrType = pd.Ptr;
                libgogo.AddParameters(tmpParam, func_t);
            } else {
                LinkError("unable to find type '",pd.TypePackage,"·",pd.TypeName,"'.");
            }
            paramType = GetNextSymToken(ld);
            ind = libgogo.StringLength(paramType);
        }
    }

}

//
// Function parsing the main part of a type in a symbol table and adding it 
// to the global symbol table.
//
func ParseSymTblType(ld *LineDesc) {
    var some_t *libgogo.TypeDesc = nil;
    var fwdStr string;
    var fwdNum uint64;
    var pkgType string;
    var pkgName string;
    var typeName string;
    var sizeStr string;
    var sizeNum uint64;
    var alignStr string;
    var alignNum uint64;
    var ind uint64;
    var tmp1 uint64;
    var tmp2 uint64;
    var tmpStr1 string;
    var tmpStr2 string;
    // maybe some more flags, indicators, temps, strings, numbers, ... ?!? ;)

    fwdStr = GetNextSymToken(ld);
    fwdNum = libgogo.StringToInt(fwdStr);
    pkgType = GetNextSymToken(ld);
    sizeStr = GetNextSymToken(ld);
    sizeNum = libgogo.StringToInt(sizeStr);
    alignStr = GetNextSymToken(ld);
    alignNum = libgogo.StringToInt(alignStr);
    ind = IsDefaultType(pkgType);
    if ind == 0 { // non-default type, try to add
        pkgName = GetPackageName(pkgType);
        typeName = GetFuncName(pkgType);
        some_t = libgogo.GetType(typeName, pkgName, GlobalTypes, 1);
        if some_t == nil {
            some_t = libgogo.NewType(typeName, pkgName, fwdNum, sizeNum, nil);
            GlobalTypes = libgogo.AppendType(some_t, GlobalTypes);
        } else {
            /* All kings of fwd declaration combinations need to be checked */
            if (some_t.ForwardDecl == 0) && (fwdNum == 0) {
                tmp1 = libgogo.GetTypeSize(some_t);
                tmp2 = libgogo.GetTypeSizeAligned(some_t);
                if (sizeNum == tmp1) && (alignNum == tmp2) {
                    LinkWarn("duplicated type (", pkgType, "). sizes matched.", "", "");
                } else {
                    libgogo.StringAppend(&tmpStr1, "new: ");
                    libgogo.StringAppend(&tmpStr1, sizeStr);
                    libgogo.StringAppend(&tmpStr1, "; defined: ");
                    tmpStr2 = libgogo.IntToString(tmp1);
                    libgogo.StringAppend(&tmpStr1, tmpStr2);
                    LinkError("duplicate type: ", pkgType, ". Incompatible sizes: ",tmpStr1, ""); 
                }
            } 
            if (some_t.ForwardDecl == 1) && (fwdNum == 0) {
                /* Fix the params of the fwd. declared type */
                some_t.ForwardDecl = 0;
                some_t.Len = sizeNum;
            } 
            if (some_t.ForwardDecl == 0) && (fwdNum == 1) {
                ; // skip since this is useless
            }
            if (some_t.ForwardDecl == 1) && (fwdNum == 1) {
                ; // skip since this is useless
            }
        }
    }
}

func GetParameterSize(pkgName string, funcName string) uint64 {
    var size uint64;
    var func_t *libgogo.TypeDesc = nil;
    func_t = libgogo.GetType(funcName, pkgName, GlobalFunctions, 1);
    if func_t != nil {
        size = libgogo.GetAlignedObjectListSize(func_t.Fields);
    } else {
        LinkError("unable to look up function '", pkgName, "·", funcName, "' in symtable.");
    }
    return size;
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
        strLen = libgogo.StringLength(ld.Line);
        size = GetParameterSize(ld.PackageName, ld.FunctionName);
        for i = 0; ld.Line[i] != '$' ; i = i +1 {
            libgogo.CharAppend(&newLine, ld.Line[i]);
        }
        libgogo.CharAppend(&newLine, ld.Line[i]);
        for i = i+1 /*Dismiss '-'*/ ;ld.Line[i]!=',';i=i+1 {
            libgogo.CharAppend(&numstr, ld.Line[i]);
        }
        oldsize = libgogo.StringToInt(numstr);
        size = size + oldsize;
        numstr = libgogo.IntToString(size);
        libgogo.StringAppend(&newLine, numstr);
        for ; i < strLen; i = i +1 {
            libgogo.CharAppend(&newLine, ld.Line[i]);
        }
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

    ResetToken();
    GetLine(&ld);
    for ;tok.id != TOKEN_EOS;{
        strCmp = libgogo.StringCompare("//Symbol table:", ld.Line);
        if strCmp == 0 {
            symtable = 1;
            GetLine(&ld); // Proceed to next line
        } 
        if symtable != 0 { // Parse symtable entries
            strCmp = libgogo.StringCompare("//End Symbol table", ld.Line);
            if strCmp == 0 {
                symtable = 0;
            } else {
                ParseSymTblLine(&ld);
            }
        } else { // Parse normal lines and fix everything
            if ld.NeedsFix != 0 {
                GetLine(&ld);
                newLine = FixOffset(&ld);
                libgogo.PrintString(newLine);
                libgogo.PrintString("\n");
            } else {
                strCmp = libgogo.StringCompare(ld.Line, "__UNLINKED_CODE");
                if strCmp != 0 {
                    libgogo.PrintString(ld.Line);
                    libgogo.PrintString("\n");
                }
            }

        }
        GetLine(&ld);
    }
}
