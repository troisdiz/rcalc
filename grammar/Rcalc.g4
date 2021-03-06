// Rcalc.g4
grammar Rcalc;

// Tokens
DOT: '.' ;
INT_NUMBER: [+-]?[0-9]+ ;
DECIMAL_NUMBER: [+-]?[0-9]*DOT[0-9]+ |  [+-]?[0-9]+(DOT[0-9]*)? ;
SCIENTIFIC_NUMBER: [+-]?[0-9]*(DOT[0-9])?[eE][+-]?[0-9]+ ;
ADD: '+';
SUB: '-';
MUL: '*';
DIV: '/';

DQUOTE: '"';
QUOTE: '\'';

CURLY_OPEN: '{';
CURLY_CLOSE: '}';

BRACKET_OPEN: '[';
BRACKET_CLOSE: ']';

PROG_OPEN: '<<';
PROG_CLOSE: '>>';

NAME: [a-zA-Z][a-zA-Z0-9]*;

WHITESPACE: [ \r\n\t]+ -> skip;


// Rules
start : instr+ EOF;

number
    : INT_NUMBER # NumberInt
    | DECIMAL_NUMBER # NumberDecimal
    | SCIENTIFIC_NUMBER # NumberScientific
    ;

instr
    : identifier                 # InstrIndentifier
    | action_or_var_call         # InstrActionOrVarCall
    | op=(ADD | SUB | MUL | DIV) # InstrOp
    | number                     # InstrNumber
    ;

identifier: QUOTE NAME QUOTE;

action_or_var_call: NAME;
