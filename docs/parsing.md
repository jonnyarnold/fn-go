# Parsing

fn runs a single-pass parser to convert an input string into an Abstract Syntax Tree (AST). To give an example, the parser converts:

```fn
x = 2 + 2
```

into the AST:

```
FunctionCall(Id(=), Id(x), FunctionCallExpression(Id(+), Number(2), Number(2)))
```

It does it in two steps: *tokenisation* and *parsing*:

## Tokenisation

The tokeniser runs through a set of rules to split the code into tokens:

1. Check if it is one of the *basic tokens*: these are tokens that do not have a value, and consist of a single character.
2. Check if we have a string. A string starts with a `"` and its value is the string until the next `"`.
3. Check if we have a number. A number starts with a numeric character and its value is the code until the next character that is not numeric or `.`.
4. Check for *symbol infix operators*: these are infix operators that consist of a single character. Output an infix operator token.
5. When these fail, eat the code until we reach a newline, comment, basic token, symbol infix operator, `"` or `.`:
    - If it is `true` or `false`, output a boolean token with the correct value.
    - If it is a keyword, output the token for that keyword.
    - If it is a *string infix operator*, output an infix operator token.
    - Otherwise, output an identifier token.

In our previous example,

```fn
x = 2 + 2
```

becomes

```
identifier[x] infix_operator[=] number[2] infix_operator[+] number[2]
```

## Parsing

The parser takes the token stream and turns it into an Abstract Syntax Tree (AST). 

### BNF-ish

```
code = (primary)*
primary = end | brackets | value
end = END_STATEMENT
brackets = BRACKET_OPEN primary BRACKET_CLOSE
value = literal | identifier | block | when | functionDefinition | functionCall | infixOp

literal = number | string | boolean

block = BLOCK_OPEN code BLOCK_CLOSE

when = WHEN BLOCK_OPEN (value block) BLOCK_CLOSE

functionDefinition = params block
params = BRACKET_OPEN (identifier (COMMA)?)* BRACKET_CLOSE

functionCall = Identifier args
args = BRACKET_OPEN (value (COMMA)?)* BRACKET_CLOSE

infixOp = value infixOp value
```

### Decision Tree

To explain the rules, here is some notation!

```
lowercase =              ## Decision point
$token                   ## Token
else                     ## If all other tokens do not match
$token => x              ## If $token is found, go to x
Expression               ## Output Expression
[Note]                   ## Operation that isn't able to be described in the notation!
```

It uses the following rules:

```
primary = 
    $end_statement                                   => [No expression]
    $identifier $number $string $boolean $( $when ${ => value
    else                                             => [Error]

value =
    $identifier
        $(    => functionCall
        else  => Identifier

    $number  => Number
    $string  => String
    $boolean => Boolean

    ${       => block
    $when    => when

    $(
        $} after $) => functionDefinition
        else        => brackets

    else => [Error]
[NB. value always calls infixOperator afterward]

infixOperator =
    $infix_operator => FunctionCall [with value afterward]
    else            => [No expression]

functionCall =
    $identifier $( => FunctionCall [with params]

params =
    $(
        $)   => [End loop, return params]
        $identifier
            $,   => [Add to params, loop]
            else => [Error]
        else => [Error]
    else => [Error]

block =
    ${
        $}   => [End loop, return Body]
        else => primary [Add to Body, loop]
    else => [Error]

when = 
    $when
        $block_open
            $block_close => [End loop, return Conditional]
            else
                value
                    block => [Add to Conditional, loop]

functionDefinition = 
    args
        block => FunctionDefinition

args = 
    $(
        $)   => [End loop, return args]
        value
            $,   => [Add to args, loop]
            else => [Error]
        else => [Error]
    else => [Error]

brackets =
    $(
        primary
            $)   => [Primary]
            else => [Error]
    else => [Error]
```