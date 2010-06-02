// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

type ObjectDesc struct {
    Name string;
    PackageName string;
    Class uint64;
    ObjType *TypeDesc;
    PtrType uint64; //If 0, type objtype, otherwise type *objtype
    Next *ObjectDesc;
};

type TypeDesc struct {
    Name string;
    PackageName string;
    ForwardDecl uint64; //If 1, type is forward declared
    Form uint64;
    Len uint64;
    Fields *ObjectDesc;
    Base *TypeDesc;
    Next *TypeDesc;
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
var CLASS_PARAMETER uint64 = 3;

//
// Forms for types
//
var FORM_SIMPLE uint64 = 1;
var FORM_STRUCT uint64 = 2;
var FORM_ARRAY uint64 = 3;
var FORM_FUNCTION uint64 = 4;

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
        for tmpObj = list; tmpObj.Next != nil; tmpObj = tmpObj.Next {
        }
        tmpObj.Next = object;
        tmpObj = list;
    }
    return tmpObj;    

		//Alternate version: PrependObject:
		/*object.Next = list;
    return object;*/
}

//
// Appends a type to the list
// 
func AppendType(objtype *TypeDesc, list *TypeDesc) *TypeDesc {
    var tmpObjType *TypeDesc = objtype;
    if list != nil {
        for tmpObjType = list; tmpObjType.Next != nil; tmpObjType = tmpObjType.Next { }
        tmpObjType.Next = objtype;
        tmpObjType = list;
    }
    return tmpObjType;

    //Alternate version: PrependType:
    /*objtype.Next = list;
    return objtype;*/
}

//
// Adds a field in form of an object descriptor to the type descriptor given
//
func AddFields(object *ObjectDesc, objtype *TypeDesc) {
    objtype.Form = FORM_STRUCT;
    objtype.Fields = AppendObject(object, objtype.Fields);
}

//
// Returns 1 if the given (struct) type has a field with the given name, or 0 otherwise
//
func HasField(name string, objtype *TypeDesc) uint64 {
    var tmpObj *ObjectDesc;
    var retVal uint64 = 0;
    tmpObj = GetField(name, objtype);
    if tmpObj != nil {
        retVal = 1;
    }
    return retVal;
}

//
// Returns the given struct's field with the name given, or nil if a field with that name does not exist
//
func GetField(name string, objtype *TypeDesc) *ObjectDesc {
    var tmpObj *ObjectDesc;
    tmpObj = GetObject(name, "", objtype.Fields);
    return tmpObj;
}

//
// Adds a parameter in form of an object descriptor to the function (type descriptor) given 
//
func AddParameters(object *ObjectDesc, fcn *TypeDesc) {
    fcn.Form = FORM_FUNCTION;
    fcn.Fields = AppendObject(object, fcn.Fields);
    fcn.Len = fcn.Len + 1; //Increase number of parameters
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
        if objtype.ForwardDecl == 0 {
            if objtype.Form == FORM_SIMPLE {
                size = objtype.Len;
            }
            if objtype.Form == FORM_STRUCT {
                for tempobj = objtype.Fields; tempobj != nil; tempobj = tempobj.Next { //Sum of all fields
                    tmpSize = GetObjectSizeAligned(tempobj);
                    size = size + tmpSize; //Add size of each field
                }
            }
            if objtype.Form == FORM_ARRAY {
                tmpSize = GetTypeSize(objtype.Base); //Get unaligned size of base type
                size = objtype.Len * tmpSize; //Array length * size of one item
                size = AlignTo64Bit(size); //Align to 64 bit
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
// Returns a 64 bit aligned address of the given one
//
func AlignTo64Bit(adr uint64) uint64 {
    return ((adr + 7) / 8) * 8;
}

//
// Returns the aligned (memory) size of the given type in bytes
//
func GetTypeSizeAligned(objtype *TypeDesc) uint64 {
    var size uint64;
    size = GetTypeSize(objtype);
    size = AlignTo64Bit(size); //64 bit alignment
    return size;
}

//
// Returns the (memory) size required by the object in bytes, considering whether or not the object is a pointer
// Same semantics as sizeof(obj)
//
func GetObjectSize(obj *ObjectDesc) uint64 {
    var size uint64;
    if obj.PtrType == 1 {
       size = 8;
    } else { //Actual type
       size = GetTypeSize(obj.ObjType);
    }
    return size;
}

//
// Returns the aligned size required by the object in bytes, considering whether or not the object is a pointer
//
func GetObjectSizeAligned(obj *ObjectDesc) uint64 {
    var size uint64;
    size = GetObjectSize(obj);
    size = AlignTo64Bit(size); //64 bit alignment
    return size;
}

//
// Calculates the (memory) offset of a given field of the specified (struct) type in bytes
// Note that the calculated size always 64 bit aligned
//
func GetObjectOffset(obj *ObjectDesc, list *ObjectDesc) uint64 {
    var offset uint64 = 0;
    var tmp *ObjectDesc;
    var tmpSize uint64;
    for tmp = list; tmp != nil; tmp = tmp.Next {
        if tmp == obj {
            break;
        }
        tmpSize = GetObjectSizeAligned(tmp);
        offset = offset + tmpSize; //Add field size
    }
    if tmp == nil {
        offset = 0; //TODO: Raise error as obj is appearently not in the list
    }
    return offset;
}

//
// Calculates the (memory) offset of a given field of the specified struct type in bytes
// Note that the calculated size always 64 bit aligned
//
func GetFieldOffset(obj *ObjectDesc, objtype *TypeDesc) uint64 {
    var offset uint64;
    offset = GetObjectOffset(obj, objtype.Fields);
    return offset;
}

//
// Calculates the (memory) size required to store the list of objects given
// Note that the calculated size is always 64 bit aligned
//
func GetAlignedObjectListSize(objList *ObjectDesc) uint64 {
    var i uint64 = 0; //Init data segment size with 0 in case there are no global objects
    var j uint64;
    var lastObj *ObjectDesc = objList; //Initialize last object with the beginning of the global variable list
    var obj *ObjectDesc;
    for obj = objList; obj != nil; obj = obj.Next { //Calculate offset of last global variable //TODO: Do this more efficiently
        i = GetObjectOffset(obj, objList);
        lastObj = obj;
    }
    if lastObj != nil { //Add size of last object to its offset in order to get the total size required for the data segment
        j = GetObjectSizeAligned(lastObj);
        i = i + j;
    }
    return i;
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
    for tmpObject = list; tmpObject != nil; tmpObject = tmpObject.Next {
        nameCompare = StringCompare(tmpObject.Name,name);
        strLen = StringLength(tmpObject.PackageName);
        packageNameCompare = StringCompare(tmpObject.PackageName,packagename);
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
func GetType(name string, packagename string, list *TypeDesc, includeforward uint64) *TypeDesc {
    var tmpType *TypeDesc;
    var retValue *TypeDesc = nil;
    var nameCompare uint64;
    var strLen uint64;
    var packageNameCompare uint64;
    for tmpType = list; tmpType != nil; tmpType = tmpType.Next {
        nameCompare = StringCompare(tmpType.Name,name);
        strLen = StringLength(tmpType.PackageName);
        packageNameCompare = StringCompare(tmpType.PackageName,packagename);
        if (nameCompare == 0) && ((strLen == 0) || (packageNameCompare == 0)) { //Empty package name indicates internal types
            if (includeforward == 1) || ((includeforward == 0) && (tmpType.ForwardDecl == 0)) {
                retValue = tmpType;
                break;
            }
        }
    }
    return retValue;
}

//
// Tries to find the given package name in one of the objects of the given list
// Returns 1 if an object has been found, or 0 otherwise
//
func FindPackageName(packagename string, list *ObjectDesc) uint64 {
    var retVal uint64 = 0;
    var tmpObj *ObjectDesc;
    var packageNameCompare uint64;
    for tmpObj = list; tmpObj != nil; tmpObj = tmpObj.Next {
        packageNameCompare = StringCompare(tmpObj.PackageName, packagename);
        if packageNameCompare == 0 {
            retVal = 1;
            break;
        }
    }
    return retVal;
}

//
// Returns the first type from the list which is forward declared
//
func GetFirstForwardDeclType(list *TypeDesc) *TypeDesc {
    var retValue *TypeDesc;
    for retValue = list; retValue != nil; retValue = retValue.Next {
        if retValue.ForwardDecl == 1 {
            break;
        }
    }
    return retValue;
}

//
// Creates a new object
//
func NewObject(name string, packagename string, class uint64) *ObjectDesc {
    var adr uint64;
    var obj *ObjectDesc;
    adr = Alloc(OBJECT_SIZE);
    obj = Uint64ToObjectDescPtr(adr);
    obj.Name = name; //TODO: Copy string?
    obj.PackageName = packagename; //Copy string?
    obj.Class = class;
    obj.ObjType = nil;
    obj.PtrType = 0;
    obj.Next = nil;
    return obj;
}

//
// Creates a new type
//
func NewType(name string, packagename string, forwarddecl uint64, len uint64, basetype *TypeDesc) *TypeDesc {
    var adr uint64;
    var objtype *TypeDesc;
    adr = Alloc(TYPE_SIZE);
    objtype = Uint64ToTypeDescPtr(adr);
    objtype.Name = name; //TODO: Copy string?
    objtype.ForwardDecl = forwarddecl;
    objtype.PackageName = packagename; //TODO: Copy string?
    if basetype != nil {
        objtype.Form = FORM_ARRAY;
    } else {
        objtype.Form = FORM_SIMPLE;
    }
    objtype.Len = len;
    objtype.Next = nil;
    objtype.Fields = nil;
    objtype.Base = basetype;
    return objtype;
}

//
// Prints a formatted output of a given list of objects, including type, size etc.
//
func PrintObjects(list *ObjectDesc) {
    var o *ObjectDesc;
    var strLen uint64;
    var tmp uint64;
    for o = list; o != nil; o = o.Next {
        PrintString("Object ");
        strLen = StringLength(o.PackageName);
        if strLen != 0 {
            PrintString(o.PackageName);
            PrintChar('.');
        }
        PrintString(o.Name);
        PrintString(" (type: ");
        if o.PtrType != 0 {
            PrintString("pointer to ");
        }
        if o.ObjType != nil {
            if o.ObjType.Base != nil {
                PrintString("array of ");
                PrintString(o.ObjType.Base.Name);
                PrintString(" of length ");
                PrintNumber(o.ObjType.Len);
                PrintString(", internally named ");
            }
            strLen = StringLength(o.ObjType.PackageName);
            if strLen != 0 {
                PrintString(o.ObjType.PackageName);
                PrintChar('.');
            }
            PrintString(o.ObjType.Name);
        } else {
            PrintString("<unknown>");
        }
        PrintString(", size: ");
        tmp = GetObjectSize(o);
        PrintNumber(tmp);
        PrintString(", offset: ");
        tmp = GetObjectOffset(o, list);
        PrintNumber(tmp);
        PrintString(")\n");
    }
}

//
// Prints a formatted output of a given list of types, including form, size etc.
//
func PrintTypes(list *TypeDesc) {
    var t *TypeDesc;
    var o *ObjectDesc;
    var strLen uint64;
    var tmp uint64;
    for t = list; t != nil; t = t.Next {
        PrintString("Type ");
        strLen = StringLength(t.PackageName);
        if strLen != 0 {
            PrintString(t.PackageName);
            PrintChar('.');
        }
        PrintString(t.Name);
        PrintString(" (size: ");
        tmp = GetTypeSize(t);
        PrintNumber(tmp);
        PrintString(")\n");
				for o = t.Fields; o != nil; o = o.Next {
            PrintString("  ");
            PrintString(o.Name);
            PrintString(" (type: ");
            if o.PtrType != 0 {
                PrintString("pointer to ");
            }
            if o.ObjType != nil {
                if o.ObjType.Base != nil {
                    PrintString("array of ");
                    PrintString(o.ObjType.Base.Name);
                    PrintString(" of length ");
                    PrintNumber(o.ObjType.Len);
                    PrintString(", internally named ");
                }
                strLen = StringLength(o.ObjType.PackageName);
                if strLen != 0 {
                    PrintString(o.ObjType.PackageName);
                    PrintChar('.');
                }
                PrintString(o.ObjType.Name);
            } else {
                PrintString("<unknown>");
            }
            PrintString(", size: ");
            tmp = GetObjectSize(o);
            PrintNumber(tmp);
            PrintString(", offset: ");
            tmp = GetObjectOffset(o, t.Fields);
            PrintNumber(tmp);
            PrintString(")\n");
        }
    }
}

//
// Prints a formatted output of a given list of functions, including parameters etc.
//
func PrintFunctions(list *TypeDesc) {
    var t *TypeDesc;
    var o *ObjectDesc;
    var strLen uint64;
    var tmp uint64;
    for t = list; t != nil; t = t.Next {
        PrintString("Function ");
        PrintString(t.PackageName);
        PrintChar('.');
        PrintString(t.Name);
        PrintString(" (number of parameters: ");
        PrintNumber(t.Len);
        PrintString(")\n");
            for o = t.Fields; o != nil; o = o.Next {
            PrintString("  ");
            PrintString(o.Name);
            PrintString(" (type: ");
            if o.PtrType != 0 {
                PrintString("pointer to ");
            }
            if o.ObjType != nil {
                if o.ObjType.Base != nil {
                    PrintString("array of ");
                    PrintString(o.ObjType.Base.Name);
                    PrintString(" of length ");
                    PrintNumber(o.ObjType.Len);
                    PrintString(", internally named ");
                }
                strLen = StringLength(o.ObjType.PackageName);
                if strLen != 0 {
                    PrintString(o.ObjType.PackageName);
                    PrintChar('.');
                }
                PrintString(o.ObjType.Name);
            } else {
                PrintString("<unknown>");
            }
            PrintString(", size: ");
            tmp = GetObjectSize(o);
            PrintNumber(tmp);
            PrintString(", offset: ");
            tmp = GetObjectOffset(o, t.Fields);
            PrintNumber(tmp);
            PrintString(")\n");
        }
    }
}
