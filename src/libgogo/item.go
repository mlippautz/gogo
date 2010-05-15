// Copyright 2010 The GoGo Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package libgogo

type Item struct {
    Mode uint64;
    Itemtype *TypeDesc;
    A uint64;
    R uint64;
    Global uint64; //If 1, the variable is global, otherwise it is local
};

//
// Pseudo constants that specify the descriptor sizes 
//
var ITEM_SIZE uint64 = 40; //5*8 bytes space for an object

//
// Modes for items
//
var MODE_VAR uint64 = 1;
var MODE_CONST uint64 = 2;
var MODE_REG uint64 = 3;

//
// Convert the uint64 value (returned from malloc) to a real item address
//
func Uint64ToItemPtr(adr uint64) *Item;

//
// Convert the item to a uint64 value
//
func ItemPtrToUint64(item *Item) uint64;

//
// Creates a new, uninitialized item
//
func NewItem() *Item {
    var adr uint64;
    var item *Item;
    adr = Alloc(ITEM_SIZE);
    item = Uint64ToItemPtr(adr);
    return item;
}

//
// Sets the given item's properties
//
func SetItem(item *Item, mode uint64, itemtype *TypeDesc, a uint64, r uint64, global uint64) {
    item.Mode = mode;
    item.Itemtype = itemtype;
    item.A = a;
    item.R = r;
    item.Global = global;
}

//
// Pushes the specified item to the given stack
//
func PushItem(item *Item, stack *Stack) {
    var itemPtr uint64;
    itemPtr = ItemPtrToUint64(item);
    Push(stack, itemPtr);
}

//
// Pops an item from the given stack
//
func PopItem(stack *Stack) *Item {
    var itemPtr uint64;
    var item *Item;
    itemPtr = Pop(stack);
    item = Uint64ToItemPtr(itemPtr);
    return item;
}
