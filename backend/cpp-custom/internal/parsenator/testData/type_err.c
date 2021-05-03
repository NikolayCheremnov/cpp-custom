// first good test

const int a = 2;

void main() {
    /*write code here
        ...
    */
    // assigments
    int variable = 20;
    short variable2 = 30;   // type err
    long int variable3 = 404;
    const int c = 1;
    int f = 1;

    // proc call
    proc(variable, variable2);

    // expressions
    variable = c - variable2;
    variable3 = (variable / variable2) % c - 10;
    // cycle
    for(int i = 1; i < 10; i = i + 1) {
        f = f * i;
    }
}

// procedure
void proc(int a, int b) {
    int c = a + b;
}

//final comment without \n