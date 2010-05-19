package test

type foo struct {
    bar1 uint64;
    bar2 string;
    bar3 *uint64;
};

type test struct {
    bar foo;
    bazz [20]string;
    barPtr *foo;
};

var fooArrInst [10]foo;
var testInst test;

func test() {
    var a uint64;
    //a = a + 1;
    //var b uint64;
    //a = testInst.bazz[1];
    //a = testInst.bazz[a];
    //a = testInst.barPtr.bar2;
    //a = fooArrInst[5].bar3;
    a = fooArrInst[a].bar3;
    //testInst.bazz[2] = fooArrInst[5].bar3;
}
