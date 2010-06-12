// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

func ZeroParameterCheck(FunctionCalled *libgogo.TypeDesc) {
    var FullFunctionName string;
    if FunctionCalled.Len == 0 { //Check if function expects parameters
        FullFunctionName = "";
        libgogo.StringAppend(&FullFunctionName, FunctionCalled.PackageName);
        libgogo.CharAppend(&FullFunctionName, '.');
        libgogo.StringAppend(&FullFunctionName, FunctionCalled.Name);
        SymbolTableError("Function expects", "no", "parameters:", FullFunctionName);
    }
}

func ParameterCheck_More(FunctionCalled *libgogo.TypeDesc, paramCount uint64) {
    var FullFunctionName string;
    var tempString string;
    if FunctionCalled.Len > paramCount { //Compare number of actual parameters
        FullFunctionName = "";
        libgogo.StringAppend(&FullFunctionName, FunctionCalled.PackageName);
        libgogo.CharAppend(&FullFunctionName, '.');
        libgogo.StringAppend(&FullFunctionName, FunctionCalled.Name);
        tempString = libgogo.IntToString(FunctionCalled.Len);
        SymbolTableError("Expecting", tempString, "parameters (more than the actual ones) for function", FullFunctionName);
    }
}

func ParameterCheck_Less(FunctionCalled *libgogo.TypeDesc, paramCount uint64) {
    var FullFunctionName string;
    var tempString string;
    if FunctionCalled.Len < paramCount { //Compare number of actual parameters
        FullFunctionName = "";
        libgogo.StringAppend(&FullFunctionName, FunctionCalled.PackageName);
        libgogo.CharAppend(&FullFunctionName, '.');
        libgogo.StringAppend(&FullFunctionName, FunctionCalled.Name);
        tempString = libgogo.IntToString(FunctionCalled.Len);
        SymbolTableError("Expecting", tempString, "parameters (less than the actual ones) for function", FullFunctionName);
    }
}

func PrintActualFunctionCall(FunctionCalled *libgogo.TypeDesc, TotalLocalVariableSize uint64, TotalParameterSize uint64) {
    if FunctionCalled.ForwardDecl == 1 {
        PrintFunctionCall(FunctionCalled.PackageName, FunctionCalled.Name, TotalLocalVariableSize, 1);
    } else {
        PrintFunctionCall(FunctionCalled.PackageName, FunctionCalled.Name, TotalParameterSize + TotalLocalVariableSize, 0);
    }
}

func GetReturnItem(FunctionCalled *libgogo.TypeDesc, TotalLocalVariableSize uint64, TotalParameterSize uint64) *libgogo.Item {
    var ReturnObject *libgogo.ObjectDesc;
    var ReturnItem *libgogo.Item;
    ReturnObject = libgogo.GetObject("return value", "", FunctionCalled.Fields); //Find return value
    if ReturnObject == nil {
        ReturnItem = nil;
    } else {
        if FunctionCalled.ForwardDecl == 1 {
            ReturnItem = ObjectToStackParameter(ReturnObject, FunctionCalled, TotalLocalVariableSize);
        } else {
            ReturnItem = ObjectToStackParameter(ReturnObject, FunctionCalled, TotalParameterSize + TotalLocalVariableSize);
        }
    }
    return ReturnItem;
}

func AddArtificialParameterIfNecessary(FunctionCalled *libgogo.TypeDesc, ExprItem *libgogo.Item) {
    var TempObject *libgogo.ObjectDesc;
    var boolFlag uint64;
    if (FunctionCalled.ForwardDecl == 1) && (FunctionCalled.Base == nil) { //Create artificial parameter from expression (based on the latter's type) if the function is called the first time without being declared
        TempObject = libgogo.NewObject("Artificial parameter", "", libgogo.CLASS_PARAMETER);
        TempObject.ObjType = ExprItem.Itemtype; //Derive type from expression
        TempObject.PtrType = ExprItem.PtrType; //Derive pointer type from expression
        if boolFlag != 0 { //& in expression forces pointer type
            if TempObject.PtrType == 0 {
                TempObject.PtrType = 1;
            } else {
                SymbolTableError("& operator on pointer type not allowed,", "", "type: pointer to", ExprItem.Itemtype.Name);
            }
        }
        libgogo.AddParameters(TempObject, FunctionCalled); //Add a new, artificial parameter
    }
}

func AddArtificialReturnValueIfNecessary(FunctionCalled *libgogo.TypeDesc, ReturnValue *libgogo.Item, ForwardDeclExpectedReturnType *libgogo.TypeDesc, ForwardDeclExpectedReturnPtrType uint64) *libgogo.Item {
    var TotalLocalVariableSize uint64;
    var TempObject *libgogo.ObjectDesc;
    if (FunctionCalled.ForwardDecl == 1) && (FunctionCalled.Base == nil) { //Create artifical return value if function is called the first time
        if ForwardDeclExpectedReturnType != nil { //Return type expected
            TempObject = libgogo.NewObject("return value", "", libgogo.CLASS_PARAMETER); //Create artificial return value
            TempObject.ObjType = ForwardDeclExpectedReturnType;
            TempObject.PtrType = ForwardDeclExpectedReturnPtrType;
            libgogo.AddParameters(TempObject, FunctionCalled); //Add a new, artificial return value
            FunctionCalled.Len = FunctionCalled.Len - 1; //Don't count parameter as input parameter
            TotalLocalVariableSize = libgogo.GetAlignedObjectListSize(LocalObjects); //Take local variable size into consideration for offset below
            ReturnValue = ObjectToStackParameter(TempObject, FunctionCalled, TotalLocalVariableSize);
        } else { //No return type expected
            ReturnValue = nil;
        }
    }
    return ReturnValue;
}

func CorrectArtificialParameterIfNecessary(FunctionCalled *libgogo.TypeDesc, ParameterLHSObject *libgogo.ObjectDesc, ExprItem *libgogo.Item) {
    if FunctionCalled.ForwardDecl == 1 { //Parameter type up-conversion
        if (ParameterLHSObject.ObjType == byte_t) && (ParameterLHSObject.PtrType == 0) && (ExprItem.Itemtype == uint64_t) && (ExprItem.PtrType == 0) { //If previous forward declaration of this parameter was of type byte, it is possible that is was a byte constant and is now of type uint64 => set to type uint64 in declaration
            ParameterLHSObject.ObjType = uint64_t;
        }
        if (ParameterLHSObject.ObjType == nil) && (ParameterLHSObject.PtrType == 1) && (ExprItem.PtrType == 1) { //If previous forward declaration of this parameter was of type unspecified pointer, it was nil and is now of type *rhs_type => set to type of RHS in declaration
            ParameterLHSObject.ObjType = ExprItem.Itemtype;
        }
    }
}
