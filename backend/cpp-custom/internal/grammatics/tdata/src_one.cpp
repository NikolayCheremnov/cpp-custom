// first good test

const int a = 2;

main() {
    /*write code here
        ...
    */
    // assigments
    int v = 2;
    short int v = 3;
    long int v = 4;
    const int c = 1;
    int f = 1;

    // proc call
    p(v, v);

    // expressions
    v = c - v;
    v = (v / v) % c - 1;
    // cycle
    for(int i = 1; i < 1; i = i + 1) {
        f = f * i;
    }
}

// procedure
void p(int a, int b) {
    int c = a + b;
}

//final comment without \n