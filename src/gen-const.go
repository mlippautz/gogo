// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

func CreateIntegerConstant(value uint64, item *libgogo.Item) {
    if (value <= 255) { //Value fits into byte_t
        libgogo.SetItem(item, libgogo.MODE_CONST, byte_t, 0, value, 0, 0); //Constant item
    } else { //Value does not fit into byte_t
        libgogo.SetItem(item, libgogo.MODE_CONST, uint64_t, 0, value, 0, 0); //Constant item
    }
}

func CreateStringConstant(str string, item *libgogo.Item) {
    var boolFlag uint64;
    var doneFlag uint64;
    var tempType *libgogo.TypeDesc;
    var tempByteArray *libgogo.ObjectDesc;
    var tempString *libgogo.ObjectDesc;
    var startAddress uint64;
    var LHSItem *libgogo.Item;
    var RHSItem *libgogo.Item;
    var s string;

    boolFlag = libgogo.StringLength(str); //Compute length
    tempType = libgogo.NewType("byteArray", ".internal", 0, boolFlag + 1, byte_t); //Create byte array type of according length (including trailing 0) to be able to address the characters
    tempByteArray = libgogo.NewObject("tempByteArray", ".internal", libgogo.CLASS_VAR); //Create object of previously declared byte array type
    tempByteArray.ObjType = tempType;
    libgogo.AppendObject(tempByteArray, GlobalObjects); //Add byte array to global objects
    startAddress = libgogo.GetObjectOffset(tempByteArray, GlobalObjects); //Calculate buffer start address
    
    SwitchOutputToDataSegment(); //Place content of byte array in data segment
    s = "String buffer start ('";
    libgogo.StringAppend(&s, str);
    libgogo.StringAppend(&s, "')");
    GenerateComment(s); //Output string in comment
    for doneFlag = 0; doneFlag < boolFlag; doneFlag = doneFlag + 1 { //Set values in data segment accordingly
        PutDataByte(startAddress + doneFlag, str[doneFlag]);
    }
    PutDataByte(startAddress + doneFlag, 0); //Add trailing 0
    GenerateComment("String buffer end");
    //SwitchOutputToCodeSegment(); //Reset to default output
    
    tempString = libgogo.NewObject("tempString", ".internal", libgogo.CLASS_VAR); //Create object for actual string
    tempString.ObjType = string_t;
    libgogo.AppendObject(tempString, GlobalObjects); //Add string to global objects
    
    SwitchOutputToInitCodeSegment(); //Initialize strings globally
    GenerateComment("Assign byte buffer to new string constant start");
    LHSItem = libgogo.NewItem();
    VariableObjectDescToItem(tempString, LHSItem, 1); //Global variable
    LHSItem.Itemtype = byte_t; //Set appropriate type for byte (array) pointer
    LHSItem.PtrType = 1;
    RHSItem = libgogo.NewItem();
    VariableObjectDescToItem(tempByteArray, RHSItem, 1); //Global variable
    RHSItem.Itemtype = byte_t; //Set appropriate type for byte (array) to be pointed at
    GenerateAssignment(LHSItem, RHSItem, 1); //tempString{first qword} = &tempByteArray
    GenerateComment("Assign byte buffer to new string constant end");
    GenerateComment("Assign string length to new string constant start");
    LHSItem = libgogo.NewItem();
    VariableObjectDescToItem(tempString, LHSItem, 1); //Global variable
    LHSItem.A = LHSItem.A + 8; //Access second qword of string to place length in there
    LHSItem.Itemtype = uint64_t; //Set appropriate type for string length
    RHSItem = libgogo.NewItem();
    libgogo.SetItem(RHSItem, libgogo.MODE_CONST, uint64_t, 0, boolFlag, 0, 0); //Constant item containing the length of the string
    GenerateAssignment(LHSItem, RHSItem, 0); //tempString{second qword} = boolFlag (string length)
    GenerateComment("Assign string length to new string constant end");
    SwitchOutputToCodeSegment(); //Reset to default output
    
    VariableObjectDescToItem(tempString, item, 1); //Global variable to return (reference to properly initialized, global string)
}
