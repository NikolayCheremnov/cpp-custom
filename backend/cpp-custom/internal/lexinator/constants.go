package lexinator

// constants

const MaxLexLen = 20
const MaxModuleLen = 10000

// lexeme type codes:

// 1)
const Id = 1
const IntConst = 2
const Assignment = 3

// 2) keywords
const Main = 1
const Int = 5
const Short = 6
const Long = 7
const Bool = 8
const For = 9
const Const = 10
const Void = 11

// 3) special signs
const Semicolon = 12
const Comma = 13
const OpeningBracket = 14
const ClosingBracket = 15
const OpeningBrace = 16
const ClosingBrace = 17

// 4) operations
const Plus = 18
const Minus = 19
const Mul = 20
const Div = 21
const Mod = 22
const Equ = 23
const LessEqu = 24
const MoreEqu = 25
const Less = 26
const More = 27

// 5) special types
const End = 0
const Err = -1
