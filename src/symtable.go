// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import "./libgogo/_obj/libgogo"

//
// List of global objects and declared types
//
var GlobalObjects *libgogo.ObjectDesc = nil;
var GlobalTypes *libgogo.TypeDesc = nil;

//
// List of function-local objects
//
var LocalObjects *libgogo.ObjectDesc = nil;

//
// List of currently processed objects and types
//
var CurrentType *libgogo.TypeDesc;
var CurrentObject *libgogo.ObjectDesc;

func NewType(name string) uint64 {
    var dontAddType uint64 = 0;
    var tempType *libgogo.TypeDesc;
    var temp byte;
    tempType = libgogo.GetType(name, CurrentPackage, GlobalTypes, 1);
    if  tempType != nil { //Check for duplicates
        temp = libgogo.IsForwardDecl(tempType);
        if temp != 0 { //Separate handling of forward declarations => unset forward declaration flag
            libgogo.UnsetForwardDecl(tempType);
            CurrentType = tempType;
            dontAddType = 1;
        } else { //Real duplicate
            SymbolTableError("duplicate type", name, "in package", CurrentPackage);
            CurrentType = tempType; //For stability/sanity of parser when continuing
            dontAddType = 1; //For stability/sanity of parser when continuing
        }
    } else {
        CurrentType = libgogo.NewType(name, CurrentPackage, 0, 0, nil);
    }
    return dontAddType;
}

func AddType(dontAddType uint64) {
    if dontAddType == 0 {
        GlobalTypes = libgogo.AppendType(CurrentType, GlobalTypes);
    }
}

func AddStructField(fieldname string) {
    var temp uint64;
    temp = libgogo.HasField(fieldname, CurrentType);
    if temp != 0 {
        SymbolTableError("duplicate", "", "field", fieldname);
    } else {
        CurrentObject = libgogo.NewObject(fieldname, "", libgogo.CLASS_FIELD); //A field has no package name
        libgogo.AddFields(CurrentObject, CurrentType);
    }
}

func SetCurrentObjectTypeToPointer() {
    libgogo.FlagObjectTypeAsPointer(CurrentObject); //Type is pointer
}

func SetCurrentObjectType(typename string, packagename string, arraydim uint64) {
    var basetype *libgogo.TypeDesc;
    var temptype *libgogo.TypeDesc;
    var CurrentTypeName string;
    var tempstr string = "";
    var temp byte;
    var boolFlag uint64;

    if InsideFunctionVarDecl == 0 {
        if InsideStructDecl == 1 { //Allow types in struct declarations which are already forward declared
            basetype = libgogo.GetType(typename, packagename, GlobalTypes, 1);
        } else {
            basetype = libgogo.GetType(typename, packagename, GlobalTypes, 0);
        }
        if basetype == nil {
            if InsideStructDecl == 1 {
                CurrentTypeName = libgogo.GetTypeName(CurrentType);
                boolFlag = libgogo.StringCompare(typename, CurrentTypeName);
                if boolFlag == 0 {
                    temp = libgogo.IsPointerType(CurrentObject);
                    if temp == 1 { //Allow pointer to own type
                        basetype = CurrentType;
                    } else {
                        SymbolTableError("A type cannot contain itself,", "", "type", typename);
                    }
                } else { //Forward declaration
                    basetype = libgogo.NewType(typename, packagename, 1, 0, nil);
                    GlobalTypes = libgogo.AppendType(basetype, GlobalTypes); //Add forward declared type to global list
                }
            } else {
                libgogo.StringAppend(&tempstr, packagename);
                libgogo.CharAppend(&tempstr, '.');
                libgogo.StringAppend(&tempstr, typename);
                SymbolTableError("Unknown", "", "type", tempstr);
            }
        }
        if arraydim == 0 { //No array
            libgogo.SetObjType(CurrentObject, basetype);
        } else { //Array
            if basetype != nil {
                CurrentTypeName = libgogo.GetTypeName(basetype);
                libgogo.StringAppend(&tempstr, CurrentTypeName);
            }
            libgogo.StringAppend(&tempstr, "Array");
            CurrentTypeName = libgogo.IntToString(arraydim); //Reuse CurrentTypeName as string representation of arraydim
            libgogo.StringAppend(&tempstr, CurrentTypeName);
            temptype = libgogo.NewType(tempstr, packagename, 0, arraydim, basetype);
            libgogo.SetObjType(CurrentObject, temptype); //Don't add array type to global list to avoid duplicate type name errors
        }
    }
}

func NewVariable(name string) {
    var TempObject *libgogo.ObjectDesc;
    CurrentObject = libgogo.NewObject(name, CurrentPackage, libgogo.CLASS_VAR);
    if InsideFunction == 0 { //Global objects
        TempObject = libgogo.GetObject(name, CurrentPackage, GlobalObjects);
        if TempObject != nil {
            SymbolTableError("duplicate", "global", "identifier", name);
        }
        GlobalObjects = libgogo.AppendObject(CurrentObject, GlobalObjects);
    } else { //Function-local objects
        TempObject = libgogo.GetObject(name, CurrentPackage, LocalObjects);
        if TempObject != nil {
            SymbolTableError("duplicate", "local", "identifier", name);
        }
        LocalObjects = libgogo.AppendObject(CurrentObject, LocalObjects);
    }
}

func EndOfFunction() {
    var temp uint64;
    temp = CheckDebugLevel(100);
    if temp == 1 { //Function-local symbol table
        libgogo.PrintString("\nFunction-local symbol table until line ");
        libgogo.PrintNumber(fileInfo[curFileIndex].lineCounter);
        libgogo.PrintString(" of ");
        libgogo.PrintString(fileInfo[curFileIndex].filename);
        libgogo.PrintString(":\n----------------------------------------------------------------------------\n");
        libgogo.PrintObjects(LocalObjects);
    }
    LocalObjects = nil; //Delete local objects
}
