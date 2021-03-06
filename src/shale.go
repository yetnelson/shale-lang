package main

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "os/exec"
    "os"
    "math"
    "strings"
    "bufio"
)

const ShaleVersion = "0.1.0"

func main() {
    var stack []interface{}
    var retPos []int
    var op func([]interface{}) ([]interface{}, string)
    ops := map[string]func([]interface{}) ([]interface{}, string) {
        // Stack manipulation
        ",": opPop,
        // Maths
        "+": opAdd,
        "-": opSub,
        "*": opMul,
        "/": opDiv,
        "%": opMod,
        "^": opPow,
        "sqrt": opSqrt,
        "sin": opSin,
        "cos": opCos,
        "tan": opTan,
        // IO
        "print": opPrint,
        "puts": opPuts,
        "stderr": opStderr,
        "input": opInput,
        "read": opRead,
        "write": opWrite,
        "sys": opSys,
        // Comparison
        "==": opEqu,
        "<": opLt,
        ">": opGt,
        "<=": opLtEqu,
        ">=": opGtEqu,
        "!=": opNotEqu,
        "!": opNot,
        "cmp": opCmp,
        "less": opLess,
        "greater": opGreater,
        "equal": opEqual,
        "assert": opAssert,
        // Data manipulation
        "len": opLen,
        "join": opJoin,
        "split": opSplit,
        "replace": opReplace,
        "getChar": opGetChar,
        "vec": opVec,
        "get": opGet,
        "set": opSet,
        "append": opAppend,
        "type": opType,
        "extend": opExtend,
        "expand": opExpand,
        "sum": opSum,
        "upcase": opUpcase,
        "downcase": opDowncase,
    }
    var variable string
    vars := make(map[string]interface{})
    vars["_version"] = ShaleVersion
    vars["_err"] = ""
    procedures := make(map[string]int)
    if len(os.Args) < 2 {
        shalePanic("No input file")
    }
    code, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        shalePanic("Could not read input file")
    }
    tokens := tokenize(string(code))
    for i := 0; i < len(tokens); i++ {
        if tokens[i][len(tokens[i])-1] == '{' {
            procedures[tokens[i][0:len(tokens[i])-1]] = i
        }
    }
    for i := 0; i < len(tokens); i++ {
        integer, err := strconv.Atoi(tokens[i])
        if err == nil {
            stack = append(stack, integer)
        } else {
            float, err := strconv.ParseFloat(tokens[i], 64)
            if err == nil {
                stack = append(stack, float)
            } else if tokens[i][0] == '"' {
                str := tokens[i][1:len(tokens[i])-1]
                stack = append(stack, str)
            } else if tokens[i][0] == '@' {
                retPos = append(retPos, i)
                i = procedures[tokens[i][1:]]
            } else if tokens[i][0] == '?' {
                // If-jumps
                if len(stack) != 1 {
                    shalePanic("?...: Invalid number of operands --- required 1")
                } else if stack[0] == true {
                    retPos = append(retPos, i)
                    i = procedures[tokens[i][1:]]
                }
                stack = nil
            } else if tokens[i][0] == '$' {
                variable = tokens[i][1:]
                stack = append(stack, vars[variable])
            } else if tokens[i][0] == ':' {
                variable = tokens[i][1:]
                switch stack[0].(type) {
                    case int:
                        vars[variable] = stack[0].(int)
                    case float64:
                        vars[variable] = stack[0].(float64)
                    case string:
                        vars[variable] = stack[0].(string)
                    case []interface{}:
                        vars[variable] = stack[0].([]interface{})
                }
                stack = nil
            } else if tokens[i][len(tokens[i])-1] == '{' {
                os.Exit(0)
            } else {
                switch tokens[i] {
                    case "Null":
                        stack = append(stack, "")
                    case "True":
                        stack = append(stack, true)
                    case "False":
                        stack = append(stack, false)
                    case "}":
                        i = retPos[len(retPos) - 1]
                        retPos = retPos[:len(retPos) - 1]
                    case ";":
                        stack = nil

                    default:
                        op = ops[tokens[i]]
                        if op == nil {
                            shalePanic("Invalid operator " + tokens[i])
                        }
                        stack, vars["_err"] = op(stack)
                }
            }
        }
    }
}

func tokenize(code string) []string {
    tokens := []string{}
    token := ""
    var tokenSet bool
    var comment bool
    inString := false
    for i := 0; i < len(code); i++ {
        if comment == true {
            if code[i] == '#' {
                comment = false
            }
        } else if code[i] == '#' && inString == false {
            comment = true
        } else if code[i] == '"' {
            token += "\""
            if inString == true {
                inString = false
            } else if inString == false {
                inString = true
                tokenSet = true
            }
        } else if code[i] == ' ' {
            if inString == false {
                if tokenSet == true {
                    tokens = append(tokens, token)
                    token = ""
                    tokenSet = false
                }
            } else if inString == true {
                token += " "
            }
        } else if code[i] == '\n' {
            if inString == true {
                token += "\n"
            } else {
                if tokenSet == true {
                    tokens = append(tokens, token)
                    token = ""
                    tokenSet = false
                }
            }
        } else {
            tokenSet = true
            token += string(code[i])
        }
        if i == len(code) - 1 && code[i] != ' ' {
            if tokenSet == true {
                tokens = append(tokens, token)
                tokenSet = false
            }
        }
    }
    return tokens
}

func shalePanic(msg string) {
    fmt.Println("Panic: " + msg)
    os.Exit(1)
}

func opPrint(stack []interface{}) ([]interface{}, string) {
    for i := 0; i < len(stack); i++ {
        if stack[i] == nil {
            fmt.Print("Null")
        } else if stack[i] == true {
            fmt.Print("True")
        } else if stack[i] == false {
            fmt.Print("False")
        } else {
            fmt.Print(stack[i])
        }
    }
    return nil, ""
}

func opPuts(stack []interface{}) ([]interface{}, string) {
    for i := 0; i < len(stack); i++ {
        if stack[i] == nil {
            fmt.Println("Null")
        } else if stack[i] == true {
            fmt.Println("True")
        } else if stack[i] == false {
            fmt.Println("False")
        } else {
            fmt.Println(stack[i])
        }
    }
    return nil, ""
}

func opGet(stack []interface{}) ([]interface{}, string) {
    var idx int
    switch stack[1].(type) {
        case int:
            idx  = stack[1].(int)
        default:
            shalePanic("get: Inoperable type")
    }
    switch stack[0].(type) {
        case []interface{}:
            vec := stack[0].([]interface{})
            var element interface{}
            if idx < 0 {
                element = vec[len(vec)+idx]
            } else {
                element = vec[idx]
            }
            return []interface{}{element}, ""
        case string:
            s := stack[0].(string)
            return []interface{}{string([]rune(s)[idx]),}, ""
        default:
            shalePanic("get: Inoperable type")
    }
    return nil, ""
}

func opSet(stack []interface{}) ([]interface{}, string) {
    vec := stack[0].([]interface{})
    idx := stack[1]
    if idx.(int) < 0 {
        vec[len(vec)+idx.(int)] = stack[2]
    } else {
        vec[idx.(int)] = stack[2]
    }
    return []interface{}{vec,}, ""
}

func opAdd(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(int) + stack[1].(int),}, ""
                case float64:
                    return []interface{}{float64(stack[0].(int)) + stack[1].(float64),}, ""
                default:
                    shalePanic("+: Inoperable type")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(float64) + float64(stack[1].(int)),}, ""
                case float64:
                    return []interface{}{stack[0].(float64) + stack[1].(float64),}, ""
                default:
                    shalePanic("+: Inoperable type")
            }
        case string:
            switch stack[1].(type) {
                case string:
                    s := stack[0].(string) + stack[1].(string)
                    return []interface{}{s,}, ""
                default:
                    shalePanic("+: Inoperable type")
            }
        default:
            shalePanic("+: Inoperable type")
    }
    return nil, ""
}

func opSub(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(int) - stack[1].(int),}, ""
                case float64:
                    return []interface{}{float64(stack[0].(int)) - stack[1].(float64),}, ""
                default:
                    shalePanic("-: Inoperable type")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(float64) - float64(stack[1].(int)),}, ""
                case float64:
                    return []interface{}{stack[0].(float64) - stack[1].(float64),}, ""
                default:
                    shalePanic("-: Inoperable type")
            }
        default:
            shalePanic("-: Inoperable type")
    }
    return nil, ""
}

func opMul(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(int) * stack[1].(int),}, ""
                case float64:
                    return []interface{}{float64(stack[0].(int)) * stack[1].(float64),}, ""
                default:
                    shalePanic("*: Inoperable type")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(float64) * float64(stack[1].(int)),}, ""
                case float64:
                    return []interface{}{stack[0].(float64) * stack[1].(float64),}, ""
                default:
                    shalePanic("*: Inoperable type")
            }
        case string:
            switch stack[1].(type) {
                case int:
                    s := strings.Repeat(stack[0].(string), stack[1].(int))
                    return []interface{}{s,}, ""
                default:
                    shalePanic("*: Inoperable type")
            }
        default:
            shalePanic("*: Inoperable type")
    }
    return nil, ""
}

func opDiv(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(int) / stack[1].(int),}, ""
                case float64:
                    return []interface{}{float64(stack[0].(int)) / stack[1].(float64),}, ""
                default:
                    shalePanic("/: Inoperable type")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(float64) / float64(stack[1].(int)),}, ""
                case float64:
                    return []interface{}{stack[0].(float64) / stack[1].(float64),}, ""
                default:
                    shalePanic("/: Inoperable type")
            }
        default:
            shalePanic("/: Inoperable type")
    }
    return nil, ""
}

func opMod(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    return []interface{}{stack[0].(int) % stack[1].(int),}, ""
                default:
                    shalePanic("%: Inoperable type")
            }
        default:
            shalePanic("%: Inoperable type")
    }
    return nil, ""
}

func opSys(stack []interface{}) ([]interface{}, string) {
    if len(stack) != 1 {
        shalePanic("sys: Incorrect number of operands --- expected 1")
    }
    output, err := exec.Command(stack[0].(string)).Output()
    return []interface{}{string(output), err}, ""
}

func opJoin(stack []interface{}) ([]interface{}, string) {
    var joinStr string
    switch stack[0].(type) {
        case []interface{}:
            if len(stack) == 2 {
                switch stack[1].(type) {
                    case string:
                        joinStr = stack[1].(string)
                    default:
                        shalePanic("join: Inoperable type")
                }
            }
        default:
            shalePanic("join: Inoperable type")
    }
    var strSlice []string
    for i := 0; i < len(stack[0].([]interface{})); i++ {
        switch stack[0].([]interface{})[0].(type) {
            case string:
                strSlice = append(strSlice, stack[0].([]interface{})[i].(string))
            default:
                shalePanic("join: Inoperable type")
        }
    }
    s := strings.Join(strSlice, joinStr)
    return []interface{}{s,}, ""
}

func opSqrt(stack []interface{}) ([]interface{}, string) {
    return []interface{}{math.Sqrt(float64(stack[0].(int)))}, ""
}

func opEqu(stack []interface{}) ([]interface{}, string) {
    if stack[0] == stack[1] {
        return []interface{}{true,}, ""
    } else {
        return []interface{}{false,}, ""
    }
}

func opLt(stack []interface{}) ([]interface{}, string) {
    var result bool
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(int) < stack[1].(int)
                case float64:
                    result = float64(stack[0].(int)) < stack[1].(float64)
                case string:
                    shalePanic("<: Cannot compare int with str")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(float64) < float64(stack[1].(int))
                case float64:
                    result = stack[0].(float64) < stack[1].(float64)
                case string:
                    shalePanic("<: Cannot compare float with str")
            }
        case string:
            switch stack[1].(type) {
                case string:
                    result = stack[0].(string) < stack[1].(string)
                case int:
                    shalePanic("<: Cannot compare str with int")
                case float64:
                    shalePanic("<: Cannot compare str with float")
            }
    }
    return []interface{}{result,}, ""
}

func opGt(stack []interface{}) ([]interface{}, string) {
    var result bool
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(int) > stack[1].(int)
                case float64:
                    result = float64(stack[0].(int)) > stack[1].(float64)
                case string:
                    shalePanic(">: Cannot compare int with str")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(float64) > float64(stack[1].(int))
                case float64:
                    result = stack[0].(float64) > stack[1].(float64)
                case string:
                    shalePanic(">: Cannot compare float with str")
            }
        case string:
            switch stack[1].(type) {
                case string:
                    result = stack[0].(string) > stack[1].(string)
                case int:
                    shalePanic(">: Cannot compare str with int")
                case float64:
                    shalePanic(">: Cannot compare str with float")
            }
    }
    return []interface{}{result,}, ""
}

func opLtEqu(stack []interface{}) ([]interface{}, string) {
    var result bool
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(int) <= stack[1].(int)
                case float64:
                    result = float64(stack[0].(int)) <= stack[1].(float64)
                case string:
                    shalePanic("<=: Cannot compare int with str")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(float64) <= float64(stack[1].(int))
                case float64:
                    result = stack[0].(float64) <= stack[1].(float64)
                case string:
                    shalePanic("<=: Cannot compare float with str")
            }
        case string:
            switch stack[1].(type) {
                case string:
                    result = stack[0].(string) <= stack[1].(string)
                case int:
                    shalePanic("<=: Cannot compare str with int")
                case float64:
                    shalePanic("<=: Cannot compare str with float")
            }
    }
    return []interface{}{result,}, ""
}

func opGtEqu(stack []interface{}) ([]interface{}, string) {
    var result bool
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(int) >= stack[1].(int)
                case float64:
                    result = float64(stack[0].(int)) >= stack[1].(float64)
                case string:
                    shalePanic(">=: Cannot compare int with str")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    result = stack[0].(float64) >= float64(stack[1].(int))
                case float64:
                    result = stack[0].(float64) >= stack[1].(float64)
                case string:
                    shalePanic(">=: Cannot compare float with str")
            }
        case string:
            switch stack[1].(type) {
                case string:
                    result = stack[0].(string) >= stack[1].(string)
                case int:
                    shalePanic(">=: Cannot compare str with int")
                case float64:
                    shalePanic(">=: Cannot compare str with float")
            }
    }
    return []interface{}{result,}, ""
}

func opNotEqu(stack []interface{}) ([]interface{}, string) {
    if stack[0] != stack[1] {
        return []interface{}{true,}, ""
    } else {
        return []interface{}{false,}, ""
    }
}

func opCmp(stack []interface{}) ([]interface{}, string) {
    var result int
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    if stack[0].(int) == stack[1].(int) {
                        result = 0
                    } else if stack[0].(int) < stack[1].(int) {
                        result = -1
                    } else if stack[0].(int) > stack[1].(int) {
                        result = 1
                    }
                case float64:
                    if float64(stack[0].(int)) == stack[1].(float64) {
                        result = 0
                    } else if float64(stack[0].(int)) < stack[1].(float64) {
                        result = -1
                    } else if float64(stack[0].(int)) > stack[1].(float64) {
                        result = 1
                    }
                case string:
                    shalePanic("<=>: Cannot compare int with str")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    if stack[0].(float64) == float64(stack[1].(int)) {
                        result = 0
                    } else if stack[0].(float64) < float64(stack[1].(int)) {
                        result = -1
                    } else if stack[0].(float64) > float64(stack[1].(int)) {
                        result = 1
                    }
                case float64:
                    if stack[0].(float64) == stack[1].(float64) {
                        result = 0
                    } else if stack[0].(float64) < stack[1].(float64) {
                        result = -1
                    } else if stack[0].(float64) > stack[1].(float64) {
                        result = 1
                    }
                case string:
                    shalePanic("<=>: Cannot compare float with str")
            }
        case string:
            switch stack[1].(type) {
                case string:
                    if stack[0].(string) == stack[1].(string) {
                        result = 0
                    } else if stack[0].(string) < stack[1].(string) {
                        result = -1
                    } else if stack[0].(string) > stack[1].(string) {
                        result = 1
                    }
                case int:
                    shalePanic("<=>: Cannot compare str with int")
                case float64:
                    shalePanic("<=>: Cannot compare str with float")
            }
    }
    return []interface{}{result,}, ""
}

func opInput(stack []interface{}) ([]interface{}, string) {
    if len(stack) == 1 {
        fmt.Print(stack[0])
    }
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    input = strings.Replace(input, "\n", "", -1)
    return []interface{}{input,}, ""
}

func opPop(stack []interface{}) ([]interface{}, string) {
    if len(stack) == 0 {
        return nil, "StackEmpty"
    }
    return stack[:len(stack) - 1], ""
}

func opClear(stack []interface{}) ([]interface{}, string) {
    return nil, ""
}

func opRead(stack []interface{}) ([]interface{}, string) {
    var brErr string
    text, err := ioutil.ReadFile(stack[0].(string))
    if err != nil {
        brErr = "ReadError"
    }
    return []interface{}{string(text),}, brErr
}

func opWrite(stack []interface{}) ([]interface{}, string) {
    var brErr string
	f, err := os.Create(stack[0].(string))
    if err != nil {
        return nil, "FileError"
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    _, err = w.WriteString(stack[1].(string))
    w.Flush()
    if err != nil {
        brErr = "WriteError"
    }
    return nil, brErr
}

func opSplit(stack []interface{}) ([]interface{}, string) {
    split := strings.Split(stack[0].(string), stack[1].(string))
    var substrings []interface{}
    for i := 0; i < len(split); i++ {
        substrings = append(substrings, split[i])
    }
    return substrings, ""
}

func opVec(stack []interface{}) ([]interface{}, string) {
    vec := []interface{}{}
    for i := 0; i < len(stack); i++ {
        if stack[i] == nil {
            vec = append(vec, nil)
        } else {
            switch stack[i].(type) {
                case int:
                    vec = append(vec, stack[i].(int))
                case float64:
                    vec = append(vec, stack[i].(float64))
                case string:
                    vec = append(vec, stack[i].(string))
                case bool:
                    vec = append(vec, stack[i].(bool))
                case []interface{}:
                    vec = append(vec, stack[i].([]interface{}))
            }
        }
    }
    return []interface{}{vec,}, ""
}

func opGetChar(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case string:
            switch stack[1].(type) {
                case int:
                    s := string(stack[0].(string)[stack[1].(int)])
                    return []interface{}{s,}, ""
                default:
                    shalePanic("getChar: Inoperable type")
            }
        default:
            shalePanic("getChar: Inoperable type")
    }
    return nil, ""
}

func opReplace(stack []interface{}) ([]interface{}, string) {
    s := strings.Replace(stack[0].(string), stack[1].(string), stack[2].(string), -1)
    return []interface{}{s,}, ""
}

func opAppend(stack []interface{}) ([]interface{}, string) {
    vec := stack[0].([]interface{})
    switch stack[1].(type) {
        case int:
            vec = append(stack[0].([]interface{}), stack[1].(int))
        case float64:
            vec =  append(stack[0].([]interface{}), stack[1].(float64))
        case string:
            vec = append(stack[0].([]interface{}), stack[1].(string))
        case []interface{}:
            vec = append(stack[0].([]interface{}), stack[1].([]interface{}))
    }
    return []interface{}{vec,}, ""
}

func opStderr(stack []interface{}) ([]interface{}, string) {
    for i := 0; i < len(stack); i++ {
        fmt.Fprintln(os.Stderr, stack[i])
    }
    return nil, ""
}

func opLen(stack []interface{}) ([]interface{}, string) {
    return []interface{}{len(stack),}, ""
}

func opSin(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            return []interface{}{math.Sin(float64(stack[0].(int))),}, ""
        case float64:
            return []interface{}{math.Sin(stack[0].(float64)),}, ""
        default:
            shalePanic("sin: Inoperable type")
    }
    return nil, ""
}

func opCos(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            return []interface{}{math.Cos(float64(stack[0].(int))),}, ""
        case float64:
            return []interface{}{math.Cos(stack[0].(float64)),}, ""
        default:
            shalePanic("sin: Inoperable type")
    }
    return nil, ""
}

func opTan(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            return []interface{}{math.Tan(float64(stack[0].(int))),}, ""
        case float64:
            return []interface{}{math.Tan(stack[0].(float64)),}, ""
        default:
            shalePanic("sin: Inoperable type")
    }
    return nil, ""
}

func opType(stack []interface{}) ([]interface{}, string) {
    var result string
    if stack[0] == nil {
        result = "null"
    } else {
        switch stack[0].(type) {
            case int:
                result = "int"
            case float64:
                result = "float"
            case string:
                result = "str"
            case bool:
                result = "bool"
            case []interface{}:
                result = "vec"
        }
    }
    return []interface{}{result,}, ""
}

func opPow(stack []interface{}) ([]interface{}, string) {
    switch stack[0].(type) {
        case int:
            switch stack[1].(type) {
                case int:
                    return []interface{}{math.Pow(float64(stack[0].(int)), float64(stack[1].(int))),}, ""
                case float64:
                    return []interface{}{math.Pow(float64(stack[0].(int)), stack[1].(float64)),}, ""
                default:
                    shalePanic("^: Inoperable type")
            }
        case float64:
            switch stack[1].(type) {
                case int:
                    return []interface{}{math.Pow(stack[0].(float64), float64(stack[1].(int))),}, ""
                case float64:
                    return []interface{}{math.Pow(stack[0].(float64), stack[1].(float64)),}, ""
                default:
                    shalePanic("^: Inoperable type")
            }
        default:
            shalePanic("^: Inoperable type")
    }
    return nil, ""
}

func opNot(stack []interface{}) ([]interface{}, string) {
    var result bool
    switch stack[0] {
        case true:
            result = false
        case false:
            result = true
        default:
            shalePanic("!: Inoperable type")
    }
    return []interface{}{result,}, ""
}

func opEqual(stack []interface{}) ([]interface{}, string) {
    var result bool
    if stack[0] == 0 {
        result = true
    } else {
        result = false
    }
    return []interface{}{result,}, ""
}

func opLess(stack []interface{}) ([]interface{}, string) {
    var result bool
    if stack[0] == -1 {
        result = true
    } else {
        result = false
    }
    return []interface{}{result,}, ""
}

func opGreater(stack []interface{}) ([]interface{}, string) {
    var result bool
    if stack[0] == 1 {
        result = true
    } else {
        result = false
    }
    return []interface{}{result,}, ""
}

func opExtend(stack []interface{}) ([]interface{}, string) {
    var vec []interface{}
    switch stack[0].(type) {
        case []interface{}:
            switch stack[1].(type) {
                case []interface{}:
                    vec = append(stack[0].([]interface{}), stack[1].([]interface{})...)
                default:
                    shalePanic("extend: Inoperabale type")
            }
        default:
            shalePanic("extend: Inoperable type")
    }
    return []interface{}{vec,}, ""
}

func opExpand(stack []interface{}) ([]interface{}, string) {
    var vec []interface{}
    switch stack[0].(type) {
        case []interface{}:
            vec = stack[0].([]interface{})
        default:
            shalePanic("expand: Inoperable type")
    }
    return vec, ""
}

func opSum(stack []interface{}) ([]interface{}, string) {
    var returnFloat bool
    var result float64
    switch stack[0].(type) {
        case []interface{}:
            for i := 0; i < len(stack[0].([]interface{})); i++ {
                switch stack[0].([]interface{})[i].(type) {
                    case int:
                        result += float64(stack[0].([]interface{})[i].(int))
                    case float64:
                        returnFloat = true
                        result += stack[0].([]interface{})[i].(float64)
                    default:
                        shalePanic("sum: Inoperable type in op 1")
                }
            }
        default:
            shalePanic("sum: Inoperable type")
    }
    if returnFloat {
        return []interface{}{result,}, ""
    }
    return []interface{}{int(result),}, ""
}

func opAssert(stack []interface{}) ([]interface{}, string) {
    if stack[0] != true {
        shalePanic("assert: Not true")
    }
    return nil, ""
}

func opUpcase(stack []interface{}) ([]interface{}, string) {
    var s string
    switch stack[0].(type) {
        case string:
            s = strings.ToUpper(stack[0].(string))
        default:
            shalePanic("upcase: Inoperable type")
    }
    return []interface{}{s,}, ""
}

func opDowncase(stack []interface{}) ([]interface{}, string) {
    var s string
    switch stack[0].(type) {
        case string:
            s = strings.ToLower(stack[0].(string))
        default:
            shalePanic("downcase: Inoperable type")
    }
    return []interface{}{s,}, ""
}
