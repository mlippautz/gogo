// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo list functions
//

package libgogo

//
// String list data structure
//
type StringList struct {
    internalList List; //List used internally
};

//
// Initializes the given string list
//
func InitializeStringList(uninitializedList *StringList) {
    InitializeList(&uninitializedList.internalList, 16); //A string has a size of 16 bytes
}

//
// Adds an item to the given string list, increasing its capacity if required
//
func AddStringItem(list *StringList, value string) {
    var newAddress uint64;
    var ptr *uint64;
    newAddress = ToUint64FromStringPtr(&value);
    ptr = ToUint64PtrFromUint64(newAddress);
    AddItem(&list.internalList, ptr);
}

//
// Returns the i-th item of a given string list without removing it
// Note that this function call fails if there are less than i items in the list or if the list is empty
//
func GetStringItemAt(list *StringList, index uint64) string {
    var returnPtr *uint64;
    var returnAddr uint64;
    var returnValue string;
    var temp uint64;
    returnPtr = GetItemAt(&list.internalList, index);
    returnAddr = ToUint64FromUint64Ptr(returnPtr);
    temp = ToUint64FromStringPtr(&returnValue);
    CopyMem(returnAddr, temp, 16); //Copy the value into the return variable
    return returnValue;
}

//
// Returns the number of items in the given list
//
func GetStringListItemCount(list *StringList) uint64 {
    return list.internalList.itemCount;
}

//
// Returns the current capacity of the given list
//
func GetStringListCapacity(list *StringList) uint64 {
    return list.internalList.capacity;
}
