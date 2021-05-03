// first good test
int main() {
    /*write code here
        ...
    */ 
    // assigments
    int variable = 20;
    short int _variable2 = 30;      // #1 invalid identify
    long int variable3ы = 404.3;    // #2 invalid identify and dot
    const c = 1;
    int tooLongIdentifyTooLongIdentify = 1; // #
    char c = 'c';                   // #3 char type not added => ' ' absent

    // proc call
    proc(variable, variable2);
    
    if (variable > 10)              // #4 if not added in this version
        variable -= 1;

    // expressions
    variable = c - variable2;
    variable3 = (variable / variable2) % c - 10;

    // cycle 1
    for(int i = 1; i <= 10; i = i + 1) {
        f = f * i;
    }

    // cycle 2
    int k = 20;
    while (k <= 20) {   // #5 while not added
        k--;            // #6 -- not added
    }

    // assigments      
    int 1a = 22;        // #7 invalid id
    int #a, _b, .c;     // # 8 bad id
    float var = 0;      // # 9 float not added
    
}

// procedure
void proc(int a, b) {
    int c = a + b;
}

// # 10 - open multiline comment
/* end of the code



