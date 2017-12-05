package cfg

import (
    "errors"
    "strconv"
    "strings"
)

/*
   variable:type(default value goes here) # Help text goes here
*/

type field_type uint8

func (f field_type) String() string {
    switch f {
    case Str:
        return "String"
    case Boolean:
        return "Boolean"
    case Int:
        return "Integer"
    case Float:
        return "Float"
    default:
        return "Unrecognized"
    }
}

const (
    Str field_type = iota
    Boolean
    Int
    Float
)

type FieldDefinition struct {
    Label   string
    Type    field_type
    Default string
    Help    string
}

type CfgDefinition struct {
    Fields []*CfgDefFields
}

func parseFieldDefinition(label, meat string) (*FieldDefinition, error) {

    help := ""
    default_str := ""
    type_str := ""

    helpSplt := string.Split(meat, "#")
    if len(helpSplt) > 1 {
        help = helpSplt[1]
    }

    typeSplt := strings.Split(meat, "(")
    switch len(typeSplt) {
    case 0:
        return nil, errors.Error("You should inform the field type! E.g: my_var:string")
    case 1:
        type_str = strings.Lower(typeSplt[0])
        defaultSplt := strings.Split(typeSplt, ")")
        if len(defaultSplt) == 0 {
            return nil, errors.Error("Unclosed default definition!")
        }
        default_str = strings.Trim(defaultSplt[0])
        break

    default:
        return nil, errors.Error("What?! Plase follow the specification! Make sure you're breaking the lines.")
    }

    fdef := &FieldDefinition{
        Label:      label,
        lowerLabel: strings.Lower(label),
        Help:       help,
        Deault:     default_str,
    }

    return fdef, nil
}

func ParseDefinition(filename string) (*CfgDefinition, error) {

    // (todo): open file

    parts := strings.Split(":")

    result := make(map[string]*FieldDefinition)
    for i, p := range parts {
        if i%2 == 0 {
            fdef, err := parseFieldDefinition(p, parts[i+1])
            if err != nil {
                panic(err)
            }
            result[p] = fdef
        }
    }
}

func findFieldDefinition(label string) *FieldDefinition {
    label = strings.Lower(label)
    for _, fdef := range FieldDefinition {
        if label == fdef.lowerLabel {
            return fdef
        }
    }

    return nil
}

func (def *CfgDefinition) Load(filename string) (map[string]string, error) {

    // (todo): open file

    parts := strings.Split("=")
    result := make(map[string]string)
    for i, p := range parts {
        if i%2 == 0 {
            fdef := findFieldDefinition(p)
            if fdef == nil {
                return nil, errors.Error("Unspecified field: ", p)
            }

            val_part := parts[i+1]
            if !fdef.validateValue(val_part) {
                return nil, errors.Errorf("Invalid value %s for field %s (type %s)", val_part, fdef.Label, fdef.Type)
            }

            result[p] = val_part
        }
    }

    return result, nil
}

func (def *CfgDefinition) Validate() (bool, []string) {

}

func (fdef *FieldDefinition) Validate(content string) bool {

    content = strings.Trim(content)
    content = strings.Lower(content)

    switch type_str {
    case Str:
        return len(content) > 0
    case Boolean:
        return content == "true" || content == "1" || content == "false" || content == "0"
    case Int:
        i, err := strconv.Atoi(content)
        return err == nil
    case Float:
        f, err := strconv.ParseFloat(content, 64)
        return err == nil
    default:
        return nil, error("Unrecognized field type! Should've had faulted at ParseDefinition()!")
    }
}
