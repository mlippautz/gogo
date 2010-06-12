// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//
// GoGo list functions
//

package libgogo

//
// List data structure
//
type List struct {
    baseAddress uint64; //Where the list starts
    itemSize uint64; //Size of one item in bytes
    itemCount uint64; //How many items there currently are in the list
    capacity uint64; //How many items there can be max. (not to be changed from outside the library)
};

//
// Initializes the given list
//
func InitializeList(uninitializedList *List) {
    uninitializedList.itemSize = 8; //Default size 64 bits = 8 bytes
    uninitializedList.capacity = 16; //Allocate 16 items by default
    uninitializedList.baseAddress = Alloc(uninitializedList.capacity * uninitializedList.baseAddress);
    uninitializedList.itemCount = 0; //Reset item count (to zero)
}

//
// Adds an item (referred to via the given pointer) to the given list, increasing its capacity if required
//
func AddItem(list *List, valuePtr *uint64) {
    var newAddress uint64;
    var valueAddr uint64;
    if (list.capacity == list.itemCount) { //Grow list if its capacity doesn't suffice to add another item
        newAddress = Alloc(list.capacity * 2 * list.itemSize); //Double the capacity
        CopyMem(list.baseAddress, newAddress, list.capacity * list.itemSize); //Copy old list items
        list.baseAddress = newAddress; //Set new address as base address
        list.capacity = list.capacity * 2; //Update (increase) capacity
    }
    valueAddr = ToUint64FromUint64Ptr(valuePtr);
    CopyMem(valueAddr, list.baseAddress + list.itemSize * list.itemCount, list.itemSize); //Append the new value by copying its value into the memory of the corresponding list item
    list.itemCount = list.itemCount + 1; //Update item count
}

//
// Removes the i-th item of the given list and returns a new address the removed item has been copied to
// Note that this function call fails if there are less than i items in the list or if the list is empty
//
func RemoveItemAt(list *List, index uint64) *uint64 {
    var returnPtr *uint64;
    var returnAddr uint64;
    var newReturnPtr uint64;
    var actualReturnPtr *uint64;
    var i uint64;
    returnPtr = GetItemAt(list, index); //Get the correspondig value value from the list
    newReturnPtr = Alloc(list.itemSize); //Allocate memory to store the item to be removed
    returnAddr = ToUint64FromUint64Ptr(returnPtr);
    CopyMem(returnAddr, newReturnPtr, list.itemSize); //Save item to be removed in order to return it
    for i = index; i < list.itemCount - 1; i = i + 1 { //Remove item by moving the following ones "backwards" in order to fill the gap caused by the deleted item
        CopyMem(list.baseAddress + (i + 1) * list.itemSize, list.baseAddress + list.itemSize * i, list.itemSize); //Move item at position i + 1 to position i
    }
    list.itemCount = list.itemCount - 1; //Update (decrease) item count
    actualReturnPtr = ToUint64PtrFromUint64(newReturnPtr);
    return actualReturnPtr; //Return value removed from list
}

//
// Returns a pointer to the i-th item of a given list without removing it
// Note that this function call fails if there are less than i items in the list or if the list is empty
//
func GetItemAt(list *List, index uint64) *uint64 {
    var returnPtr *uint64;
    if (list.itemCount <= index) { //Check if there are enough items in the list
        PrintString("Tried to GetItemAt(");
        PrintNumber(index);
        PrintString(") from a list with only ");
        PrintNumber(list.itemCount);
        ExitError(" items - out of bounds error", 125);
    }
    returnPtr = ToUint64PtrFromUint64(list.baseAddress + index * list.itemSize);
    return returnPtr; //Get the value on position i of the list (e.g. its corresponding memory) into the return variable
}

//
// Returns the number of items in the given list
//
func GetListItemCount(list *List) uint64 {
    return list.itemCount;
}

//
// Returns the current capacity of the given list
//
func GetListCapacity(list *List) uint64 {
    return list.capacity;
}

//
// Sets the list's item size
//
func SetListItemSize(list *List, itemSize uint64) {
    list.itemSize = itemSize;
}
