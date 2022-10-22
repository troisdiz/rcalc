// Rcalc.g4
grammar Rcalc;

// Tokens
DOT: '.' ;
INT_NUMBER: [+-]?[0-9]+ ;
DECIMAL_NUMBER: [+-]?[0-9]*DOT[0-9]+ |  [+-]?[0-9]+(DOT[0-9]*)? ;
SCIENTIFIC_NUMBER: [+-]?[0-9]*(DOT[0-9])?[eE][+-]?[0-9]+ ;
OP_ADD: '+';
OP_SUB: '-';
OP_MUL: '*';
OP_DIV: '/';
OP_TEST_EQUAL: '==';
OP_TEST_LT: '<';
OP_TEST_GT: '>';
OP_TEST_LET: '<=';
OP_TEST_GET: '>=';

DQUOTE: '"';
QUOTE: '\'';

CURLY_OPEN: '{';
CURLY_CLOSE: '}';

BRACKET_OPEN: '[';
BRACKET_CLOSE: ']';

PROG_OPEN: '<<';
PROG_CLOSE: '>>';

KW_START: 'start';
KW_FOR: 'for';
KW_NEXT: 'next';

KW_IF: 'if';
KW_THEN: 'then';
KW_ELSE: 'else';
KW_END: 'end';

NAME: [a-zA-Z][a-zA-Z0-9]*;

WHITESPACE: [ \r\n\t]+ -> skip;

// Rules
start : instr+ EOF;

instr
    : action_or_var_call         # InstrActionOrVarCall
    | op                         # InstrOp
    | variable                   # InstrVariable
    | if_then_else               # InstIfThenElse
    | start_next_loop            # InstrStartNextLoop
    | for_next_loop              # InstrForNextLoop
    | program_declaration        # InstrProgramDeclaration
    | local_var_creation         # InstrLocalVarCreation
    ;

op
    : OP_ADD | OP_SUB | OP_MUL | OP_DIV
    | OP_TEST_EQUAL | OP_TEST_GT | OP_TEST_GET | OP_TEST_LT | OP_TEST_LET
    ;

if_then_else
    : KW_IF instr+ KW_THEN instr+ (KW_ELSE instr+)* KW_END ;

start_next_loop: KW_START instr+ KW_NEXT ;
for_next_loop: KW_FOR variableDeclaration instr+ KW_NEXT ;
program_declaration: PROG_OPEN instr+ PROG_CLOSE ;
local_var_creation
    : '->' variableDeclaration+ program_declaration # LocalVarCreationProgram
//    | '->' variableDeclaration+ identifier          # LocalVarCreationAlgebraicExpr
    ;

variableDeclaration: NAME #DeclarationVariable;

variable
    : number     # VariableNumber
    | identifier # VariableIdentifier
    | list       # VariableList
    | vector     # VariableVector
    ;

number
    : INT_NUMBER        # NumberInt
    | DECIMAL_NUMBER    # NumberDecimal
    | SCIENTIFIC_NUMBER # NumberScientific
    ;

identifier: QUOTE NAME QUOTE ;

list : CURLY_OPEN variable* CURLY_CLOSE ;

vector : BRACKET_OPEN variable* BRACKET_CLOSE ;

action_or_var_call: NAME;
