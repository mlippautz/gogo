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

//
// Basic types for reference
//
var uint64_t *libgogo.TypeDesc = nil;
var byte_t *libgogo.TypeDesc = nil;
var string_t *libgogo.TypeDesc = nil;

//
// Nil pointer for reference
//
var nilPtr *libgogo.ObjectDesc = nil;

//
// Initializes the global symbol table (objects and types)
//
func InitSymbolTable() {
    //Default data types
    uint64_t = libgogo.NewType("uint64", "", 0, 8, nil);
    GlobalTypes = libgogo.AppendType(uint64_t, GlobalTypes);
    byte_t = libgogo.NewType("byte", "", 0, 1, nil);
    GlobalTypes = libgogo.AppendType(byte_t, GlobalTypes);
    string_t = libgogo.NewType("string", "", 0, 16, nil);
    GlobalTypes = libgogo.AppendType(string_t, GlobalTypes);

    //Default objects
    nilPtr = libgogo.NewObject("nil", "", libgogo.CLASS_VAR);
    libgogo.SetObjType(nilPtr, nil);
    libgogo.FlagObjectTypeAsPointer(nilPtr); //nil is a pointer to no specified type (universal)
    GlobalObjects = libgogo.AppendObject(nilPtr, GlobalObjects);
}

//
// Prints the global symbol table (objects and types)
// subject to the current debug level
//
func PrintGlobalSymbolTable() {
    var temp uint64;
    temp = CheckDebugLevel(100);
    if  temp == 1 {
        libgogo.PrintString("\nGlobal symbol table:\n");
        libgogo.PrintString("--------------------\n");
        libgogo.PrintTypes(GlobalTypes);
        libgogo.PrintObjects(GlobalObjects);
    }
}

//
// Prints the local symbol table (objects only) subject
// to the current debug level
//
func PrintLocalSymbolTable() {
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
}

//
// Determines whether a type is a basic type (uint64, byte, string)
// Returns 1 if t is a basic type, or 0 if not
//
func IsBasicDataType(t *libgogo.TypeDesc) uint64 {
    var retVal uint64 = 0;
    if (t == uint64_t) || (t == byte_t) || (t == string_t) {
        retVal = 1;
    }
    return retVal;
}

//
// Checks for undefined types which were forward declared
//
func UndefinedForwardDeclaredTypeCheck() {
    var temptype *libgogo.TypeDesc;
    var tempstring string;
    temptype = libgogo.GetFirstForwardDeclType(GlobalTypes);
    if temptype != nil {
        tempstring = libgogo.GetTypeName(temptype);
        SymbolTableError("undefined", "", "type", tempstring);
    }
}

//
// Creates a new type with the specified name to the symbol table
// The newly created type is not added to the global symbol table,
// but a flag is returned whether this has to be done later when
// the whole type is specified using AddType
//
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

//
// Adds a newly created type built by NewType when all of its fields
// have been added. The parameter dontAddType specifies whether the
// the type actually needs to be added to the global symbol tables
//
func AddType(dontAddType uint64) {
    if dontAddType == 0 {
        GlobalTypes = libgogo.AppendType(CurrentType, GlobalTypes);
    }
}

//
// Adds a fields with the speicified name to a new type created by
// NewType. The field's type has to be set using SetCurrentObjectType
//
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

//
// Flags the type the current object or field to be a pointer and
// not the type itself
//
func SetCurrentObjectTypeToPointer() {
    libgogo.FlagObjectTypeAsPointer(CurrentObject); //Type is pointer
}

//
// Sets the type of the current object or field using type
// [arraydim]packagename.typename where arraydim and packagename
// are optional. In case of arraydim != 0, a new type is created
// to represent an array type of the specified size
//
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

//
// Adds a new variable to the global symbol table with the specified
// name. The variable type has to be set using SetCurrentObjectType
//
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

//
// When called when the reaching the end of a function during parsing, the
// local symbol table is purged and optional debug information about the
// symbol table before purging are printed
//
func EndOfFunction() {
    PrintLocalSymbolTable(); //Print local symbol table
    LocalObjects = nil; //Delete local objects
}