// Rcalc.g4
grammar Rcalc;

// Tokens
NUMBER: [0-9]+;
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

instr
    : identifier
    | action_or_var_call
    | ADD
    | SUB
    | MUL
    | DIV
    ;

identifier: QUOTE NAME QUOTE;

action_or_var_call: NAME;
