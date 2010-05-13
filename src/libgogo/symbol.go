// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

type ObjectDesc struct {
    name string;
    packagename string;
    class uint64;
    objtype *TypeDesc;
    ptrtype byte; //If 0, type objtype, otherwise type *objtype
    next *ObjectDesc;
};

type TypeDesc struct {
    name string;
    packagename string;
    forwarddecl byte; //If 1, type is forward declared
    form uint64;
    len uint64;
    fields *ObjectDesc;
    base *TypeDesc;
    next *TypeDesc;
};

//
// Pseudo constants that specify the descriptor sizes 
//
var OBJECT_SIZE uint64 = 64; //8*8 bytes space for an object
var TYPE_SIZE uint64 = 80;  //10*8 bytes space for a type

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
}

//
// Returns 0 if the given (struct) type has a field with the given name, or 1 otherwise
//
func HasField(name string, objtype *TypeDesc) uint64 {
    var tmpObj *ObjectDesc;
    var retVal uint64 = 0;
    tmpObj = GetObject(name, "", objtype.fields);
    if tmpObj != nil {
        retVal = 1;
    }
    return retVal;
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

//
// Returns 0 if the given type if forward declared, or 1 if it is not
//
func IsForwardDecl(objtype *TypeDesc) byte {
    return objtype.forwarddecl;
}

//
// Unsets the forward declaration flag of the given type, making it a "normal" type
//
func UnsetForwardDecl(objtype *TypeDesc) {
    objtype.forwarddecl = 0;
}

//
// Returns the (memory) size of the given type in bytes which is required when hypothetically allocating one variable of this type
// Same semantics as sizeof(objtype)
//
func GetTypeSize(objtype *TypeDesc) uint64 {
    var size uint64 = 0;
    var tempobj *ObjectDesc;
    var tmpSize uint64;
    if objtype != nil {
        if objtype.forwarddecl == 0 {
            if objtype.form == FORM_SIMPLE {
                size = objtype.len;
            }
            if objtype.form == FORM_STRUCT {
                for tempobj = objtype.fields; tempobj != nil; tempobj = tempobj.next { //Sum of all fields
                    tmpSize = GetObjectSize(tempobj);
                    size = size + tmpSize; //Add size of each field
                }
            }
            if objtype.form == FORM_ARRAY {
                tmpSize = GetTypeSize(objtype.base);
                size = objtype.len * tmpSize; //Array length * size of one item
            }
        } else {
            ; //TODO: if type is only forward declared => error!
        }
    } else {
        ; //TODO: if objtype is nil => error!
    }
    return size;
}

//
// Returns the (memory) size required by the object in bytes, considering whether or not the object is a pointer
// Same semantics as sizeof(obj)
//
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
// Calculates the (memory) offset of a given field of the specified (struct) type in bytes
//
func GetObjectOffset(obj *ObjectDesc, list *ObjectDesc) uint64 {
    var offset uint64 = 0;
    var tmp *ObjectDesc;
    var tmpSize uint64;
    for tmp = list; tmp != nil; tmp = tmp.next {
        if tmp == obj {
            break;
        }
        tmpSize = GetObjectSize(tmp);
        offset = offset + tmpSize;
    }
    if tmp == nil {
        offset = 0; //TODO: Raise error as obj is appearently not within the list
    }
    return offset;
}

//
// Fetches an object with a specific identifier or nil if it is not in the specified list
//
func GetObject(name string, packagename string, list *ObjectDesc) *ObjectDesc {
    var tmpObject *ObjectDesc;
    var retValue *ObjectDesc = nil;
    var nameCompare uint64;
    var strLen uint64;
    var packageNameCompare uint64;
    for tmpObject = list; tmpObject != nil; tmpObject = tmpObject.next {
        nameCompare = StringCompare(tmpObject.name,name);
        strLen = StringLength(tmpObject.packagename);
        packageNameCompare = StringCompare(tmpObject.packagename,packagename);
        if (nameCompare == 0) && ((strLen == 0) || (packageNameCompare == 0)) { //Empty package name indicates internal types
            retValue = tmpObject;
            break;
        }
    }
    return retValue;
}

//
// Fetches a type with a given name or nil if it is not in the specified list
// If includeforward is 0, no forward declared types will be returned; otherwise, forward declared types will also be returned
//
func GetType(name string, packagename string, list *TypeDesc, includeforward byte) *TypeDesc {
    var tmpType *TypeDesc;
    var retValue *TypeDesc = nil;
    var nameCompare uint64;
    var strLen uint64;
    var packageNameCompare uint64;
    for tmpType = list; tmpType != nil; tmpType = tmpType.next {
        nameCompare = StringCompare(tmpType.name,name);
        strLen = StringLength(tmpType.packagename);
        packageNameCompare = StringCompare(tmpType.packagename,packagename);
        if (nameCompare == 0) && ((strLen == 0) || (packageNameCompare == 0)) { //Empty package name indicates internal types
            if (includeforward == 1) || ((includeforward == 0) && (tmpType.forwarddecl == 0)) {
                retValue = tmpType;
                break;
            }
        }
    }
    return retValue;
}

//
// Returns the first type from the list which is forward declared
//
func GetFirstForwardDeclType(list *TypeDesc) *TypeDesc {
    var retValue *TypeDesc;
    for retValue = list; retValue != nil; retValue = retValue.next {
        if retValue.forwarddecl == 1 {
            break;
        }
    }
    return retValue;
}

//
// Creates a new object
//
func NewObject(name string, packagename string, class uint64) *ObjectDesc {
    var adr uint64 = Alloc(OBJECT_SIZE);
    var obj *ObjectDesc = Uint64ToObjectDescPtr(adr);
    obj.name = name; //TODO: Copy string?
    obj.packagename = packagename; //Copy string?
    obj.class = class;
    obj.objtype = nil;
    obj.ptrtype = 0;
    obj.next = nil;
    return obj;
}

//
// Creates a new type
//
func NewType(name string, packagename string, forwarddecl byte, len uint64, basetype *TypeDesc) *TypeDesc {
    var adr uint64 = Alloc(TYPE_SIZE);
    var objtype *TypeDesc = Uint64ToTypeDescPtr(adr);
    objtype.name = name; //TODO: Copy string?
    objtype.forwarddecl = forwarddecl;
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

//
// Prints a formatted output of a given list of objects, including type, size etc.
//
func PrintObjects(list *ObjectDesc) {
    var o *ObjectDesc;
    for o = list; o != nil; o = o.next {
        PrintString("Object ");
        if StringLength(o.packagename) != 0 {
            PrintString(o.packagename);
            PrintChar('.');
        }
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
        PrintString(", offset: ");
        PrintNumber(GetObjectOffset(o, list));
        PrintString(")\n");
    }
}

//
// Prints a formatted output of a given list of types, including form, size etc.
//
func PrintTypes(list *TypeDesc) {
    var t *TypeDesc;
    var o *ObjectDesc;
    for t = list; t != nil; t = t.next {
        PrintString("Type ");
        if StringLength(t.packagename) != 0 {
            PrintString(t.packagename);
            PrintChar('.');
        }
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
            } else {
                PrintString("<unknown>");
            }
            PrintString(", size: ");
            PrintNumber(GetObjectSize(o));
            PrintString(", offset: ");
            PrintNumber(GetObjectOffset(o, t.fields));
            PrintString(")\n");
        }
    }
}
