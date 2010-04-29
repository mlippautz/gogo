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
var OBJECT_SIZE uint64 = 32; // 4*8 bytes space for an object
var TYPE_SIZE uint64 = 32;  // 4*8 bytes space for a type

//
// Classes for objects
//
var CLASS_VAR uint64 = 1;
var CLASS_TYPE uint64 = 2;

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
// Prepend an object to the list
//
func PrependObject(object *ObjectDesc, list *ObjectDesc) *ObjectDesc {
    var tmpObject *ObjectDesc;
    if list != nil {
        list.next = object;
        tmpObject = list;
    } else {
        tmpObject = object;
    }
    return tmpObject;
}

//
// Prepend type to the list
// 
func PreprendType(objtype *TypeDesc, list *TypeDesc) *TypeDesc {
    var tmpObject *TypeDesc;
    if list != nil {
        list.next = objtype;
        tmpObject = list;
    } else {
        tmpObject = objtype;
    }
    return tmpObject;
}

//
// Fetches an object with a specific identifier
//
func GetObject(name string, list *ObjectDesc) *ObjectDesc {
    var tmpObject *ObjectDesc;
    for tmpObject = list; tmpObject != nil; tmpObject = list.next {
        if StringCompare(tmpObject.name,name) == 0 {
            break;
        }
    }
    return tmpObject;
}

//
// Fetches a type with a given name
//
func GetType(name string, list *TypeDesc) *TypeDesc {
    var tmpType *TypeDesc;
    for tmpType = list; tmpType != nil; tmpType = list.next {
        if StringCompare(tmpType.name,name) == 0 {
            break;
        }
    }
    return tmpType;
}

//
// Creates a new object
//
func NewObject(name string) *ObjectDesc {
    var adr uint64 = Alloc(OBJECT_SIZE);
    var obj *ObjectDesc = Uint64ToObjectDescPtr(adr);
    obj.next = nil;
    return obj;
}

//
// Creates a new type
//
func NewType(name string) *TypeDesc {
    var adr uint64 = Alloc(TYPE_SIZE);
    var objtype *TypeDesc = Uint64ToTypeDescPtr(adr);
    objtype.next = nil;
    return objtype;
}
