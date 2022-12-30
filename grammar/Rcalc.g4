// Rcalc.g4
grammar Rcalc;

// Tokens
DOT: '.' ;
INT_NUMBER: [-]?[0-9]+ ;
DECIMAL_NUMBER: [-]?[0-9]*DOT[0-9]+ |  [-]?[0-9]+(DOT[0-9]*)? ;
SCIENTIFIC_NUMBER: [-]?[0-9]*(DOT[0-9])?[eE][+-]?[0-9]+ ;

OP_ADD: '+';
OP_SUB: '-';
OP_MUL: '*';
OP_DIV: '/';
OP_POW: '^' ;

OP_TEST_EQUAL: '==';
OP_TEST_LT: '<';
OP_TEST_GT: '>';
OP_TEST_LET: '<=';
OP_TEST_GET: '>=';

DQUOTE: '"';
QUOTE: '\'';
COMMA: ',';

PAREN_OPEN: '(' ;
PAREN_CLOSE: ')' ;

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

NAME: [a-zA-Z_][a-zA-Z0-9_]*;

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

program_declaration:
    PROG_OPEN instr+ PROG_CLOSE  # ProgramDeclaration
    ;

local_var_creation
    : '->' variableDeclaration+ program_declaration # LocalVarCreationProgram
//    | '->' variableDeclaration+ identifier          # LocalVarCreationAlgebraicExpr
    ;

variableDeclaration: NAME #DeclarationVariable;

variable
    : number                      # VariableNumber
    | quoted_algebraic_expression # VariableAlgebraicExpression
    | list                        # VariableList
    | vector                      # VariableVector
    ;

number
    : INT_NUMBER        # NumberInt
    | DECIMAL_NUMBER    # NumberDecimal
    | SCIENTIFIC_NUMBER # NumberScientific
    ;

quoted_algebraic_expression: QUOTE alg_expression QUOTE ;

alg_expression
   : alg_mulExpression ((OP_ADD | OP_SUB) alg_mulExpression)* # AlgExprAddSub
   ;

alg_mulExpression
   : alg_powExpression ((OP_MUL | OP_DIV) alg_powExpression)* # AlgExprMulDiv
   ;

alg_powExpression
   : alg_signedAtom (OP_POW alg_signedAtom)* #AlgExprPow
   ;

alg_signedAtom
   : OP_ADD alg_signedAtom # AlgExprAddSignedAtom
   | OP_SUB alg_signedAtom # AlgExprSubSignedAtom
   | alg_func_call         # AlgExprFuncAtom
   | alg_atom              # AlgExprAtom
   ;

alg_atom
   : number                                # AlgExprNumber
   | alg_variable                          # AlgExprVariable
   | PAREN_OPEN alg_expression PAREN_CLOSE # AlgExprParen
   ;

alg_variable
   : NAME
   ;

alg_func_call
   : function_name=NAME PAREN_OPEN alg_expression (COMMA alg_expression)* PAREN_CLOSE # AlgExprFuncCall
   ;

list : CURLY_OPEN variable* CURLY_CLOSE ;

vector : BRACKET_OPEN variable* BRACKET_CLOSE ;

action_or_var_call: NAME;
