package main

func main() {
    var x string = "";
    libgogo.StringAppend(&x, "Hello");
    libgogo.StringAppend(&x, " world!\n");
    libgogo.PrintString(x);
    libgogo.PrintNumber(123);
    libgogo.PrintChar('\n');
}
