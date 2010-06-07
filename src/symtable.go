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
var GlobalFunctions *libgogo.TypeDesc = nil;

//
// List of function-local objects
//
var LocalObjects *libgogo.ObjectDesc = nil;

//
// List of currently processed objects and types
//
var CurrentType *libgogo.TypeDesc;
var CurrentObject *libgogo.ObjectDesc;
var CurrentFunction *libgogo.TypeDesc;

//
// Basic types for reference
//
var uint64_t *libgogo.TypeDesc = nil;
var byte_t *libgogo.TypeDesc = nil;
var string_t *libgogo.TypeDesc = nil;
var bool_t *libgogo.TypeDesc = nil;

//
// Nil pointer for reference
//
var nilPtr *libgogo.ObjectDesc = nil;

//
// Initializes the global symbol table (objects and types)
//
func InitSymbolTable() {
    var main_init_fcn *libgogo.TypeDesc;

    //Default data types
    uint64_t = libgogo.NewType("uint64", "", 0, 8, nil);
    GlobalTypes = libgogo.AppendType(uint64_t, GlobalTypes);
    byte_t = libgogo.NewType("byte", "", 0, 1, nil);
    GlobalTypes = libgogo.AppendType(byte_t, GlobalTypes);
    string_t = libgogo.NewType("string", "", 0, 16, byte_t); //Hybrid string type: 16 byte value on the one hand, byte array on the other in order to allow item/char access
    GlobalTypes = libgogo.AppendType(string_t, GlobalTypes);
    bool_t = libgogo.NewType("bool", "", 0, 8, nil);
    GlobalTypes = libgogo.AppendType(bool_t, GlobalTypes);

    //Default objects
    nilPtr = libgogo.NewObject("nil", "", libgogo.CLASS_VAR);
    nilPtr.ObjType = nil;
    nilPtr.PtrType = 1; //nil is a pointer to no specified type (universal)
    GlobalObjects = libgogo.AppendObject(nilPtr, GlobalObjects);
    
    //Default functions
    main_init_fcn = libgogo.NewType("init", "main", 0, 0, nil); //main.init (avoid external redeclaraction)
    GlobalFunctions = libgogo.AppendType(main_init_fcn, GlobalFunctions);
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
        libgogo.PrintFunctions(GlobalFunctions);
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
        libgogo.PrintObjects(CurrentFunction.Fields);
        libgogo.PrintString("--- End of parameters, begin of local variables ---\n");
        libgogo.PrintObjects(LocalObjects);
    }
}

//
// Checks for undefined types which were forward declared
//
func UndefinedForwardDeclaredTypeCheck() {
    var temptype *libgogo.TypeDesc;
    if Compile != 0 {
		temptype = libgogo.GetFirstForwardDeclType(GlobalTypes);
		if temptype != nil {
		    SymbolTableError("undefined", "", "type", temptype.Name);
		}
        temptype = libgogo.GetFirstForwardDeclType(GlobalFunctions);
		if temptype != nil {
		    SymbolTableError("undefined", "", "function", temptype.Name);
		}
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
    if Compile != 0 {
		tempType = libgogo.GetType(name, CurrentPackage, GlobalTypes, 1);
		if  tempType != nil { //Check for duplicates
		    if tempType.ForwardDecl != 0 { //Separate handling of forward declarations => unset forward declaration flag
		        tempType.ForwardDecl = 0;
		        CurrentType = tempType;
		        dontAddType = 1;
		    } else { //Real duplicate
		        SymbolTableError("duplicate type", name, "in package", CurrentPackage);
		    }
		} else {
		    CurrentType = libgogo.NewType(name, CurrentPackage, 0, 0, nil);
		}
    }
    return dontAddType;
}

//
// Adds a newly created type built by NewType when all of its fields
// have been added. The parameter dontAddType specifies whether the
// the type actually needs to be added to the global symbol tables
//
func AddType(dontAddType uint64) {
    if (Compile != 0) && (dontAddType == 0) {
        GlobalTypes = libgogo.AppendType(CurrentType, GlobalTypes);
    }
}

//
// Adds a fields with the speicified name to a new type created by
// NewType. The field's type has to be set using SetCurrentObjectType
//
func AddStructField(fieldname string) {
    var temp uint64;
    if Compile != 0 {
		temp = libgogo.HasField(fieldname, CurrentType);
		if temp != 0 {
		    SymbolTableError("duplicate", "", "field", fieldname);
		} else {
		    CurrentObject = libgogo.NewObject(fieldname, "", libgogo.CLASS_FIELD); //A field has no package name
		    libgogo.AddFields(CurrentObject, CurrentType);
		}
    }
}

//
// Flags the type the current object or field to be a pointer and
// not the type itself
//
func SetCurrentObjectTypeToPointer() {
    if Compile != 0 {
        CurrentObject.PtrType = 1; //Type is pointer
    }
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
    var tempstr string = "";
    var tempstr2 string;
    var boolFlag uint64;

    if Compile != 0 {
	    if InsideStructDecl == 1 { //Allow types in struct declarations which are already forward declared
	        basetype = libgogo.GetType(typename, packagename, GlobalTypes, 1);
	    } else {
	        basetype = libgogo.GetType(typename, packagename, GlobalTypes, 0);
	    }
	    if basetype == nil {
	        if InsideStructDecl == 1 {
	            boolFlag = libgogo.StringCompare(typename, CurrentType.Name);
	            if boolFlag == 0 {
	                if CurrentObject.PtrType == 1 { //Allow pointer to own type
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
	        CurrentObject.ObjType = basetype;
	    } else { //Array
	        if basetype != nil {
	            libgogo.StringAppend(&tempstr, basetype.Name);
	        }
	        libgogo.StringAppend(&tempstr, "Array");
	        tempstr2 = libgogo.IntToString(arraydim);
	        libgogo.StringAppend(&tempstr, tempstr2);
	        temptype = libgogo.NewType(tempstr, packagename, 0, arraydim, basetype);
	        CurrentObject.ObjType = temptype; //Don't add array type to global list to avoid duplicate type name errors
	    }
    }
}

//
// Adds a new variable to the global symbol table with the specified
// name. The variable type has to be set using SetCurrentObjectType
//
func NewVariable(name string) {
    var TempObject *libgogo.ObjectDesc;
    var temp uint64;
    if Compile != 0 {
		CurrentObject = libgogo.NewObject(name, CurrentPackage, libgogo.CLASS_VAR);
		if InsideFunction == 0 { //Global objects or function parameters
		    if InsideFunctionVarDecl == 0 { //Global objects
		        TempObject = libgogo.GetObject(name, CurrentPackage, GlobalObjects);
		        if TempObject != nil {
		            SymbolTableError("duplicate", "global", "identifier", name);
		        }
		        GlobalObjects = libgogo.AppendObject(CurrentObject, GlobalObjects);
		    } else { //Function parameters
		        temp = libgogo.HasField(name, CurrentFunction);
		        if temp != 0 {
		            SymbolTableError("duplicate", "parameter", "name", name);
		        } else {
		            CurrentObject.Class = libgogo.CLASS_PARAMETER;
		            CurrentObject.PackageName = ""; //A parameter has no package name
		            libgogo.AddParameters(CurrentObject, CurrentFunction);
		        }
		    }
		} else { //Function-local objects
	        TempObject = libgogo.GetObject(name, CurrentPackage, LocalObjects);
	        if TempObject != nil {
	            SymbolTableError("duplicate", "local", "identifier", name);
	        }
	        TempObject = libgogo.GetObject(name, CurrentPackage, CurrentFunction.Fields);
	        if TempObject != nil {
	            SymbolTableError("There is already a parameter", "", "named", name);
	        }
	        LocalObjects = libgogo.AppendObject(CurrentObject, LocalObjects);
		}
    }
}

//
// When called when the reaching the end of a function during parsing, the
// local symbol table is purged and optional debug information about the
// symbol table before purging are printed
//
func EndOfFunction() {
    if Compile != 0 {
        PrintLocalSymbolTable(); //Print local symbol table
        LocalObjects = nil; //Delete local objects
    }
}

//
// Creates a new variable item from a variable object, considering the according address offsets
// kind = 0: Local variable, 1: Global variable, 2: Local parameter
//
func VariableObjectDescToItem(obj *libgogo.ObjectDesc, item *libgogo.Item, kind uint64) {
    var tempAddr uint64;
    var size uint64;
    if kind == 0 { //Local variable
        tempAddr = libgogo.GetObjectOffset(obj, LocalObjects);
        size = libgogo.GetObjectSizeAligned(obj);
        tempAddr = tempAddr + size - 8; //Due to sign of the offset (p.e. -8(SP)), the offsets starts at the last byte and end at the first one
    } else { //Global variable or local parameter
        if kind == 1 { //Global variable
            tempAddr = libgogo.GetObjectOffset(obj, GlobalObjects);
        } else { //Local parameter (kind = 2)
            tempAddr = libgogo.GetObjectOffset(obj, CurrentFunction.Fields);
        }
    }
    libgogo.SetItem(item, libgogo.MODE_VAR, obj.ObjType, obj.PtrType, tempAddr, 0, kind); //Varible item of given kind
}

//
// Converts an object representing a parameter of a function call into an item with the correct address offset
//
func ObjectToStackParameter(obj *libgogo.ObjectDesc, FunctionCalled *libgogo.TypeDesc, stackoffset uint64) *libgogo.Item {
    var OldLocalObjects *libgogo.ObjectDesc;
    var ReturnItem *libgogo.Item;
    var ObjSize uint64;
    ObjSize = libgogo.GetObjectSizeAligned(obj);
    OldLocalObjects = LocalObjects; //Save pointer to local objects
    LocalObjects = FunctionCalled.Fields; //Use parameters with local object offsets
    ReturnItem = libgogo.NewItem();
    VariableObjectDescToItem(obj, ReturnItem, 0); //Treat parameter as if it was a local object
    ReturnItem.A = stackoffset - 8 - ReturnItem.A - 8 + ObjSize; //Add offset (total size of parameters and variables); compensate both local offsets (stackoffset and ReturnItem.A) by subtracting -8 for each, then adding the parameter size
    if FunctionCalled.ForwardDecl == 1 {
        ReturnItem.LinkerInformation = "##"; //Invalid characters to make assembly impossible
        libgogo.StringAppend(&ReturnItem.LinkerInformation, FunctionCalled.PackageName);
        libgogo.StringAppend(&ReturnItem.LinkerInformation, "·");
        libgogo.StringAppend(&ReturnItem.LinkerInformation, FunctionCalled.Name);
        libgogo.StringAppend(&ReturnItem.LinkerInformation, "##");
        ReturnItem.A = ReturnItem.A + 100000; //Add dummy offset of 100000 to bias negative calculations in order to prevent underflow
    }
    LocalObjects = OldLocalObjects; //Restore old local objects pointer
    return ReturnItem;
}

//
// Adds a new function to the global symbol table with the specified
// name
//
func NewFunction(name string, packagename string, forwarddecl uint64) *libgogo.TypeDesc {
    var TempType *libgogo.TypeDesc;
    var DontAppend uint64 = 0;
    if Compile != 0 {
		CurrentFunction = libgogo.NewType(name, packagename, forwarddecl, 0, nil);
        TempType = libgogo.GetType(name, packagename, GlobalFunctions, 1);
        if TempType != nil {
            if TempType.ForwardDecl == 1 { //Unset forward declaration
                TempType.ForwardDecl = forwarddecl;
                DontAppend = 1;
            } else {
                SymbolTableError("duplicate function", name, "in package", name);
            }
        }
        if DontAppend == 0 {
            GlobalFunctions = libgogo.AppendType(CurrentFunction, GlobalFunctions);
        }
    }
    return CurrentFunction;
}
