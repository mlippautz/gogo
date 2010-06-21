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
// Parameter counter when parsing implementations of previously forward declared functions
//
var fwdParamCount uint64 = 0;

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
// Flag indicating whether file needs linking
//
var NeedsLink uint64 = 0;

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
    var TempObject *libgogo.ObjectDesc;

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
                    NeedsLink = 1; // We need linking because of a fwd. declared type
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
	    
  		if (InsideFunction == 0) && (InsideFunctionVarDecl == 1) { //Function parameters
  		    if (CurrentFunction.ForwardDecl == 0) && (CurrentFunction.Base != nil) { //Check implementation of forward declaration
  		        TempObject = libgogo.GetParameterAt(fwdParamCount, CurrentFunction);
  		        CorrectArtificialParameterForced(CurrentFunction, TempObject, CurrentObject.ObjType, CurrentObject.PtrType, 0); //Perform up-conversion if necessary
  	            if (TempObject.ObjType != CurrentObject.ObjType) || (TempObject.PtrType != CurrentObject.PtrType) {
  	                tempstr = libgogo.IntToString(fwdParamCount);
                    SymbolTableError("Parameter", tempstr, "has been forward declared with different type, function", CurrentFunction.Name);
                }
  		    }
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
		        if (CurrentFunction.ForwardDecl == 0) && (CurrentFunction.Base != nil) { //Check implementation of forward declaration
		            fwdParamCount = fwdParamCount + 1; //One more parameter
		            ParameterCheck_Less(CurrentFunction, fwdParamCount); //Check number of parameters
    		        TempObject = libgogo.GetParameterAt(fwdParamCount, CurrentFunction);
    		        temp = libgogo.HasField(name, CurrentFunction); //Check for duplicate parameter names
	                if temp != 0 {
	                    SymbolTableError("duplicate", "parameter", "name", name);
	                }
    		        TempObject.PackageName = CurrentFunction.PackageName;
    		        TempObject.Name = name; //Set name of current parameter
		            //A type check can only be performed as soon as the type is parsed (see parser.go)
		        } else {
	                temp = libgogo.HasField(name, CurrentFunction);
	                if temp != 0 {
	                    SymbolTableError("duplicate", "parameter", "name", name);
	                } else {
	                    CurrentObject.Class = libgogo.CLASS_PARAMETER;
	                    CurrentObject.PackageName = CurrentFunction.PackageName;
	                    libgogo.AddParameters(CurrentObject, CurrentFunction);
	                }
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
    if FunctionCalled.ForwardDecl == 1 { //Add linker information
        ReturnItem.LinkerInformation = "##1##";
        libgogo.StringAppend(&ReturnItem.LinkerInformation, FunctionCalled.PackageName);
        libgogo.StringAppend(&ReturnItem.LinkerInformation, "·");
        libgogo.StringAppend(&ReturnItem.LinkerInformation, FunctionCalled.Name);
        libgogo.StringAppend(&ReturnItem.LinkerInformation, "##");
        ReturnItem.A = ReturnItem.A + 100000; //Bias to avoid underflows
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
    var NewFct *libgogo.TypeDesc;
    if Compile != 0 {
		NewFct = libgogo.NewType(name, packagename, forwarddecl, 0, nil);
		if forwarddecl == 0 {
		    CurrentFunction = NewFct;
		} else {
            NeedsLink = 1; // We need linking because of a fwd declared function
        }
        TempType = libgogo.GetType(name, packagename, GlobalFunctions, 1);
        if TempType != nil {
            if TempType.ForwardDecl == 1 { //Unset forward declaration
                TempType.ForwardDecl = forwarddecl;
                DontAppend = 1; //Don't add function again
                NewFct = TempType;
                CurrentFunction = NewFct;
                fwdParamCount = 0; //Start over with new parameters
            } else {
                SymbolTableError("duplicate function", name, "in package", packagename);
            }
        }
        if DontAppend == 0 {
            GlobalFunctions = libgogo.AppendType(NewFct, GlobalFunctions);
        }
    }
    return NewFct;
}

//
// Adds a return parameter to the given function if necessary and returns it
//
func AddReturnParameter(CurrentFunction *libgogo.TypeDesc, ReturnValuePseudoObject *libgogo.ObjectDesc) *libgogo.ObjectDesc {
    var TempObject *libgogo.ObjectDesc;
    if (CurrentFunction.ForwardDecl == 0) && (CurrentFunction.Base != nil) { //Implementation of forward declared function
        TempObject = libgogo.GetObject("return value", "", CurrentFunction.Fields); //Check if there is a return value
        if TempObject != nil { //Type check
            //TODO (maybe): Actually required down-conversion
            //CorrectArtificialParameterForced(CurrentFunction, TempObject, ReturnValuePseudoObject.ObjType, ReturnValuePseudoObject.PtrType, 0); //Perform up-conversion if necessary
            if (TempObject.ObjType != ReturnValuePseudoObject.ObjType) || (TempObject.PtrType != ReturnValuePseudoObject.PtrType) {
                SymbolTableError("Function has been forward declared with", "different", "return value type, function", CurrentFunction.Name);
            }
        }
        ReturnValuePseudoObject = TempObject; //Return return value object of forward declaration
    } else { //Default return value declaration
        if ReturnValuePseudoObject != nil {
            libgogo.AddParameters(ReturnValuePseudoObject, CurrentFunction); //Treat return value like an additional parameter at the end of the parameter list
            CurrentFunction.Len = CurrentFunction.Len - 1; //Don't count parameter as input parameter
        }
    }
    return ReturnValuePseudoObject;
}

//
// Checks if there are sufficient parameters for a previously forward declared function
//
func FwdDeclCheckIfNecessary() {
    if (CurrentFunction.ForwardDecl == 0) && (CurrentFunction.Base != nil) {
        ParameterCheck_More(CurrentFunction, fwdParamCount); //Check parameter count
    }
}

//
// Fails if the current function declaration has no return value, but is required to have one through previous forward declarations
//
func AssertReturnTypeIfNecessary(ReturnValuePseudoObject *libgogo.ObjectDesc) {
    var TempObject *libgogo.ObjectDesc;
    if (CurrentFunction.ForwardDecl == 0) && (CurrentFunction.Base != nil) { //Implementation of forward declared function
        TempObject = libgogo.GetObject("return value", "", CurrentFunction.Fields); //Check if there is a return value
        if (TempObject != nil) && (ReturnValuePseudoObject == nil) {
            SymbolTableError("Function has been forward declared with", "a", "return value, function", CurrentFunction.Name);
        }
    }
}

//
// Applies a name selector to a function defined by its package name
// Returns a new forward declared function if there is no function name matching in the package given
//
func ApplyFunctionSelector(PackageFunction *libgogo.TypeDesc, name string) *libgogo.TypeDesc {
    var boolFlag uint64;
    var tempFcn *libgogo.TypeDesc;
    var ReturnedFunction *libgogo.TypeDesc = nil;
    boolFlag = libgogo.StringLength(PackageFunction.Name);
    if boolFlag != 0 {
        SymbolTableError("Cannot apply a selector to", "a", "function, function", PackageFunction.Name);
    } else {
        tempFcn = libgogo.GetType(name, PackageFunction.PackageName, GlobalFunctions, 1); //Check global functions
	    if tempFcn == nil { //New forward declaration
	        tempFcn = NewFunction(name, PackageFunction.PackageName, 1);
	        ReturnedFunction = tempFcn; //Assign new tempFcn pointer to PackageFunction (outside this function)
	    }
        PackageFunction.Name = tempFcn.Name;
        PackageFunction.Len = tempFcn.Len;
        PackageFunction.Fields = tempFcn.Fields;
        PackageFunction.ForwardDecl = tempFcn.ForwardDecl;
        PackageFunction.Base = tempFcn.Base;
    }
    return ReturnedFunction;
}

//
// Sets a function's name based on its package name and resets its package name to the current package
// Returns a new forward declared function if there is no function name matching in current package
//
func PackageFunctionToNameFunction(PackageFunction *libgogo.TypeDesc) *libgogo.TypeDesc {
    var tempFcn *libgogo.TypeDesc;
    var ReturnedFunction *libgogo.TypeDesc = nil;
    tempFcn = libgogo.GetType(PackageFunction.PackageName, CurrentPackage, GlobalFunctions, 1); //Check global functions
    if tempFcn == nil { //New forward declaration
        tempFcn = NewFunction(PackageFunction.PackageName, CurrentPackage, 1);
        ReturnedFunction = tempFcn; //Assign new tempFcn pointer to PackageFunction (outside this function)
    }
    PackageFunction.Name = tempFcn.Name;
    PackageFunction.PackageName = tempFcn.PackageName;
    PackageFunction.Len = tempFcn.Len;
    PackageFunction.Fields = tempFcn.Fields;
    PackageFunction.ForwardDecl = tempFcn.ForwardDecl;
    PackageFunction.Base = tempFcn.Base;
    return ReturnedFunction;
}

func VariableOrFieldAccess(item *libgogo.Item, packagename string, name string) {
    var boolFlag uint64;
    var tempObject *libgogo.ObjectDesc;
    var tempList *libgogo.ObjectDesc;
    if (item.Itemtype == nil) && (item.A == 0) && (item.R == 0) { //Item undefined, only package known => Find object
        tempObject = libgogo.GetObject(name, packagename, LocalObjects); //Check local objects
	    tempList = LocalObjects;
	    if tempObject == nil {
	        tempObject = libgogo.GetObject(name, packagename, GlobalObjects); //Check global objects
	        tempList = GlobalObjects;
	    }
	    if tempObject == nil {
	        SymbolTableError("Package", packagename, "has no variable named", name);
	    }
	    if tempList == LocalObjects { //Local
	        VariableObjectDescToItem(tempObject, item, 0); //Local variable
	    } else { //Global
	        VariableObjectDescToItem(tempObject, item, 1); //Global variable
  	    }
    } else { //Field access
        if Compile != 0 {
            if item.Itemtype == nil {
                SymbolTableError("Type has", "no", "fields:", "?");
            }
            if item.Itemtype.Form != libgogo.FORM_STRUCT { //Struct check
                SymbolTableError("Type is", "not a", "struct:", item.Itemtype.Name);
            } else {
                boolFlag = libgogo.HasField(name, item.Itemtype);
                if boolFlag == 0 { //Field check
                    SymbolTableError("Type", item.Itemtype.Name, "has no field named", name);
                } else {
                    tempObject = libgogo.GetField(name, item.Itemtype);
                    boolFlag = libgogo.GetFieldOffset(tempObject, item.Itemtype); //Calculate offset
                    GenerateFieldAccess(item, boolFlag);
                    item.Itemtype = tempObject.ObjType; //Set item type to field type
                    item.PtrType = tempObject.PtrType;
                }
            }
        }
    }
}

//
// Checks whether array access is possible to the given item
//
func ArrayAccessCheck(item *libgogo.Item, packagename string)  {
    var boolFlag uint64;
    if Compile != 0 {
        if (item.Itemtype == nil) && (item.A == 0) && (item.R == 0) {
            SymbolTableError("No array access possible to", "", "package", packagename);
        }
        if item.Itemtype.Form != libgogo.FORM_ARRAY { //Array check
            SymbolTableError("Type is", "not an", "array:", item.Itemtype.Name);
        }
        if item.Itemtype == string_t { //Derefer string address at offset 0 to access actual byte array of characters
            boolFlag = item.PtrType; //Save old value of PtrType
            item.PtrType = 1; //Force deref.
            DereferItemIfNecessary(item); //Actual deref.
            item.PtrType = boolFlag; //Restore old value of PtrType
        }
    }
}

//
// Tries to find a variable based on the information collected in the item given
//
func FindIdentifierAndParseSelector(item *libgogo.Item) {
    var boolFlag uint64;
    var tempObject *libgogo.ObjectDesc;
    var tempList *libgogo.ObjectDesc;
    var packagename string;
    if Compile != 0 {
		//Token can be package name
		boolFlag = libgogo.FindPackageName(tok.strValue, GlobalObjects); //Check global objects
		if boolFlag == 0 {
		    boolFlag = libgogo.FindPackageName(tok.strValue, LocalObjects); //Check local objects
		}
		if (boolFlag == 0) && (CurrentFunction != nil) {
		    boolFlag = libgogo.FindPackageName(tok.strValue, CurrentFunction.Fields); //Check local parameters
		}
		if boolFlag == 0 { //Token is not package name, but identifier
			tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, LocalObjects); //Check local objects
			tempList = LocalObjects;
			if tempObject == nil {
                if CurrentFunction != nil {
           			tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, CurrentFunction.Fields); //Check local parameters
        			tempList = CurrentFunction.Fields;
    			}
    			if tempObject == nil {
    				tempObject = libgogo.GetObject(tok.strValue, CurrentPackage, GlobalObjects); //Check global objects
	    			tempList = GlobalObjects;
	    		}
			}
			if tempObject == nil {
				SymbolTableError("Undefined", "", "variable", tok.strValue);
			}
			if tempList == LocalObjects { //Local
    			VariableObjectDescToItem(tempObject, item, 0); //Local variable
			} else { //Global or parameter
			    if tempList == GlobalObjects { //Global
    				VariableObjectDescToItem(tempObject, item, 1); //Global variable
    			} else { //Parameter
    				VariableObjectDescToItem(tempObject, item, 2); //Local parameter
    	        }
			}
		    ParseSelector(item, CurrentPackage); //Parse selectors for an object in the current package
		} else { //Token is package name
		    libgogo.SetItem(item, 0, nil, 0, 0, 0, 0); //Mark item as not being set
		    packagename = tok.strValue; //Save package name
		    ParseSelector(item, tok.strValue); //Parse selectors for an undefined object in the given package
		    if (item.Itemtype == nil) && (item.A == 0) && (item.R == 0) {
		        SymbolTableError("Cannot use package", "", "as a variable:", packagename);
		    }
		}
    } else {
        ParseSelector(item, CurrentPackage);
    }
}

//
// Tries to find a function based on the information collected in the item given; if there is none, an according forward declaration will be made
//
func FindIdentifierAndParseSelector_FunctionCall(FunctionCalled *libgogo.TypeDesc) *libgogo.TypeDesc {
    var boolFlag uint64;
    var tempFcn *libgogo.TypeDesc;
    if Compile != 0 {
		//Token can be package name
		boolFlag = libgogo.FindTypePackageName(tok.strValue, GlobalFunctions); //Check global functions
		if boolFlag == 0 { //Token is not package name, but identifier
			tempFcn = libgogo.GetType(tok.strValue, CurrentPackage, GlobalFunctions, 1); //Check global functions
            if tempFcn == nil { //New forward declaration
                FunctionCalled.PackageName = tok.strValue; //Set package name
                tempFcn = ParseSelector_FunctionCall(FunctionCalled); //Parse selector for an undefined function in the given package
                boolFlag = libgogo.StringLength(tempFcn.Name);
                if boolFlag == 0 {
		            SymbolTableError("Cannot use package", "", "as a function:", tempFcn.PackageName);
		        }
            } //else: tempFcn is already the return value
        } else { //Token is package name
            FunctionCalled.PackageName = tok.strValue; //Set package name
            tempFcn = ParseSelector_FunctionCall(FunctionCalled); //Parse selector for an undefined function in the given package
            boolFlag = libgogo.StringLength(tempFcn.Name);
            if boolFlag == 0 {
		        SymbolTableError("Cannot use package", "", "as a function:", tempFcn.PackageName);
		    }
        }
    } else {
        tempFcn = ParseSelector_FunctionCall(FunctionCalled);
    }
    return tempFcn;
}

//
// Checks whether a given type list contains any forward declarations.
// Returns: 1 if there are any, 0 otherwise
//
func ContainsFwdDecls(list *libgogo.TypeDesc) uint64 {
    var tmpObjType *libgogo.TypeDesc;
    var retValue uint64 = 0;
    if list != nil {
        for tmpObjType = list; tmpObjType.Next != nil; tmpObjType = tmpObjType.Next { 
            if tmpObjType.ForwardDecl != 0 {
                retValue = 1;
            }
        }
    }
    return retValue;
}
