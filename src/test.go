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

var initTest byte = 1 + 2 * 3;
var initTest2 uint64 = initTest + 1;
var fooArrInst [10]foo;
var testInst test;
var testPtr *test;

func test(bla uint64, blub string) {
    var a uint64 = 5;
    var b uint64 = 1;
    var c string = "Test";
    var d *string;
    var e byte;
    var f byte;
    b = e + f + 1 - a;
    a = 1 + b * a;
    c = fooArrInst[a].bar2;
    c[1] = testPtr.barPtr.bar2[1];
    d[2] = 'x';
    a = fooArrInst[b].bar1 + testInst.barPtr.bar1;
}
