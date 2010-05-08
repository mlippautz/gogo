// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

type ObjectDesc struct {
    name string;
    class uint64;
    objtype *TypeDesc;
    next *ObjectDesc;
};

type TypeDesc struct {
    name string;
    form uint64;
    len uint64;
    fields *ObjectDesc;
    base *TypeDesc;
    next *TypeDesc;
};

//
// Pseudo constants that specify the descriptor sizes 
//
var OBJECT_SIZE uint64 = 32+8; //4*8 bytes space for an object, extra 8 for the string length
var TYPE_SIZE uint64 = 48+8;  //6*8 bytes space for a type, extra 8 for the string length

//
// Classes for objects
//
var CLASS_VAR uint64 = 1;
var CLASS_TYPE uint64 = 2;
var CLASS_FIELD uint64 = 3;

//
// Available types
//
var TYPE_UINT64 uint64 = 1;
var TYPE_STRING uint64 = 2;
var TYPE_STRUCT uint64 = 3;
var TYPE_ARRAY uint64 = 4;

//
// List of global objects and declared types
//
var GlobalObjects *ObjectDesc = nil;
var Types *TypeDesc = nil;

//
// Convert the uint64 value (returned from malloc) to a real object address
//
func Uint64ToObjectDescPtr(adr uint64) *ObjectDesc;

//
// Convert the uint64 value (returned from malloc) to a real type object address
//
func Uint64ToTypeDescPtr(adr uint64) *TypeDesc;


//
// Appends an object to the list
//
func AppendObject(object *ObjectDesc, list *ObjectDesc) *ObjectDesc {
    var tmpObj *ObjectDesc = object;
    if list != nil {
        for tmpObj = list; tmpObj.next != nil; tmpObj = tmpObj.next {
        }
        tmpObj.next = object;
        tmpObj = list;
    }
    return tmpObj;    

		//Alternate version: PrependObject:
		/*object.next = list;
    return object;*/
}

//
// Appends a type to the list
// 
func AppendType(objtype *TypeDesc, list *TypeDesc) *TypeDesc {
    var tmpObjType *TypeDesc = objtype;
    if list != nil {
        for tmpObjType = list; tmpObjType.next != nil; tmpObjType = tmpObjType.next { }
        tmpObjType.next = objtype;
        tmpObjType = list;
    }
    return tmpObjType;

    //Alternate version: PrependType:
    /*objtype.next = list;
    return objtype;*/
}

//
// Add a field in form of an object descriptor to the type descriptor given
//
func AddFields(object *ObjectDesc, objtype *TypeDesc) {
    objtype.fields = AppendObject(object, objtype.fields);
}

//
// Sets an object's type
//
func SetObjType(object *ObjectDesc, objtype *TypeDesc) {
    object.objtype = objtype;
}

//
// Fetches an object with a specific identifier or nil if it is not in the specified list
//
func GetObject(name string, list *ObjectDesc) *ObjectDesc {
    var tmpObject *ObjectDesc;
    var retValue *ObjectDesc = nil;
    for tmpObject = list; tmpObject != nil; tmpObject = tmpObject.next {
        if StringCompare(tmpObject.name,name) == 0 {
            retValue = tmpObject;
            break;
        }
    }
    return retValue;
}

//
// Fetches a type with a given name or nil if it is not in the specified list
//
func GetType(name string, list *TypeDesc) *TypeDesc {
    var tmpType *TypeDesc;
    var retValue *TypeDesc = nil;
    for tmpType = list; tmpType != nil; tmpType = tmpType.next {
        if StringCompare(tmpType.name,name) == 0 {
            retValue = tmpType;
            break;
        }
    }
    return retValue;
}

//
// Creates a new object
//
func NewObject(name string, class uint64) *ObjectDesc {
    var adr uint64 = Alloc(OBJECT_SIZE);
    var obj *ObjectDesc = Uint64ToObjectDescPtr(adr);
    obj.name = name; //TODO: Copy string?
    obj.objtype = nil;
    obj.class = class;
    obj.next = nil;
    return obj;
}

//
// Creates a new type
//
func NewType(name string, len uint64) *TypeDesc {
    var adr uint64 = Alloc(TYPE_SIZE);
    var objtype *TypeDesc = Uint64ToTypeDescPtr(adr);
    objtype.name = name; //TODO: Copy string?
    objtype.len = len;
    objtype.next = nil;
    objtype.fields = nil;
    return objtype;
}

func PrintObjects(list *ObjectDesc) {
    var o *ObjectDesc;
    for o = list; o != nil; o = o.next {
        PrintString("Object ");
        PrintString(o.name);
        PrintString(" (type ");
        if o.objtype != nil {
            PrintString(o.objtype.name);
        } else {
            PrintString("<unknown>");
        }
        PrintString(")\n");
    }
}

func PrintTypes(list *TypeDesc) {
    var t *TypeDesc;
    var o *ObjectDesc;
    for t = list; t != nil; t = t.next {
        PrintString("Type ");
        PrintString(t.name);
        if t.len != 0 {
            PrintString(" (length: ");
            PrintNumber(t.len);
            PrintString(")");
        }
        PrintString("\n");
				for o = t.fields; o != nil; o = o.next {
            PrintString("  ");
            PrintString(o.name);
            PrintString(" (type ");
            if o.objtype != nil {
                PrintString(o.objtype.name);
            } else {
                PrintString("<unknown>");
            }
            PrintString(")\n");
        }
    }
}
