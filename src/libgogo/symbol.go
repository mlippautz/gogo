// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

type ObjectDesc struct {
    name string;
    class uint64;
    objtype *TypeDesc;
    ptrtype byte; //If 0, type objtype, otherwise type *objtype
    next *ObjectDesc;
};

type TypeDesc struct {
    name string;
    packagename string;
    form uint64;
    len uint64;
    fields *ObjectDesc;
    base *TypeDesc;
    next *TypeDesc;
};

//
// Pseudo constants that specify the descriptor sizes 
//
var OBJECT_SIZE uint64 = 48; //6*8 bytes space for an object, extra 8 for the string length
var TYPE_SIZE uint64 = 72;  //9*8 bytes space for a type, extra 8 for the string length

//
// Classes for objects
//
var CLASS_VAR uint64 = 1;
var CLASS_FIELD uint64 = 2;

//
// Forms for types
//
var FORM_SIMPLE uint64 = 1;
var FORM_STRUCT uint64 = 2;
var FORM_ARRAY uint64 = 3;

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
    objtype.form = FORM_STRUCT;
    objtype.fields = AppendObject(object, objtype.fields);
    objtype.len = objtype.len + GetTypeSize(objtype);
}

//
// Sets an object's type
//
func SetObjType(object *ObjectDesc, objtype *TypeDesc) {
    object.objtype = objtype;
}

//
// Marks the type of the given object as pointer
//
func FlagObjectTypeAsPointer(object *ObjectDesc) {
    object.ptrtype = 1;
}

//
// Returns whether an object's type is a pointer
//
func IsPointerType(object *ObjectDesc) byte {
    return object.ptrtype;
}

//
// Returns a type's name
//
func GetTypeName(objtype *TypeDesc) string {
    return objtype.name;
}

//
// Returns an object's name
//
func GetObjectName(obj *ObjectDesc) string {
    return obj.name;
}

func GetTypeSize(objtype *TypeDesc) uint64 {
    var size uint64 = 0;
    var tempobj *ObjectDesc;
    if objtype != nil {
        if objtype.form == FORM_SIMPLE {
            size = objtype.len;
        }
        if objtype.form == FORM_STRUCT {
            for tempobj = objtype.fields; tempobj != nil; tempobj = tempobj.next { //Sum of all fields
                size = size + GetObjectSize(tempobj); //Add size of each field
            }
        }
        if objtype.form == FORM_ARRAY {
            size = objtype.len * GetTypeSize(objtype.base); //Array length * size of one item
        }
    } //TODO: if objtype is nil => error!
    return size;
}

func GetObjectSize(obj *ObjectDesc) uint64 {
    var size uint64;
    if obj.ptrtype == 1 { //Pointer only, not the whole type
       size = 8;
    } else { //Actual type
       size = GetTypeSize(obj.objtype);
    }
    return size;
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
func GetType(name string, packagename string, list *TypeDesc) *TypeDesc {
    var tmpType *TypeDesc;
    var retValue *TypeDesc = nil;
    for tmpType = list; tmpType != nil; tmpType = tmpType.next {
        if (StringCompare(tmpType.name,name) == 0) && ((StringLength(tmpType.packagename) == 0) || (StringCompare(tmpType.packagename,packagename) == 0)) { //Empty package name indicates internal types
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
    obj.class = class;
    obj.objtype = nil;
    obj.ptrtype = 0;
    obj.next = nil;
    return obj;
}

//
// Creates a new type
//
func NewType(name string, packagename string, len uint64, basetype *TypeDesc) *TypeDesc {
    var adr uint64 = Alloc(TYPE_SIZE);
    var objtype *TypeDesc = Uint64ToTypeDescPtr(adr);
    objtype.name = name; //TODO: Copy string?
    objtype.packagename = packagename; //TODO: Copy string?
    if basetype != nil {
        objtype.form = FORM_ARRAY;
    } else {
        objtype.form = FORM_SIMPLE;
    }
    objtype.len = len;
    objtype.next = nil;
    objtype.fields = nil;
    objtype.base = basetype;
    return objtype;
}

func PrintObjects(list *ObjectDesc) {
    var o *ObjectDesc;
    for o = list; o != nil; o = o.next {
        PrintString("Object ");
        PrintString(o.name);
        PrintString(" (type: ");
        if o.ptrtype != 0 {
            PrintString("pointer to ");
        }
        if o.objtype != nil {
            if o.objtype.base != nil {
                PrintString("array of ");
                PrintString(o.objtype.base.name);
                PrintString(" of length ");
                PrintNumber(o.objtype.len);
                PrintString(", internally named ");
            }
            if StringLength(o.objtype.packagename) != 0 {
                PrintString(o.objtype.packagename);
                PrintChar('.');
            }
            PrintString(o.objtype.name);
        } else {
            PrintString("<unknown>");
        }
        PrintString(", size: ");
        PrintNumber(GetObjectSize(o));
        PrintString(")\n");
    }
}

func PrintTypes(list *TypeDesc) {
    var t *TypeDesc;
    var o *ObjectDesc;
    for t = list; t != nil; t = t.next {
        PrintString("Type ");
        PrintString(t.name);
        PrintString(" (size: ");
        PrintNumber(GetTypeSize(t));
        PrintString(")\n");
				for o = t.fields; o != nil; o = o.next {
            PrintString("  ");
            PrintString(o.name);
            PrintString(" (type: ");
            if o.ptrtype != 0 {
                PrintString("pointer to ");
            }
            if o.objtype != nil {
                if o.objtype.base != nil {
                    PrintString("array of ");
                    PrintString(o.objtype.base.name);
                    PrintString(" of length ");
                    PrintNumber(o.objtype.len);
                    PrintString(", internally named ");
                }
                if StringLength(o.objtype.packagename) != 0 {
                    PrintString(o.objtype.packagename);
                    PrintChar('.');
                }
                PrintString(o.objtype.name);
                PrintString(", size: ");
                PrintNumber(GetTypeSize(o.objtype));
            } else {
                PrintString("<unknown>");
            }
            PrintString(")\n");
        }
    }
}
